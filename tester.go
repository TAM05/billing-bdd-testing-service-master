package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/utilitywarehouse/uwgolib/log"

	"github.com/DATA-DOG/godog"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	pubsub "github.com/utilitywarehouse/go-pubsub"
	"github.com/utilitywarehouse/uw-cdr/cdr"
)

type tester struct {
	sink pubsub.MessageSink
	repo pricedEventFinder

	waitDur   time.Duration
	tags      string
	currentID string
}

func withWaitDuration(d time.Duration) func(t *tester) {
	return func(t *tester) {
		t.waitDur = d
	}
}

func withSink(s pubsub.MessageSink) func(t *tester) {
	return func(t *tester) {
		t.sink = s
	}
}
func withRepo(r pricedEventFinder) func(t *tester) {
	return func(t *tester) {
		t.repo = r
	}
}

func withCurrentID(id string) func(t *tester) {
	return func(t *tester) {
		t.currentID = id
	}
}

func withTags(tags string) func(t *tester) {
	return func(t *tester) {
		t.tags = tags
	}
}

func newTester(options ...func(*tester)) tester {
	t := tester{
		waitDur: 10 * time.Second,
		tags:    "~@test",
	}
	for _, o := range options {
		o(&t)
	}
	return t
}

var commTypeRegexp = regexp.MustCompile(`^(\w+) \((\d+) \w+\)$`)

func (t *tester) spawnUsageEvent(moCLI, retailBand, commType string) error {
	s := commTypeRegexp.FindStringSubmatch(commType)
	if len(s) != 3 {
		return fmt.Errorf("Invalid 'Communication type' format %v", commType)
	}

	usage, err := strconv.Atoi(s[2])
	if err != nil {
		return fmt.Errorf("Invalid usage: %v", s[2])
	}
	service, err := getService(s[1], usage)
	if err != nil {
		return err
	}

	rawEvent := cdr.UsageRecord{
		EventId:    t.currentID,
		EventStart: time.Now().UTC().Format(time.RFC3339),
		RetailBand: retailBand,
		Service:    &service,
		Subscriber: &cdr.Subscriber{
			Cli: moCLI,
		},
		Provider: &cdr.Provider{
			ID:   "ee",
			Type: "??",
		},
	}

	b, err := json.Marshal(rawEvent)
	if err != nil {
		return err
	}

	return t.sink.PutMessage(pubsub.SimpleProducerMessage(b))
}

func (t *tester) getTotalRate(rate string) error {
	expected, err := decimal.NewFromString(rate)
	if err != nil {
		log.Infof("Event %v. %v", t.currentID, err.Error())
		return err
	}

	deadline := time.Now().Add(t.waitDur)
	for time.Now().Before(deadline) {
		time.Sleep(time.Second)
		pricedEvent, err := t.repo.getPricedEvent(t.currentID)
		if err != nil {
			log.Infof("Event %v. %v", t.currentID, err.Error())
			continue
		}

		if !pricedEvent.Price.Value.Equal(expected) {
			s := fmt.Sprintf("Expected %v, but actual price is %v", expected, pricedEvent.Price.Value)
			log.Warnf("Event %v. %v", t.currentID, s)
			return errors.New(s)
		}
		return nil
	}
	msg := "Get total rate failure. Wait duration has passed"
	log.Warnf("Event %v. %v. Rating either have not taken place or is incorrect", t.currentID, msg)
	return errors.New(msg)
}

func getService(serviceType string, usage int) (cdr.Service, error) {
	otherParty := cdr.DirectionOtherParty{
		OtherParty: &cdr.OtherParty{
			Cli: "Some other party CLI",
		},
	}

	switch cdr.ServiceType(serviceType) {
	case cdr.VOICE:
		return cdr.Service{
			Id: cdr.ServiceType(serviceType),
			Voice: &cdr.VoiceService{
				DirectionOtherParty: otherParty,
				DurationSeconds:     usage,
			},
		}, nil
	case cdr.SMS:
		return cdr.Service{
			Id: cdr.ServiceType(serviceType),
			Sms: &cdr.SmsService{
				DirectionOtherParty: otherParty,
			},
		}, nil
	case cdr.MMS:
		return cdr.Service{
			Id: cdr.ServiceType(serviceType),
			Mms: &cdr.MmsService{
				DirectionOtherParty: otherParty,
			},
		}, nil
	case cdr.DATA:
		return cdr.Service{
			Id: cdr.ServiceType(serviceType),
			Data: &cdr.DataService{
				VolumeKB: usage,
			},
		}, nil
	default:
		return cdr.Service{}, fmt.Errorf("Unknown service type: %v", serviceType)
	}
}

func newFeatureContext(t *tester) func(*godog.Suite) {
	return func(s *godog.Suite) {
		s.Step(`^MO is on tariff .*$`, func() error { return nil })
		s.Step(`^MO with number ([+0-9a-zA-Z_]+) connects to MT via (\w+) and (.*)$`, t.spawnUsageEvent)
		s.Step(`^get result and ([0-9.]+)`, t.getTotalRate)

		s.BeforeScenario(func(interface{}) {
			uid := uuid.New().String() + "-test"
			t.currentID = uid
			log.Printf("Generated event with uuid %v", t.currentID)
		})
		s.AfterScenario(func(interface{}, error) {
			//TODO remove raw & priced usage from the db
		})
	}
}
