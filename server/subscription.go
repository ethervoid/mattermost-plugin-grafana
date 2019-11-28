package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

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
	TimeRange         int
	PanelWidth        int
	PanelHeight       int
	LastTimeRefreshed int
}

// Subscribe blabla
func (p *Plugin) Subscribe(teamID string, channelName string, panelURL string, timeRange string) error {
	timeRangeInMinutes, err := strconv.Atoi(timeRange)
	if err != nil {
		return err
	}
	timeRangeInSeconds := timeRangeInMinutes * 60
	channel, appError := p.API.GetChannelByName(teamID, channelName, false)
	if appError != nil {
		return errors.New(appError.Error())
	}
	newSubscription := &Subscription{
		ChannelID:         channel.Id,
		ChannelName:       channel.Name,
		PanelURL:          panelURL,
		TimeRange:         timeRangeInSeconds,
		PanelWidth:        400,
		PanelHeight:       200,
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
	parsedURL, err := p.prepareSubscriptionURL(subscription)
	if err != nil {
		return err
	}
	panelImage, err := p.loadImageFromURL(parsedURL)
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

func (p *Plugin) prepareSubscriptionURL(subscription *Subscription) (string, error) {
	urlParsed, err := url.Parse(subscription.PanelURL)
	if err != nil {
		return "", err
	}
	urlQuery := urlParsed.Query()
	now := time.Now()
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", subscription.TimeRange))
	then := now.Add(-duration)
	urlQuery.Set("from", fmt.Sprintf("%d", then.UnixNano()/1000000))
	urlQuery.Set("to", fmt.Sprintf("%d", now.UnixNano()/1000000))
	urlQuery.Set("width", fmt.Sprintf("%d", subscription.PanelWidth))
	urlQuery.Set("height", fmt.Sprintf("%d", subscription.PanelHeight))
	urlParsed.RawQuery = urlQuery.Encode()
	p.API.LogInfo("Parsed URL: " + urlParsed.String())
	return urlParsed.String(), nil
}
