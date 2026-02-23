//go:build wasip1

package main

import (
	"github.com/navidrome/navidrome/plugins/pdk/go/lifecycle"
	"github.com/navidrome/navidrome/plugins/pdk/go/scheduler"

	navisync "github.com/JEFF7712/NaviSync/sync"
)

// NaviSyncPlugin implements the Navidrome lifecycle and scheduler interfaces.
type NaviSyncPlugin struct{}

func (p *NaviSyncPlugin) OnInit() error {
	return navisync.OnInit()
}

func (p *NaviSyncPlugin) OnCallback(req scheduler.SchedulerCallbackRequest) error {
	return navisync.OnCallback(req)
}

func main() {}

func init() {
	plugin := &NaviSyncPlugin{}
	lifecycle.Register(plugin)
	scheduler.Register(plugin)
}
