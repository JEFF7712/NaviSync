# NaviSync

Sync your Spotify playlists to Navidrome automatically.

NaviSync is a native Navidrome plugin that periodically fetches your Spotify playlists and recreates them in Navidrome by matching tracks against your local library.

## Quick Start

### 1. Install Plugin

1. Download `navisync.ndp` from [Releases](https://github.com/JEFF7712/NaviSync/releases)
2. Copy it to your Navidrome plugins directory
3. Enable plugins in `navidrome.toml`:
   ```toml
   [Plugins]
   Enabled = true
   ```
4. Restart Navidrome
5. In Navidrome UI, go to **Settings > Plugins** and assign users to NaviSync

### 2. Create a Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/)
2. Create a new app
3. Add redirect URI: `http://localhost:8888/callback` (for the token helper script)
4. Save your **Client ID** and **Client Secret**

### 3. Get a Refresh Token

Use the included helper script to obtain a Spotify refresh token:

```bash
node scripts/get_token.js <CLIENT_ID> <CLIENT_SECRET>
```

Follow the browser prompt to authorize, then copy the refresh token.

### 4. Configure Plugin

In Navidrome UI, go to **Settings > Plugins > NaviSync**:

- **Spotify Client ID**: Your app's Client ID
- **Spotify Client Secret**: Your app's Client Secret
- **Spotify Refresh Token**: The token from step 3
- **Sync Interval**: Cron expression (default: `0 */6 * * *` = every 6 hours)
- **Playlists Filter**: Leave empty to sync all, or comma-separated playlist names

## How It Works

1. On each scheduled sync, NaviSync fetches all users assigned to the plugin
2. For each user, it retrieves their Spotify refresh token from Navidrome's KVStore (falling back to the global config)
3. It authenticates with Spotify and fetches the user's playlists
4. For each playlist (optionally filtered), it searches Navidrome's library via the Subsonic API to match tracks by artist and title
5. Matched tracks are assembled into a Navidrome playlist (created or updated)

## Features

- Automatic scheduled sync via Navidrome's scheduler
- Per-user Spotify tokens stored in Navidrome's KVStore
- Track matching: exact case-insensitive match, with fuzzy contains-based fallback
- Playlist filtering by name
- Automatic Spotify token refresh with rotation support
- Manual sync and connection test triggers via plugin settings

## Troubleshooting

**Plugin won't load?**
- Ensure you're running Navidrome with plugin support enabled
- Check logs: set `Plugins.LogLevel = "debug"` in your Navidrome config
- Verify the `.ndp` file is in the plugins directory

**No users found?**
- Users must be explicitly assigned to the plugin in **Settings > Plugins > NaviSync**

**Tracks not matching?**
- NaviSync matches by artist name and track title against your Navidrome library
- Ensure your library has accurate metadata (tools like MusicBrainz Picard help)
- Check debug logs for "Unmatched" entries to see which tracks couldn't be found

**Token errors?**
- Verify your Spotify Client ID and Client Secret are correct
- Re-run the token helper script to get a fresh refresh token
- Check that `api.spotify.com` and `accounts.spotify.com` are allowed in the plugin's HTTP permissions

## Development

### Requirements

- Go 1.23+
- TinyGo 0.34+

### Build

```bash
# Quick build
make build

# Or use the helper script
./build.sh
```

This produces `navisync.ndp` (a zip containing `manifest.json` + `plugin.wasm`).

### Project Structure

```
main.go              - Plugin registration (lifecycle + scheduler)
sync/sync.go         - Sync orchestration (OnInit, OnCallback, PerformSync)
navidrome/navidrome.go - Navidrome host function wrappers (Subsonic API, KVStore, Users)
spotify/client.go    - Spotify API client (playlists, tracks, token refresh)
manifest.json        - Plugin manifest with permissions and config schema
```

## License

MIT
