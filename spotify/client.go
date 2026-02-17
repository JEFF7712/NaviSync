package spotify

import (
	"encoding/json"
	"fmt"
	"github.com/extism/go-pdk"
)

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{Token: token}
}

func (c *Client) SetToken(token string) {
	c.Token = token
}

// Simplified structs for Spotify API response
type Playlist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URI  string `json:"uri"`
}

type PlaylistsResponse struct {
	Items []Playlist `json:"items"`
	Next  string     `json:"next"`
}

type Track struct {
	ID     string `json:"id"`
	Title  string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	ISRC   string `json:"isrc"`
}

type PlaylistTracksResponse struct {
	Items []struct {
		Track struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ExternalIDs struct {
				ISRC string `json:"isrc"`
			} `json:"external_ids"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Name string `json:"name"`
			} `json:"album"`
		} `json:"track"`
	} `json:"items"`
	Next string `json:"next"`
}

// GetPlaylists fetches the current user's playlists.
func (c *Client) GetPlaylists() ([]Playlist, error) {
	// TODO: Handle pagination
	req := pdk.NewHTTPRequest(pdk.MethodGet, "https://api.spotify.com/v1/me/playlists")
	req.SetHeader("Authorization", "Bearer "+c.Token)
	
	res := req.Send()
	if res.Status() != 200 {
		return nil, fmt.Errorf("failed to get playlists: status %d", res.Status())
	}
	
	var response PlaylistsResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return nil, err
	}
	
	return response.Items, nil
}

// GetPlaylistTracks fetches tracks from a playlist.
func (c *Client) GetPlaylistTracks(playlistID string) ([]Track, error) {
	// TODO: Handle pagination
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)
	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("Authorization", "Bearer "+c.Token)
	
	res := req.Send()
	if res.Status() != 200 {
		return nil, fmt.Errorf("failed to get tracks for playlist %s: status %d", playlistID, res.Status())
	}
	
	var response PlaylistTracksResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return nil, err
	}
	
	var tracks []Track
	for _, item := range response.Items {
		track := Track{
			ID:    item.Track.ID,
			Title: item.Track.Name,
			ISRC:  item.Track.ExternalIDs.ISRC,
		}
		if len(item.Track.Artists) > 0 {
			track.Artist = item.Track.Artists[0].Name
		}
		track.Album = item.Track.Album.Name
		tracks = append(tracks, track)
	}
	
	return tracks, nil
}

// RefreshToken refreshes the access token using the refresh token (not implemented fully).
func (c *Client) RefreshToken() (string, error) {
	// This requires client_id/secret and a POST request to accounts.spotify.com
	// For now, we assume the token is valid or refreshed externally.
	// Implementing OAuth flow inside WASM is tricky without host support for secure storage/secrets.
	return c.Token, nil
}
