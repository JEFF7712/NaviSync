package navidrome

import (
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

// --- KVStore ---

// GetUserToken retrieves the Spotify OAuth token.
// Priority:
// 1. KVStore (not fully implemented yet, skipping)
// 2. Global Config (spotify_refresh_token)
func GetUserToken(userID string) (string, error) {
	// 1. Try Config (Simplest for single-user/admin setup)
	token, _ := pdk.GetConfig("spotify_refresh_token")
	if token != "" {
		return token, nil
	}

	// 2. Try KVStore (Placeholder implementation)
	// In the future, we will fetch this from the host.
	// For now, we rely on the config.
	
	return "", fmt.Errorf("no refresh token found in config 'spotify_refresh_token'")
}

// SetUserToken saves the Spotify OAuth token.
func SetUserToken(userID string, token string) error {
	// Placeholder: We can't update the config from here.
	// Real implementation would use KVStore.
	pdk.Log(pdk.LogDebug, "Mock: Saving new token for user "+userID)
	return nil
}

// --- Scheduler ---

// Schedule registers a callback to be called at the specified interval.
func Schedule(cron string, callback string) error {
	// Placeholder for scheduler host function
	pdk.Log(pdk.LogInfo, "Scheduling " + callback + " at " + cron)
	return nil
}

// --- Users ---

// GetUsers returns all users in Navidrome.
func GetUsers() ([]User, error) {
	// Placeholder for users host function
	// In a real plugin, we'd call `navidrome_get_users`
	return []User{{ID: "admin", Username: "admin"}}, nil
}

// --- Subsonic / Library ---

// FindTrack attempts to find a track in Navidrome by ISRC or fuzzy match.
func FindTrack(isrc, artist, title string) (*Track, error) {
	// Placeholder
	// 1. Search by ISRC if available
	// 2. Search by "artist title"
	return &Track{ID: "123", Title: title, Artist: artist}, nil
}

// UpdatePlaylist creates or updates a playlist for a user.
func UpdatePlaylist(userID, name string, trackIDs []string) error {
	// Placeholder
	// 1. Check if playlist exists
	// 2. If not, create it
	// 3. Update tracks
	pdk.Log(pdk.LogInfo, fmt.Sprintf("Updated playlist %s for user %s with %d tracks", name, userID, len(trackIDs)))
	return nil
}
