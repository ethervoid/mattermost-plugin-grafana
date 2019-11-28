package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "grafana",
		DisplayName:      "Grafana",
		Description:      "Integration with Grafana.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: connect, disconnect, list, register",
		AutoCompleteHint: "[command]",
	}
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

// ExecuteCommand blabla
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	parameters := []string{}
	action := ""
	if len(split) > 1 {
		action = split[1]
	}
	if len(split) > 2 {
		parameters = split[2:]
	}

	if command != "/grafana" {
		return &model.CommandResponse{}, nil
	}

	switch action {
	case "subscribe":
		//https://play.grafana.org/render/d-solo/000000012/grafana-play-home?orgId=1&from=&to=1574720321201&panelId=2&width=1000&height=500&tz=Europe%2FMadrid
		if len(parameters) != 3 {
			p.postCommandResponse(args, "You have to pass the panel channel, URL and refresh interval parameters")
			return &model.CommandResponse{}, nil
		}
		err := p.Subscribe(args.TeamId, parameters[0], parameters[1], parameters[2])
		if err != nil {
			p.postCommandResponse(args, err.Error())
			return &model.CommandResponse{}, nil
		}
		err = p.RefreshSubscriptions()
		if err != nil {
			p.postCommandResponse(args, err.Error())
			return &model.CommandResponse{}, nil
		}
		p.postCommandResponse(args, "Subscription created succesfully")
	case "list":
		subsText, err := p.showSubscriptions()
		if err != nil {
			p.postCommandResponse(args, err.Error())
			return &model.CommandResponse{}, nil
		}
		p.postCommandResponse(args, subsText)
	case "unsubscribe":
		if len(parameters) != 1 {
			p.postCommandResponse(args, "You have to pass the channel parameter")
			return &model.CommandResponse{}, nil
		}
		err := p.Unsubscribe(args.TeamId, parameters[0])
		if err != nil {
			p.postCommandResponse(args, err.Error())
			return &model.CommandResponse{}, nil
		}
		p.postCommandResponse(args, "Subscription removed succesfully")
	case "refresh":
		p.RefreshSubscriptions()
	}

	return &model.CommandResponse{}, nil

}

func (p *Plugin) showSubscriptions() (string, error) {
	subscriptions, err := p.GetSubscriptions()
	if err != nil {
		return "", err
	}
	txt := ""
	if len(subscriptions) == 0 {
		txt = "Currently there are no subscriptions"
	} else {
		txt = "### Subscriptions\n"
	}
	for _, sub := range subscriptions {
		txt += fmt.Sprintf("* `%s` - %s\n", strings.Trim(sub.ChannelName, "/"), sub.PanelURL)
	}
	return txt, nil
}
