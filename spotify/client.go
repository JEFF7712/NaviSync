package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/extism/go-pdk"
	"time"
)

type Client struct {
	Token        string
	ClientID     string
	ClientSecret string
}

func NewClient(token, clientID, clientSecret string) *Client {
	return &Client{
		Token:        token,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
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

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetPlaylists fetches the current user's playlists with pagination.
func (c *Client) GetPlaylists() ([]Playlist, error) {
	var allPlaylists []Playlist
	nextURL := "https://api.spotify.com/v1/me/playlists?limit=50"

	for nextURL != "" {
		req := pdk.NewHTTPRequest("GET", nextURL)
		req.SetHeader("Authorization", "Bearer "+c.Token)
		
		res := req.Send()
		if res.Status() == 429 {
			pdk.Log(pdk.LogWarn, "Spotify Rate Limit (429) hit. Waiting...")
			time.Sleep(5 * time.Second) // Simple backoff
			continue
		}
		if res.Status() != 200 {
			return nil, fmt.Errorf("failed to get playlists: status %d. Body: %s", res.Status(), string(res.Body()))
		}
		
		var response PlaylistsResponse
		if err := json.Unmarshal(res.Body(), &response); err != nil {
			return nil, err
		}
		
		allPlaylists = append(allPlaylists, response.Items...)
		nextURL = response.Next
	}
	
	return allPlaylists, nil
}

// GetPlaylistTracks fetches tracks from a playlist with pagination.
func (c *Client) GetPlaylistTracks(playlistID string) ([]Track, error) {
	var allTracks []Track
	nextURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?limit=100", playlistID)

	for nextURL != "" {
		req := pdk.NewHTTPRequest("GET", nextURL)
		req.SetHeader("Authorization", "Bearer "+c.Token)
		
		res := req.Send()
		if res.Status() == 429 {
			pdk.Log(pdk.LogWarn, "Spotify Rate Limit (429) hit. Waiting...")
			time.Sleep(5 * time.Second)
			continue
		}
		if res.Status() != 200 {
			return nil, fmt.Errorf("failed to get tracks for playlist %s: status %d", playlistID, res.Status())
		}
		
		var response PlaylistTracksResponse
		if err := json.Unmarshal(res.Body(), &response); err != nil {
			return nil, err
		}
		
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
			allTracks = append(allTracks, track)
		}
		nextURL = response.Next
	}
	
	return allTracks, nil
}

// RefreshToken refreshes the access token using the refresh token.
func (c *Client) RefreshToken() (string, error) {
	if c.ClientID == "" || c.ClientSecret == "" {
		return "", fmt.Errorf("spotify_client_id or spotify_client_secret not configured")
	}

	url := "https://accounts.spotify.com/api/token"
	req := pdk.NewHTTPRequest("POST", url)
	
	auth := base64.StdEncoding.EncodeToString([]byte(c.ClientID + ":" + c.ClientSecret))
	req.SetHeader("Authorization", "Basic "+auth)
	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	
	body := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", c.Token)
	req.SetBody([]byte(body))
	
	res := req.Send()
	if res.Status() == 429 {
		return "", fmt.Errorf("Spotify Rate Limit (429) on token refresh. Please wait.")
	}
	if res.Status() != 200 {
		return "", fmt.Errorf("failed to refresh token: status %d. Body: %s", res.Status(), string(res.Body()))
	}
	
	var response TokenResponse
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return "", err
	}
	
	return response.AccessToken, nil
}
