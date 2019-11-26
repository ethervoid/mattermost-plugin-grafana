package main

import (
	"bytes"
	"encoding/json"
	"strconv"
)

const (
	// SUBSCRIPTIONSKEY blabla
	SUBSCRIPTIONSKEY = "subscriptions"
)

// Subscription blabla
type Subscription struct {
	ChannelID         string
	PanelURL          string
	RefreshPeriod     int
	LastTimeRefreshed int
}

// Subscribe blabla
func (p *Plugin) Subscribe(panelURL string, channel string, refreshPeriod string) error {
	refreshPeriodNumeric, err := strconv.Atoi(refreshPeriod)
	if err != nil {
		return err
	}
	newSubscription := &Subscription{
		ChannelID:         channel,
		PanelURL:          panelURL,
		RefreshPeriod:     refreshPeriodNumeric,
		LastTimeRefreshed: -1,
	}
	subscriptions, err := p.GetSubscriptions()
	if err != nil {
		return err
	}

	if subscriptions == nil {
		subscriptions = []*Subscription{newSubscription}
	}
	exists := false
	for index, s := range subscriptions {
		if s.ChannelID == channel {
			subscriptions[index] = newSubscription
			exists = true
			break
		}
	}

	if !exists {
		subscriptions = append(subscriptions, newSubscription)
	}

	err = p.storeSubscriptions(subscriptions)
	if err != nil {
		return err
	}
	return nil
}

func (p *Plugin) storeSubscriptions(subscriptions []*Subscription) error {
	b, err := json.Marshal(subscriptions)
	if err != nil {
		return err
	}
	p.API.KVSet(SUBSCRIPTIONSKEY, b)
	return nil
}

// GetSubscription blabla
func (p *Plugin) GetSubscription(channelID string) (*Subscription, error) {
	subscriptions, err := p.GetSubscriptions()
	if err != nil {
		return nil, err
	}

	for index, s := range subscriptions {
		if s.ChannelID == channelID {
			return subscriptions[index], nil
		}
	}

	return nil, nil
}

// GetSubscriptions blabla
func (p *Plugin) GetSubscriptions() ([]*Subscription, error) {
	var subscriptions []*Subscription

	value, err := p.API.KVGet(SUBSCRIPTIONSKEY)
	if err != nil {
		return nil, err
	}

	if value == nil {
		subscriptions = []*Subscription{}
	} else {
		json.NewDecoder(bytes.NewReader(value)).Decode(&subscriptions)
	}

	return subscriptions, nil
}

// Unsubscribe blabla
func (p *Plugin) Unsubscribe(channelID string) error {
	subscriptions, err := p.GetSubscriptions()
	if err != nil {
		return err
	}

	for index, s := range subscriptions {
		if s.ChannelID == channelID {
			subscriptions = append(subscriptions[:index], subscriptions[index+1:]...)
		}
	}

	err = p.storeSubscriptions(subscriptions)
	if err != nil {
		return err
	}

	return nil
}
