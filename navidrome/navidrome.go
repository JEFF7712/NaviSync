package navidrome

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/extism/go-pdk"
	"github.com/navidrome/navidrome/plugins/pdk/go/host"
)

// User is the Navidrome user type from the PDK host.
// Fields: UserName, Name, IsAdmin.
type User = host.User

// Track represents a matched track in Navidrome.
type Track struct {
	ID     string
	Title  string
	Artist string
	Album  string
}

// --- Subsonic API response types ---

type SubsonicResponse struct {
	Subsonic SubsonicBody `json:"subsonic-response"`
}

type SubsonicBody struct {
	Status        string         `json:"status"`
	Error         *SubsonicError `json:"error,omitempty"`
	SearchResult3 *SearchResult3 `json:"searchResult3,omitempty"`
	Playlists     *PlaylistList  `json:"playlists,omitempty"`
	Playlist      *Playlist      `json:"playlist,omitempty"`
}

type SubsonicError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SearchResult3 struct {
	Song []Song `json:"song"`
}

type Song struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

type PlaylistList struct {
	Playlist []Playlist `json:"playlist"`
}

type Playlist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	SongCount int    `json:"songCount"`
}

// --- Host function wrappers ---

// GetUsers returns all users assigned to this plugin.
func GetUsers() ([]User, error) {
	return host.UsersGetUsers()
}

// GetUserToken retrieves the Spotify refresh token for a user.
// Checks KVStore first, then falls back to global config.
func GetUserToken(userName string) (string, error) {
	key := "spotify_token:" + userName
	data, exists, err := host.KVStoreGet(key)
	if err == nil && exists && len(data) > 0 {
		return string(data), nil
	}

	// Fall back to global config
	token, ok := pdk.GetConfig("spotify_refresh_token")
	if ok && token != "" {
		return token, nil
	}

	return "", fmt.Errorf("no refresh token found for user %s", userName)
}

// SetUserToken persists a Spotify refresh token for a user in KVStore.
func SetUserToken(userName string, token string) error {
	key := "spotify_token:" + userName
	return host.KVStoreSet(key, []byte(token))
}

// makeSubsonicRequest calls Navidrome's internal Subsonic API and parses the response.
func makeSubsonicRequest(endpoint, userName string, params *url.Values) (*SubsonicResponse, error) {
	uri := fmt.Sprintf("/rest/%s?u=%s", endpoint, userName)
	if params != nil {
		encoded := params.Encode()
		if encoded != "" {
			uri += "&" + encoded
		}
	}

	respStr, err := host.SubsonicAPICall(uri)
	if err != nil {
		return nil, fmt.Errorf("SubsonicAPICall error for %s: %w", endpoint, err)
	}

	var resp SubsonicResponse
	if err := json.Unmarshal([]byte(respStr), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse %s response: %w", endpoint, err)
	}

	if resp.Subsonic.Status != "ok" {
		if resp.Subsonic.Error != nil {
			return nil, fmt.Errorf("subsonic error (%d): %s", resp.Subsonic.Error.Code, resp.Subsonic.Error.Message)
		}
		return nil, fmt.Errorf("subsonic status not ok for %s", endpoint)
	}

	return &resp, nil
}

// FindTrack searches for a track in Navidrome by artist and title.
func FindTrack(userName, artist, title string) (*Track, error) {
	params := url.Values{
		"query":       []string{title},
		"songCount":   []string{"20"},
		"artistCount": []string{"0"},
		"albumCount":  []string{"0"},
	}

	resp, err := makeSubsonicRequest("search3", userName, &params)
	if err != nil {
		return nil, err
	}

	if resp.Subsonic.SearchResult3 == nil || len(resp.Subsonic.SearchResult3.Song) == 0 {
		return nil, nil
	}

	lowerTitle := strings.ToLower(title)
	lowerArtist := strings.ToLower(artist)

	// Try exact match first (case-insensitive)
	for _, song := range resp.Subsonic.SearchResult3.Song {
		if strings.ToLower(song.Title) == lowerTitle && strings.ToLower(song.Artist) == lowerArtist {
			return &Track{ID: song.ID, Title: song.Title, Artist: song.Artist, Album: song.Album}, nil
		}
	}

	// Fuzzy fallback: contains-based matching
	for _, song := range resp.Subsonic.SearchResult3.Song {
		songTitle := strings.ToLower(song.Title)
		songArtist := strings.ToLower(song.Artist)
		titleMatch := strings.Contains(songTitle, lowerTitle) || strings.Contains(lowerTitle, songTitle)
		artistMatch := strings.Contains(songArtist, lowerArtist) || strings.Contains(lowerArtist, songArtist)
		if titleMatch && artistMatch {
			return &Track{ID: song.ID, Title: song.Title, Artist: song.Artist, Album: song.Album}, nil
		}
	}

	return nil, nil
}

// UpdatePlaylist creates or updates a playlist for a user.
func UpdatePlaylist(userName, name string, trackIDs []string) error {
	// Find existing playlist
	playlistsResp, err := makeSubsonicRequest("getPlaylists", userName, &url.Values{})
	if err != nil {
		return fmt.Errorf("failed to get playlists: %w", err)
	}

	var existingID string
	if playlistsResp.Subsonic.Playlists != nil {
		for _, pl := range playlistsResp.Subsonic.Playlists.Playlist {
			if pl.Name == name {
				existingID = pl.ID
				break
			}
		}
	}

	// Create or update playlist
	params := url.Values{
		"songId": trackIDs,
	}
	if existingID != "" {
		params.Set("playlistId", existingID)
	} else {
		params.Set("name", name)
	}

	_, err = makeSubsonicRequest("createPlaylist", userName, &params)
	if err != nil {
		return fmt.Errorf("failed to create/update playlist %s: %w", name, err)
	}

	return nil
}
