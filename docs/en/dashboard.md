# Dashboard

The ALTCHA Dashboard is a web UI that visualizes `/challenge` and `/verify` request statistics. It runs as a separate binary (`/dashboard`) and shares the same PostgreSQL database with the API server.

## Architecture

```
[Browser] ──→ [Dashboard :9000] ──→ [PostgreSQL]
                                          ↑
[Browser] ──→ [API Server :3000] ────────┘ (records events)
```

- **API Server**: When `POSTGRES_URL` is set, records `/challenge` and `/verify` requests to the PostgreSQL `events` table
- **Dashboard**: Queries the same PostgreSQL to visualize statistics

## Quick Start

```bash
docker compose up --build
```

- API: http://localhost:3000
- Dashboard: http://localhost:9000 (default credentials: `admin`/`admin`)

## Environment Variables

### Analytics (API Server)

| Variable | Required | Default | Description |
|---|---|---|---|
| POSTGRES_URL | | | PostgreSQL connection URL. Enables analytics when set |
| GEOIP_DB | | | Path to GeoLite2-Country.mmdb. Enables location statistics |

### Dashboard

| Variable | Required | Default | Description |
|---|---|---|---|
| POSTGRES_URL | Yes | | PostgreSQL connection URL |
| DASHBOARD_PORT | | `9000` | Dashboard server port |
| AUTH_PROVIDER | Yes | | Authentication method: `basic` or `keycloak` |

### Basic Auth

| Variable | Description |
|---|---|
| AUTH_USERNAME | Basic auth username |
| AUTH_PASSWORD | Basic auth password |

### Keycloak (OIDC) Auth

| Variable | Required | Description |
|---|---|---|
| AUTH_ISSUER | Yes | Keycloak realm URL (e.g., `https://keycloak.example.com/realms/myrealm`) |
| AUTH_CLIENT_ID | Yes | OIDC client ID |
| AUTH_CLIENT_SECRET | | Client secret (not needed for public clients) |
| AUTH_PKCE | | Enable PKCE (default: `true`) |
| AUTH_AUTHORIZATION_ENDPOINT | | Authorization endpoint override |
| AUTH_TOKEN_ENDPOINT | | Token endpoint override |
| AUTH_END_SESSION_ENDPOINT | | Logout endpoint override |
| AUTH_JWKS_URI | | JWKS URI override |

### Access Control

| Variable | Description |
|---|---|
| AUTH_ALLOWED_USERS | Allowed usernames (comma-separated) |
| AUTH_ALLOWED_GROUPS | Allowed groups (comma-separated) |
| AUTH_ALLOWED_ROLES | Allowed roles (comma-separated) |

When none are set, all authenticated users are allowed. When any are set, they are evaluated with OR logic.

## GeoIP Setup

Country-level statistics require the MaxMind GeoLite2-Country database. This is optional — when not configured, the Locations section is hidden and all other features work normally.

1. Register a free account at [MaxMind](https://dev.maxmind.com/geoip/geolite2-free-geolocation-data) (license key required for download)
2. Download GeoLite2-Country.mmdb
3. Set the `GEOIP_DB` environment variable to the file path

```env
GEOIP_DB=GeoLite2-Country.mmdb
```

> The mmdb file (~6MB) is a binary and may have redistribution restrictions under MaxMind's license, so it is not committed to git (`*.mmdb` is in `.gitignore`). In Docker/K8s environments, provide it via volume mount.

## API Endpoints

The dashboard provides these internal APIs (authentication required):

- `GET /api/summary?from=YYYY-MM-DD&to=YYYY-MM-DD` — KPI summary
- `GET /api/timeseries?from=YYYY-MM-DD&to=YYYY-MM-DD` — Daily trends
- `GET /api/locations?from=YYYY-MM-DD&to=YYYY-MM-DD` — Country/continent statistics

## Kubernetes Deployment

```yaml
# Run /dashboard via command in dashboard-deploy.yaml
command: ["/dashboard"]
```

The dashboard should only be accessible from the internal network using the `nginx-internal` Ingress class.

## Dashboard Features

- **KPI Cards**: Challenges, Verified, Failed, Avg Latency, 4XX Errors, 5XX Errors, Total Requests
- **Trend Chart**: Mixed chart with daily request counts (bar) and average latency (line)
- **Location Stats**: Request distribution by continent/country (when GeoIP is configured)
- **Date Range**: 7 days / 30 days / 90 days / custom selection
