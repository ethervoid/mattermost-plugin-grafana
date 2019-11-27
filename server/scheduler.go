package main

import "time"

import "github.com/mattermost/mattermost-server/model"

// RunScheduler blabla
func (p *Plugin) RunScheduler() {
	go func() {
		<-time.NewTimer(10 * time.Second).C
		p.RefreshSubscriptions()
		p.RunScheduler()
	}()
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

	p.API.PublishWebSocketEvent("update_grafana_subscription",
		payload,
		//&model.WebsocketBroadcast{ChannelId: subscription.ChannelID})
		&model.WebsocketBroadcast{UserId: "e91wpdsqkjb7zm9dci8tj9ms9y"})

	p.API.LogInfo("Image sent through the websocket...")

	return nil
}
