# Concert Manager

Backend service for the Beacon concert management app. Provides REST API for concert discovery, artist tracking, and event recommendations.

## Building and Running

### Prerequisites

- Go 1.21+
- Environment configuration file (`env.sh`)

### Makefile Targets

**Build both TUI and server:**
```bash
make all
```

**Build server only:**
```bash
make server
```

**Build TUI only:**
```bash
make tui
```

**Run server (builds first if needed):**
```bash
make server
make runserver
```

**Run TUI (builds first if needed):**
```bash
make tui
make runtui
```

**Pass arguments to executables:**
```bash
make runserver ARGS="--port 8080"
make runtui ARGS="--verbose"
```

### Build Artifacts

All binaries are built to the `build/` directory:
- `build/cm-server` - HTTP REST API server
- `build/cm-tui` - Terminal UI client

### Environment Configuration

Create an `env.sh` file in the concert-manager directory with required environment variables:

```bash
export CM_PROJ_ID=""
export CM_LOG_LEVEL="DEBUG"
export CM_TICKETMASTER_API_KEY=""
export CM_SPOTIFY_AUTH_TOKEN=""
export CM_SPOTIFY_REFRESH_TOKEN=""
export CM_LASTFM_API_KEY=""
export CM_ALERT_EMAIL=""
export CM_GMAIL_USER=""
export CM_GMAIL_PASSWORD=""

```

The `env.sh` file is automatically sourced when using `make runserver` or `make runtui`.

## Deployment

For deployment and management scripts, see [scripts/README.md](scripts/README.md).
