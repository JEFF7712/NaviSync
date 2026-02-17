package navidrome

import (
	"encoding/json"
	"fmt"

	"github.com/extism/go-pdk"
)

// User represents a Navidrome user.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Track represents a track in Navidrome.
type Track struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	ISRC   string `json:"isrc"`
}

// --- KVStore (Mock) ---

// GetUserToken retrieves the Spotify OAuth token.
func GetUserToken(userID string) (string, error) {
	// Priority: Global Config (spotify_refresh_token)
	token, _ := pdk.GetConfig("spotify_refresh_token")
	if token != "" {
		return token, nil
	}

	return "", fmt.Errorf("no refresh token found for user %s. Please enter it in the plugin settings.", userID)
}

// SetUserToken saves the Spotify OAuth token.
func SetUserToken(userID string, token string) error {
	pdk.Log(pdk.LogDebug, "Mock: SetUserToken not yet linked to real host DB")
	return nil
}

// --- Scheduler (Mock) ---

// Schedule registers a callback.
func Schedule(cron string, callback string) error {
	pdk.Log(pdk.LogInfo, "Mock: Scheduled "+callback+" at "+cron)
	return nil
}

// --- Users (Mock) ---

// GetUsers returns all users.
func GetUsers() ([]User, error) {
	// For now, only return the admin user as a placeholder.
	// Real implementation requires verifying the exact WasmImport name for Navidrome's user service.
	return []User{{ID: "admin", Username: "admin"}}, nil
}

// --- Subsonic / Library (Mock) ---

// FindTrack attempts to find a track in Navidrome.
func FindTrack(isrc, artist, title string) (*Track, error) {
	// This is a placeholder. Real integration will call Navidrome's search.
	pdk.Log(pdk.LogDebug, fmt.Sprintf("Mock: Searching for %s - %s", artist, title))
	return nil, nil // Return nil so it doesn't try to add non-existent tracks
}

// UpdatePlaylist creates or updates a playlist for a user.
func UpdatePlaylist(userID, name string, trackIDs []string) error {
	pdk.Log(pdk.LogInfo, fmt.Sprintf("Mock: Syncing playlist %s for user %s (%d tracks matched)", name, userID, len(trackIDs)))
	return nil
}
