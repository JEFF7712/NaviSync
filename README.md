# NaviSync

A Navidrome plugin that periodically synchronizes your Spotify playlists to Navidrome.

## Features

- **OAuth 2.0 Integration**: Securely connects to Spotify using standard OAuth flow.
- **Scheduled Sync**: Automatically syncs playlists at a configurable interval (default: every 6 hours).
- **Intelligent Matching**: Matches tracks using ISRC codes (primary) and fuzzy artist/title matching (fallback).
- **Two-way Sync (One-way for now)**: Currently supports syncing Spotify -> Navidrome.
- **Configurable**: Choose specific playlists to sync or sync all.

## Installation

1.  Download the latest `navisync.ndp` release.
2.  Place the `.ndp` file in your Navidrome `plugins` directory (usually alongside your music library or in a dedicated folder).
3.  Ensure your `navidrome.toml` has plugins enabled:

    ```toml
    Plugins.Enable = true
    Plugins.AutoReload = true # Optional, for development
    ```

4.  Restart Navidrome.

## Configuration

After installation, go to **Settings > Plugins > NaviSync** in the Navidrome UI to configure:

-   **Spotify Client ID**: Your Spotify App Client ID.
-   **Spotify Client Secret**: Your Spotify App Client Secret.
-   **Sync Interval**: Cron expression (default: `0 */6 * * *`).
-   **Playlists Filter**: Comma-separated list of playlists to sync (e.g., "My Top Songs, Discover Weekly"). Leave empty to sync all.

### Spotify App Setup

1.  Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/).
2.  Create a new App.
3.  Set the **Redirect URI** to match your Navidrome instance's callback URL (e.g., `https://music.yourdomain.com/auth/spotify/callback`). *Note: The exact callback URL format depends on how Navidrome exposes plugin auth endpoints. Check Navidrome docs.*
4.  Copy the **Client ID** and **Client Secret** to the plugin configuration.

## Development

### Prerequisites

-   Go 1.21+
-   TinyGo 0.30+
-   `make`
-   `zip`

### Build

```bash
make build
```

This produces `navisync.ndp`.

### Debugging

Enable debug logging in Navidrome:

```toml
Plugins.LogLevel = "debug"
```

Check the Navidrome logs for plugin output.
