package main

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/utilitywarehouse/go-pubsub/mockqueue"
	"github.com/utilitywarehouse/uw-lib-billing/usage"
)

func TestHandler(t *testing.T) {
	var testCases = []struct {
		name             string
		calculatedPrice  string
		expectedResponse string
	}{
		{
			name:             "happy case",
			calculatedPrice:  "1.55",
			expectedResponse: `{"result":"success"}`,
		},
		{
			name:             "failing case",
			calculatedPrice:  "1.56",
			expectedResponse: `{"result":"failure"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tester := newTester(
				withSink(mockqueue.NewMockQueue()),
				withRepo(stubPricedEventsRepo{usage.Event{Price: &usage.Amount{Value: newDecimal(t, tc.calculatedPrice)}}}),
				withWaitDuration(time.Millisecond),
				withTags("@test"),
			)

			w := httptest.NewRecorder()

			handler(tester)(w, nil)

			assert.Equal(t, tc.expectedResponse, w.Body.String())
		})
	}
}

type stubPricedEventsRepo struct {
	event usage.Event
}

func (s stubPricedEventsRepo) getPricedEvent(id string) (usage.Event, error) {
	return s.event, nil
}
