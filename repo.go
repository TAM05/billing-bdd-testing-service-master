package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/utilitywarehouse/uw-lib-billing/usage"
)

type pricedEventFinder interface {
	getPricedEvent(string) (usage.Event, error)
}

type pricedEventsRepo struct {
	client *http.Client
	addr   string
}

func newPricedEventsRepo(client *http.Client, addr string) pricedEventsRepo {
	return pricedEventsRepo{client, addr}
}

func (r pricedEventsRepo) getPricedEvent(id string) (usage.Event, error) {
	resp, err := r.client.Get(r.addr + "/api/1.0/events/priced/" + url.PathEscape(id))
	if err != nil {
		return usage.Event{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return usage.Event{}, fmt.Errorf("Expected response 200 got %d", resp.StatusCode)
	}
	var event usage.Event
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		return usage.Event{}, err
	}
	return event, nil
}
