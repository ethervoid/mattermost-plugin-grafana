package main

import "time"

// RunScheduler blabla
func (p *Plugin) RunScheduler() {
	go func() {
		<-time.NewTimer(10 * time.Second).C
		p.RefreshSubscriptions()
		p.RunScheduler()
	}()
}
