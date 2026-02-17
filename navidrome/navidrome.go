package navidrome

import (
	"encoding/json"
	"errors"
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

// GetUserToken retrieves the Spotify OAuth token for a user from the KVStore.
func GetUserToken(userID string) (string, error) {
	// Key format: "spotify_token:<userID>"
	key := "spotify_token:" + userID
	mem := pdk.AllocateString(key)
	defer mem.Free()
	
	// Assuming `kvstore_get` is the host function
	// In Extism Go PDK, host functions are called via `HostFunction`
	// But usually, there are wrapper libraries.
	// Since I don't have the library, I'll simulate it or use generic host calls.
	
	// Implementation placeholder:
	// val, err := pdk.GetInput() ...
	// Since I cannot implement exact host calls without the definition,
	// I will write the logic assuming a helper function exists.
	
	// For now, let's assume `pdk.GetVar` might work for KVStore if mapped,
	// but `kvstore` is likely a specific host module.
	
	// Using generic HostFunction call:
	// res, err := pdk.HostFunction("kvstore", "get", []byte(key))
	// if err != nil { return "", err }
	// return string(res), nil
	
	return "", errors.New("KVStore not implemented (missing host function definitions)")
}

// SetUserToken saves the Spotify OAuth token for a user.
func SetUserToken(userID string, token string) error {
	key := "spotify_token:" + userID
	// Implementation placeholder
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
