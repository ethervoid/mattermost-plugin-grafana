package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/mattermost/mattermost-server/model"
)

const (
	// SUBSCRIPTIONSKEY blabla
	SUBSCRIPTIONSKEY = "subscriptions"
)

// Subscription blabla
type Subscription struct {
	ChannelID         string
	ChannelName       string
	PanelURL          string
	RefreshPeriod     int
	LastTimeRefreshed int
}

// Subscribe blabla
func (p *Plugin) Subscribe(teamID string, channelName string, panelURL string, refreshPeriod string) error {
	refreshPeriodNumeric, err := strconv.Atoi(refreshPeriod)
	if err != nil {
		return err
	}
	channel, appError := p.API.GetChannelByName(teamID, channelName, false)
	if appError != nil {
		return errors.New(appError.Error())
	}
	newSubscription := &Subscription{
		ChannelID:         channel.Id,
		ChannelName:       channel.Name,
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
		if s.ChannelID == channel.Id {
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
func (p *Plugin) Unsubscribe(teamID string, channelName string) error {
	subscriptions, err := p.GetSubscriptions()
	if err != nil {
		return err
	}
	channel, appError := p.API.GetChannelByName(teamID, channelName, false)
	if appError != nil {
		return errors.New(appError.Error())
	}

	for index, s := range subscriptions {
		if s.ChannelID == channel.Id {
			subscriptions = append(subscriptions[:index], subscriptions[index+1:]...)
		}
	}

	err = p.storeSubscriptions(subscriptions)
	if err != nil {
		return err
	}

	payload := make(map[string]interface{})
	payload["channelTarget"] = channel.Id

	p.API.PublishWebSocketEvent("remove_grafana_subscription",
		payload,
		&model.WebsocketBroadcast{ChannelId: channel.Id})

	return nil
}

// RefreshSubscriptions blabla
func (p *Plugin) RefreshSubscriptions() error {
	// TODO Get current time and refresh
	subscriptions, err := p.GetSubscriptions()
	if err != nil {
		return err
	}

	if len(subscriptions) > 0 {
		p.RefreshSubscription(subscriptions[0])
	}

	return nil
}

// RefreshSubscription blabla
func (p *Plugin) RefreshSubscription(subscription *Subscription) error {
	// TODO Retrieve the Grafana image and send it again using the websocket
	panelImage, err := p.loadImageFromURL(subscription.PanelURL)
	if err != nil {
		return err
	}

	panelImageBase64, err := p.encodeImageToBase64(panelImage)
	if err != nil {
		return err
	}

	payload := make(map[string]interface{})
	payload["image"] = panelImageBase64
	payload["channelTarget"] = subscription.ChannelID

	p.API.PublishWebSocketEvent("update_grafana_subscription",
		payload,
		&model.WebsocketBroadcast{ChannelId: subscription.ChannelID})

	p.API.LogInfo("Image sent through the websocket...")

	return nil
}
