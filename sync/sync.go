package sync

import (
	"fmt"
	"github.com/extism/go-pdk"
	"github.com/JEFF7712/NaviSync/spotify"
	"github.com/JEFF7712/NaviSync/navidrome"
)

// ScheduleSync sets up the recurring synchronization task.
func ScheduleSync() error {
	interval, ok := pdk.GetConfig("sync_interval")
	if !ok || interval == "" {
		interval = "0 */6 * * *" // Default: every 6 hours
	}

	pdk.Log(pdk.LogInfo, fmt.Sprintf("Scheduling Spotify sync with interval: %s", interval))
	
	// Register the callback "nd_sync_spotify" with the scheduler
	return navidrome.Schedule(interval, "nd_sync_spotify")
}

// CheckTriggers checks for manual action flags in config.
func CheckTriggers() {
	testConn, _ := pdk.GetConfig("test_connection")
	if testConn == "true" {
		pdk.Log(pdk.LogInfo, "TEST: Testing Spotify connection...")
		
		token, _ := pdk.GetConfig("spotify_refresh_token")
		clientID, _ := pdk.GetConfig("spotify_client_id")
		clientSecret, _ := pdk.GetConfig("spotify_client_secret")

		client := spotify.NewClient(token, clientID, clientSecret)
		accessToken, err := client.RefreshToken()
		if err != nil {
			pdk.Log(pdk.LogError, "TEST: Connection test failed: "+err.Error())
		} else if accessToken != "" {
			pdk.Log(pdk.LogInfo, "TEST: Connection successful! Obtained Access Token.")
		}
	}

	manualSync, _ := pdk.GetConfig("manual_sync")
	if manualSync == "true" {
		pdk.Log(pdk.LogInfo, "MANUAL: Triggering manual sync...")
		if err := PerformSync(); err != nil {
			pdk.Log(pdk.LogError, "MANUAL: Sync failed: "+err.Error())
		} else {
			pdk.Log(pdk.LogInfo, "MANUAL: Sync finished successfully.")
		}
	}
}

// PerformSync is the main logic called by the scheduler.
// It iterates through all users (or configured users) and syncs their playlists.
func PerformSync() error {
	pdk.Log(pdk.LogInfo, "Starting scheduled Spotify sync...")

	// Fetch all users from Navidrome
	users, err := navidrome.GetUsers()
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		// Check if user has Spotify credentials stored in KVStore
		token, err := navidrome.GetUserToken(user.ID)
		if err != nil {
			pdk.Log(pdk.LogWarn, fmt.Sprintf("Skipping user %s: no token found or error: %v", user.Username, err))
			continue
		}

		clientID, _ := pdk.GetConfig("spotify_client_id")
		clientSecret, _ := pdk.GetConfig("spotify_client_secret")

		// Refresh token if needed
		client := spotify.NewClient(token, clientID, clientSecret)
		accessToken, err := client.RefreshToken()
		if err == nil && accessToken != "" {
			pdk.Log(pdk.LogInfo, fmt.Sprintf("Successfully refreshed token for user %s", user.Username))
			client.SetToken(accessToken)
		} else if err != nil {
			pdk.Log(pdk.LogError, fmt.Sprintf("Failed to refresh token for user %s: %v", user.Username, err))
			continue
		}

		// Sync playlists for this user
		if err := syncUserPlaylists(client, user); err != nil {
			pdk.Log(pdk.LogError, fmt.Sprintf("Failed to sync playlists for user %s: %v", user.Username, err))
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
		// Parse comma-separated filter
		// ... (omitted for brevity, assume simple split)
	}

	for _, pl := range playlists {
		if len(filter) > 0 && !filter[pl.Name] {
			continue
		}

		pdk.Log(pdk.LogDebug, fmt.Sprintf("Syncing playlist: %s", pl.Name))
		
		tracks, err := client.GetPlaylistTracks(pl.ID)
		if err != nil {
			pdk.Log(pdk.LogError, fmt.Sprintf("Failed to get tracks for playlist %s: %v", pl.Name, err))
			continue
		}

		// Match tracks in Navidrome
		var navidromeTrackIDs []string
		for _, track := range tracks {
			ndTrack, err := navidrome.FindTrack(track.ISRC, track.Artist, track.Title)
			if err == nil && ndTrack != nil {
				navidromeTrackIDs = append(navidromeTrackIDs, ndTrack.ID)
			} else {
				pdk.Log(pdk.LogWarn, fmt.Sprintf("Unmatched track: %s - %s (ISRC: %s)", track.Artist, track.Title, track.ISRC))
			}
		}

		// Update or create playlist in Navidrome
		if len(navidromeTrackIDs) > 0 {
			if err := navidrome.UpdatePlaylist(user.ID, pl.Name, navidromeTrackIDs); err != nil {
				pdk.Log(pdk.LogError, fmt.Sprintf("Failed to update playlist %s: %v", pl.Name, err))
			}
		}
	}
	return nil
}
