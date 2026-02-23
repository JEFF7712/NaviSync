package sync

import (
	"fmt"
	"strings"

	"github.com/extism/go-pdk"
	"github.com/navidrome/navidrome/plugins/pdk/go/host"
	"github.com/navidrome/navidrome/plugins/pdk/go/scheduler"

	"github.com/JEFF7712/NaviSync/navidrome"
	"github.com/JEFF7712/NaviSync/spotify"
)

const (
	syncPayload  = "sync-spotify"
	syncSchedule = "navisync-schedule"
)

// OnInit is called when the plugin initializes.
// Schedules recurring sync and checks manual triggers.
func OnInit() error {
	interval, ok := pdk.GetConfig("sync_interval")
	if !ok || interval == "" {
		interval = "0 */6 * * *"
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("Scheduling Spotify sync with interval: %s", interval))

	_, err := host.SchedulerScheduleRecurring(interval, syncPayload, syncSchedule)
	if err != nil {
		return fmt.Errorf("failed to schedule sync: %w", err)
	}

	CheckTriggers()

	pdk.Log(pdk.LogInfo, "NaviSync plugin initialized")
	return nil
}

// OnCallback handles scheduler callbacks.
func OnCallback(req scheduler.SchedulerCallbackRequest) error {
	switch req.Payload {
	case syncPayload:
		return PerformSync()
	default:
		pdk.Log(pdk.LogWarn, fmt.Sprintf("Unknown callback payload: %s", req.Payload))
		return nil
	}
}

// CheckTriggers checks for manual action flags in config.
func CheckTriggers() {
	testConn, _ := pdk.GetConfig("test_connection")
	if testConn == "true" {
		pdk.Log(pdk.LogInfo, "Testing Spotify connection...")

		token, _ := pdk.GetConfig("spotify_refresh_token")
		clientID, _ := pdk.GetConfig("spotify_client_id")
		clientSecret, _ := pdk.GetConfig("spotify_client_secret")

		client := spotify.NewClient(token, clientID, clientSecret)
		newRefresh, err := client.RefreshToken()
		if err != nil {
			pdk.Log(pdk.LogError, "Connection test failed: "+err.Error())
		} else {
			pdk.Log(pdk.LogInfo, "Connection test successful!")
			if newRefresh != "" {
				pdk.Log(pdk.LogInfo, "Spotify returned a new refresh token")
			}
		}
	}

	manualSync, _ := pdk.GetConfig("manual_sync")
	if manualSync == "true" {
		pdk.Log(pdk.LogInfo, "Triggering manual sync...")
		if err := PerformSync(); err != nil {
			pdk.Log(pdk.LogError, "Manual sync failed: "+err.Error())
		} else {
			pdk.Log(pdk.LogInfo, "Manual sync finished successfully.")
		}
	}
}

// PerformSync is the main sync logic called by the scheduler.
func PerformSync() error {
	pdk.Log(pdk.LogInfo, "Starting Spotify sync...")

	users, err := navidrome.GetUsers()
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	clientID, _ := pdk.GetConfig("spotify_client_id")
	clientSecret, _ := pdk.GetConfig("spotify_client_secret")

	for _, user := range users {
		token, err := navidrome.GetUserToken(user.UserName)
		if err != nil {
			pdk.Log(pdk.LogWarn, fmt.Sprintf("Skipping user %s: %v", user.UserName, err))
			continue
		}

		client := spotify.NewClient(token, clientID, clientSecret)
		newRefresh, err := client.RefreshToken()
		if err != nil {
			pdk.Log(pdk.LogError, fmt.Sprintf("Failed to authenticate user %s: %v", user.UserName, err))
			continue
		}

		// Persist new refresh token if Spotify returned one (token rotation)
		if newRefresh != "" {
			if err := navidrome.SetUserToken(user.UserName, newRefresh); err != nil {
				pdk.Log(pdk.LogWarn, fmt.Sprintf("Failed to save new refresh token for %s: %v", user.UserName, err))
			}
		}

		pdk.Log(pdk.LogInfo, fmt.Sprintf("Authenticated Spotify for user %s", user.UserName))

		if err := syncUserPlaylists(client, user); err != nil {
			pdk.Log(pdk.LogError, fmt.Sprintf("Failed to sync playlists for user %s: %v", user.UserName, err))
		}
	}

	pdk.Log(pdk.LogInfo, "Spotify sync completed.")
	return nil
}

func syncUserPlaylists(client *spotify.Client, user navidrome.User) error {
	playlists, err := client.GetPlaylists()
	if err != nil {
		return err
	}

	filterStr, _ := pdk.GetConfig("playlists_filter")
	filter := make(map[string]bool)
	if filterStr != "" {
		for _, f := range strings.Split(filterStr, ",") {
			filter[strings.TrimSpace(f)] = true
		}
	}

	for _, pl := range playlists {
		if len(filter) > 0 && !filter[pl.Name] {
			continue
		}

		pdk.Log(pdk.LogInfo, fmt.Sprintf("Syncing playlist [%s] for user %s", pl.Name, user.UserName))

		tracks, err := client.GetPlaylistTracks(pl.ID)
		if err != nil {
			pdk.Log(pdk.LogError, fmt.Sprintf("Failed to get tracks for playlist %s: %v", pl.Name, err))
			continue
		}

		var trackIDs []string
		for _, track := range tracks {
			ndTrack, err := navidrome.FindTrack(user.UserName, track.Artist, track.Title)
			if err == nil && ndTrack != nil {
				trackIDs = append(trackIDs, ndTrack.ID)
			} else {
				pdk.Log(pdk.LogDebug, fmt.Sprintf("Unmatched: %s - %s", track.Artist, track.Title))
			}
		}

		if len(trackIDs) > 0 {
			if err := navidrome.UpdatePlaylist(user.UserName, pl.Name, trackIDs); err != nil {
				pdk.Log(pdk.LogError, fmt.Sprintf("Failed to update playlist %s: %v", pl.Name, err))
			} else {
				pdk.Log(pdk.LogInfo, fmt.Sprintf("Synced %d/%d tracks to [%s]", len(trackIDs), len(tracks), pl.Name))
			}
		} else {
			pdk.Log(pdk.LogWarn, fmt.Sprintf("No tracks matched for playlist [%s]", pl.Name))
		}
	}
	return nil
}
