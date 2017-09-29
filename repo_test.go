package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/utilitywarehouse/uw-lib-billing/usage"
)

func TestGetPricedEvent(t *testing.T) {
	id := uuid.New().String()
	event := usage.Event{
		EventID: id,
		Device:  usage.Device{ID: "+447700900077"},
	}
	b, err := json.Marshal(event)
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/1.0/events/priced/"+id {
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer ts.Close()

	repo := newPricedEventsRepo(http.DefaultClient, ts.URL)

	actual, err := repo.getPricedEvent(id)
	require.NoError(t, err)
	assert.Equal(t, event, actual)
}
