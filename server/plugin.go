package main

import (
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// NewPlugin returns an instance of a Plugin
func NewPlugin() *Plugin {
	return &Plugin{}
}

// OnActivate blabla
func (p *Plugin) OnActivate() error {
	p.API.RegisterCommand(getCommand())
	p.RunScheduler()
	return nil
}
