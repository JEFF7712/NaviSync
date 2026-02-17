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

// Result wrapper for host calls
type hostResult struct {
	Exists bool            `json:"exists"`
	Value  json.RawMessage `json:"value"`
}

// --- KVStore ---

// GetUserToken retrieves the Spotify OAuth token from KVStore.
func GetUserToken(userID string) (string, error) {
	// First check config for manual override
	token, _ := pdk.GetConfig("spotify_refresh_token")
	if token != "" {
		return token, nil
	}

	// Try Host KVStore
	key := "spotify_token:" + userID
	res, err := pdk.HostFunction("kvstore", "get", []byte(key))
	if err != nil {
		return "", fmt.Errorf("kvstore error: %w", err)
	}

	var result hostResult
	if err := json.Unmarshal(res, &result); err != nil {
		return "", err
	}

	if !result.Exists {
		return "", fmt.Errorf("no token found for user %s", userID)
	}

	return string(result.Value), nil
}

// SetUserToken saves the Spotify OAuth token.
func SetUserToken(userID string, token string) error {
	key := "spotify_token:" + userID
	_, err := pdk.HostFunction("kvstore", "set", []byte(fmt.Sprintf(`{"key":"%s","value":"%s"}`, key, token)))
	return err
}

// --- Scheduler ---

// Schedule registers a callback to be called at the specified interval.
func Schedule(cron string, callback string) error {
	payload := fmt.Sprintf(`{"cron":"%s","callback":"%s"}`, cron, callback)
	_, err := pdk.HostFunction("scheduler", "schedule_recurring", []byte(payload))
	return err
}

// --- Users ---

// GetUsers returns all users in Navidrome.
func GetUsers() ([]User, error) {
	res, err := pdk.HostFunction("navidrome", "get_users", nil)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(res, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// --- Subsonic / Library ---

// FindTrack attempts to find a track in Navidrome by ISRC or fuzzy match.
func FindTrack(isrc, artist, title string) (*Track, error) {
	query := fmt.Sprintf(`{"isrc":"%s","artist":"%s","title":"%s"}`, isrc, artist, title)
	res, err := pdk.HostFunction("subsonic", "find_track", []byte(query))
	if err != nil {
		return nil, err
	}

	var track Track
	if err := json.Unmarshal(res, &track); err != nil {
		if string(res) == "null" {
			return nil, nil
		}
		return nil, err
	}
	return &track, nil
}

// UpdatePlaylist creates or updates a playlist for a user.
func UpdatePlaylist(userID, name string, trackIDs []string) error {
	tracksJSON, _ := json.Marshal(trackIDs)
	payload := fmt.Sprintf(`{"userId":"%s","name":"%s","trackIds":%s}`, userID, name, string(tracksJSON))
	_, err := pdk.HostFunction("subsonic", "update_playlist", []byte(payload))
	return err
}
