package main

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pubsub "github.com/utilitywarehouse/go-pubsub"
	"github.com/utilitywarehouse/go-pubsub/mockqueue"
	"github.com/utilitywarehouse/uw-cdr/cdr"
	"github.com/utilitywarehouse/uw-lib-billing/usage"
)

func TestSpawnUsageEvent(t *testing.T) {
	mq := mockqueue.NewMockQueue()
	id := uuid.New().String()
	tester := newTester(withSink(mq), withCurrentID(id))

	err := tester.spawnUsageEvent("+440123401234", "m121", "voice (65 seconds)")
	require.NoError(t, err)

	var event cdr.UsageRecord
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	mq.ConsumeMessages(ctx, func(m pubsub.ConsumerMessage) error { return json.Unmarshal(m.Data, &event) }, nil)

	assert.Equal(t, id, event.EventId)
	assert.Equal(t, event.RetailBand, "m121")
	assert.Equal(t, *event.Service, cdr.Service{
		Id: "voice",
		Voice: &cdr.VoiceService{
			DirectionOtherParty: cdr.DirectionOtherParty{
				OtherParty: &cdr.OtherParty{
					Cli: "Some other party CLI",
				},
			},
			DurationSeconds: 65,
		},
	})
	_, err = time.Parse(time.RFC3339, event.EventStart)
	assert.NoError(t, err)
}

func TestGetTotalRate(t *testing.T) {
	var cases = []struct {
		name        string
		charge      string
		currentID   string
		event       usage.Event
		expectedErr error
	}{
		{
			name:      "happy case",
			charge:    "0.45",
			currentID: "foobar",
			event: usage.Event{
				EventID: "foobar",
				Price: &usage.Amount{
					Value: newDecimal(t, "0.45"),
				},
			},
			expectedErr: nil,
		},
		{
			name:      "price differs",
			charge:    "10.1",
			currentID: "foobar",
			event: usage.Event{
				EventID: "foobar",
				Price: &usage.Amount{
					Value: newDecimal(t, "0.45"),
				},
			},
			expectedErr: errors.New("Expected 10.1, but actual price is 0.45"),
		},
		{
			name:      "event not found",
			charge:    "10.1",
			currentID: "booha",
			event: usage.Event{
				EventID: "foobar",
				Price: &usage.Amount{
					Value: newDecimal(t, "0.45"),
				},
			},
			expectedErr: errors.New("Get total rate failure. Wait duration has passed"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tester := newTester(
				withRepo(mockPricedEventsRepo{map[string]usage.Event{tc.event.EventID: tc.event}}),
				withWaitDuration(time.Millisecond),
				withCurrentID(tc.currentID),
			)

			err := tester.getTotalRate(tc.charge)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func newDecimal(t *testing.T, price string) decimal.Decimal {
	val, err := decimal.NewFromString(price)
	if err != nil {
		t.Fatal(err)
	}
	return val
}

type mockPricedEventsRepo struct {
	store map[string]usage.Event
}

func (m mockPricedEventsRepo) getPricedEvent(id string) (usage.Event, error) {
	ev, ok := m.store[id]
	if !ok {
		return usage.Event{}, errors.New("event not found")
	}
	return ev, nil
}
