package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/utilitywarehouse/go-operational-health-checks/healthcheck"
	pubsub "github.com/utilitywarehouse/go-pubsub"
	"github.com/utilitywarehouse/go-pubsub/kafka"
	"github.com/utilitywarehouse/uw-lib-billing/server"
	"github.com/utilitywarehouse/uwgolib/log"
)

const appName = "billing-end-to-end-testing-service"
const appDescription = "Runs end-to-end behaviour driven tests defined"

var revision = "overridden at build time"

func main() {
	run(os.Args, nil)
}

func run(args []string, ready chan<- bool) {
	app := cli.App(appName, appDescription)
	brokers := app.String(cli.StringOpt{
		Name:   "brokers",
		Value:  "localhost:9092",
		Desc:   "Comma separated array of broker host and port",
		EnvVar: "BROKERS",
	})
	destinationTopic := app.String(cli.StringOpt{
		Name:   "destination-topic",
		Value:  "uw.billing.usage-events.raw",
		Desc:   "Topic to write usage events to",
		EnvVar: "DESTINATION_TOPIC",
	})
	pricedEventsAPIURL := app.String(cli.StringOpt{
		Name:   "priced-events-api-url",
		Value:  "http://uw-service-priced-events:80",
		Desc:   "Address of the priced events API service",
		EnvVar: "PRICED_EVENTS_API_URL",
	})
	waitDuration := app.String(cli.StringOpt{
		Name:   "wait-duration",
		Value:  "10s",
		Desc:   "Duration to wait before checking the rated event",
		EnvVar: "WAIT_DURATION",
	})
	serverConfig := server.ExternalHTTPConfig(app)
	log.ExternalConfig(app)

	app.Action = func() {
		ctx := context.Background()
		errors := make(chan error, 10)
		go handleErrors(errors)

		sink := initMessageSink(*brokers, *destinationTopic)

		d, err := time.ParseDuration(*waitDuration)
		if err != nil {
			log.Fatal(err)
		}
		httpClient := initHTTPClient()
		repo := newPricedEventsRepo(httpClient, *pricedEventsAPIURL)
		t := newTester(withSink(sink), withRepo(repo), withWaitDuration(d))

		m := mux.NewRouter()
		m.PathPrefix("/api/1.0/run-test").Methods(http.MethodPost).HandlerFunc(prometheus.InstrumentHandlerFunc("/api/1.0/run-test", handler(t)))

		server.AddStatusEndpoints(m, server.GetInstrumentedStatus(appName, appDescription, revision).
			AddChecker("priced events api check", healthcheck.NewHTTPAPIHealthCheck(httpClient, *pricedEventsAPIURL+"/__/health", "End-to-end tests are going to fail")).
			AddChecker("broker check", healthcheck.NewPubSubCheck(sink, "End-to-end tests are going to fail")))

		s := server.NewHTTPServer(serverConfig, m)
		go func() {
			log.Infof("Listening on %s", serverConfig.Addr)
			if err := s.Start(); err != nil {
				log.Error(err)
			}
		}()

		if ready != nil {
			close(ready)
		}

		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		<-ch
		log.Infof("Shutting down application...")
		s.Shutdown(ctx)

		if err := sink.Close(); err != nil {
			log.Errorf("Error while closing kafka forwarder: error=(%v).", err)
		}
	}
	log.Infof("Starting %s version %s ", appName, revision)
	err := app.Run(args)
	if err != nil {
		log.Fatalf("App exited with error: [%+v]", err)
	}
	log.Infof("%s shuts down normally", appName)
}

func initHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 128,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
		},
	}
}

func handler(t tester) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status := godog.RunWithOptions("godogs", newFeatureContext(&t), godog.Options{
			Format:    "pretty",
			Paths:     []string{"features"},
			Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
			Tags:      t.tags,
		})
		if status != 0 {
			w.Write([]byte(`{"result":"failure"}`))
		} else {
			w.Write([]byte(`{"result":"success"}`))
		}
	}
}

func initMessageSink(brokers, destinationTopic string) pubsub.MessageSink {
	sink, err := kafka.NewMessageSink(
		kafka.MessageSinkConfig{
			Topic:   destinationTopic,
			Brokers: strings.Split(brokers, ","),
			KeyFunc: func(m pubsub.ProducerMessage) []byte {
				return nil
			},
		})
	if err != nil {
		log.Fatalf("Could not create kafka sink: error=(%v)", err)
	}
	return sink
}

func onError(errors chan error, err error) {
	select {
	case errors <- err:
	default:
	}
}

func handleErrors(errors chan error) {
	for e := range errors {
		log.Errorf("error=(%v)", e)
	}
}
