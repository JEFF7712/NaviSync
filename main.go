package main

import (
	"github.com/extism/go-pdk"
	"github.com/JEFF7712/NaviSync/sync"
)

//export nd_on_init
func nd_on_init() int32 {
	// Initialize the plugin
	// This function is called when the plugin is loaded by Navidrome.
	// We can use this to register scheduled tasks.
	
	// Register the sync task with the scheduler
	// The schedule is retrieved from the configuration "sync_interval"
	// Default is "0 */6 * * *" (every 6 hours)
	err := sync.ScheduleSync()
	if err != nil {
		pdk.Log(pdk.LogError, "Failed to schedule sync: "+err.Error())
		return 1
	}

	// Check for manual triggers on startup/reload
	sync.CheckTriggers()

	pdk.Log(pdk.LogInfo, "Navidrome Spotify Sync plugin initialized")
	return 0
}

func main() {}
