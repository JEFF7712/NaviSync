# NaviSync

Sync your Spotify playlists to Navidrome automatically.

## Quick Start

### 1. Install Plugin

1. Download `navisync.ndp` from [Releases](https://github.com/JEFF7712/NaviSync/releases)
2. Copy it to your Navidrome plugins folder
3. Add to `navidrome.toml`:
   ```toml
   [Plugins]
   Enabled = true
   ```
4. Restart Navidrome

### 2. Configure Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/)
2. Create a new app
3. Add redirect URI: `http://localhost:8888/callback` (for local OAuth flow)
4. Save your **Client ID** and **Client Secret**

### 3. Configure Plugin

In Navidrome UI → **Settings > Plugins > NaviSync**:

- **Spotify Client ID**: Your app's Client ID
- **Spotify Client Secret**: Your app's Client Secret  
- **Spotify Redirect URI**: `http://localhost:8888/callback`
- **Sync Interval**: `0 */6 * * *` (every 6 hours, or customize)
- **Playlists Filter**: Leave empty to sync all, or comma-separated names

### 4. Authenticate

The first sync will prompt you to authorize via Spotify OAuth. Follow the browser prompt to grant access.

## How It Works

- **Scheduled sync** runs on your configured interval (default: every 6 hours)
- **Smart matching**: Tracks matched by ISRC codes, with fuzzy fallback for artist/title
- **Persistent tokens**: OAuth tokens stored securely in Navidrome's KVStore
- **Per-user**: Each Navidrome user can connect their own Spotify account

## Features

✅ OAuth 2.0 authentication  
✅ Automatic token refresh  
✅ ISRC-based track matching  
✅ Configurable sync intervals  
✅ Filter specific playlists  
✅ Per-user configuration  

## Troubleshooting

**Plugin won't enable?**
- Check Navidrome logs: `Plugins.LogLevel = "debug"` in config
- Verify all users are assigned to the plugin in Settings

**Tracks not matching?**
- Ensure your Navidrome library has accurate metadata (use MusicBrainz Picard)
- Check logs for unmatched tracks

**OAuth fails?**
- Verify redirect URI matches exactly in Spotify app settings
- Ensure `http://localhost:8888/callback` is accessible from where you run the OAuth flow

## Development

Build from source:
```bash
make build  # Requires TinyGo 0.30+
```

## License

MIT
