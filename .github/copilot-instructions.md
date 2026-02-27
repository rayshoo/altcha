# AI coding agent instructions for `altcha`

## Overview

- Purpose: Dockerized ALTCHA challenge/verify microservice using Go + Echo. Provides `/challenge` and `/verify` used by the ALTCHA widget. Optional demo UI.
- Key libs: `github.com/altcha-org/altcha-lib-go`, `github.com/labstack/echo/v4`, `github.com/joho/godotenv`.
- Entrypoint: `cmd/server/main.go`.

## Repo layout

- `cmd/server/main.go`: Entrypoint; loads .env, parses config, starts API and optional demo server.
- `pkg/config/config.go`: Config struct and env-var parsing with defaults.
- `pkg/handler/challenge.go`: `GET /challenge` handler.
- `pkg/handler/verify.go`: `GET /verify` handler with in-memory record cache.
- `pkg/handler/demo.go`: Demo page serving and proxy handlers.
- `pkg/middleware/security.go`: CSP header middleware for demo server.
- `pkg/server/server.go`: Echo server creation and route registration.
- `web/demo/index.html`: Demo UI page.
- `Dockerfile`: multi-stage Go build; copies `.env` and `web/` into image.
- `compose.yaml`: exposes 3000 (API) and 8000 (demo); sets `SECRET` (default is placeholder).
- `Makefile`: `build`, `run`, `dev`, `docker-build`, `docker-up`, `clean`, `lint`.

## Build & run

- Local (Go):
  - `make build && make run`
  - Development: `make dev`
- Docker Compose: `docker compose up --build`
  - Override secret:
    - PowerShell: `$env:SECRET = "<long-random>"; docker compose up --build`
    - Unix: `SECRET="<long-random>" docker compose up --build`

## Configuration (env)

- `SECRET` (required): HMAC key for ALTCHA. Default `$ecret.key` is unsafe; code logs a warning if used.
- `PORT`: API port (default 3000).
- `EXPIREMINUTES`: challenge expiry minutes (default 10).
- `COMPLEXITY`: PoW complexity / max number for difficulty (default 1000000).
- `MAXRECORDS`: in-memory single-use token cache size (default 1000).
- `CORS_ORIGIN`: comma-separated allowed origins; defaults to `*` if unset.
- `RATE_LIMIT`: requests per second per IP (0 or unset = unlimited).
- `STORE`: token store backend: `memory` (default), `sqlite`, `redis`.
- `SQLITE_PATH`: SQLite file path (default `data/altcha.db`, used when STORE=sqlite).
- `REDIS_URL`: Redis connection URL (default `redis://localhost:6379`, used when STORE=redis).
- `REDIS_CLUSTER`: set `true` for cluster mode (ElastiCache, Valkey); also auto-detected when REDIS_URL contains commas.
- `LOG_LEVEL`: `info` (API logs only, default) or `debug` (API + demo logs).
- `DEMO`: when `true`, serve demo on 8000 with CSP middleware.
- `.env` is loaded by `godotenv` at runtime; Dockerfile also copies `.env` into image.

## API contracts (keep stable)

- `GET /` → `204 No Content` (liveness).
- `GET /health` → `200 OK` JSON with status, version, go runtime.
- `GET /challenge` → `200 OK` JSON from `altcha.CreateChallenge()`.
- `GET /verify?altcha=<payload>` → `202 Accepted` on success, `417 Expectation Failed` on invalid or reused token.
- Reuse prevention uses an in-memory `recordCache` (size = `MAXRECORDS`); cache clears on restart/scaling.
- CORS defaults to `*`; configurable via `CORS_ORIGIN`. Demo uses strict CSP.

## Patterns & conventions

- Standard Go project layout: `cmd/`, `pkg/`.
- Echo framework for HTTP; minimal error handling by design (status-only API).
- Keep endpoints and status codes as-is to preserve client integrations and docs.
- When adding env vars or endpoints, update `README.md`, `.env.example`, and this file.

## CI/CD

- GitHub Actions: `.github/workflows/docker-publish.yml` builds multi-arch images (amd64/arm64) with Buildx/QEMU.
- Publishes to GHCR `ghcr.io/<owner>/<repo>` on pushes to `main` and version tags `v*.*.*`.
- Uses `docker/metadata-action` for tags/labels; caches via GHA cache.

## Common tasks (examples)

- Test verify manually:
  - PowerShell: `curl "http://localhost:3000/verify?altcha=$([uri]::EscapeDataString($payload))" -Method GET -UseBasicParsing`
  - Unix: `curl -G --data-urlencode "altcha=$payload" http://localhost:3000/verify -i`
- Enable demo: set `DEMO=true` and open `http://localhost:8000`.

## Gotchas

- Do not ship with default `SECRET`.
- In-memory token cache is not shared across replicas; use a shared store if you scale (out of scope here).
- The demo proxy posts to `/test` and calls API locally at `http://localhost:3000`.
