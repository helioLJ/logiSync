# LogiSync

Minimal logistics tracking pipeline: Go API + Redis Streams worker + Postgres, plus a mock carrier portal and Playwright scraper.

## Features

- Async job processing via Redis Streams.
- Postgres-backed job status, results, and artifact metadata.
- Providers: `dummy`, `mock_portal_scrape` (Playwright against mock portal).
- Local artifacts stored with S3-ready keys.

## Architecture

```mermaid
flowchart LR
  Client[Client / CLI] -->|POST /v1/tracking/jobs| API[Go API]
  Client -->|GET /v1/jobs/{id}| API
  Client -->|GET /v1/tracking/results/{id}| API

  API -->|insert job| PG[(Postgres)]
  API -->|XADD tracking:jobs| Redis[(Redis Streams)]

  Worker[Go Worker] -->|XREADGROUP| Redis
  Worker -->|mark RUNNING/DONE/FAILED| PG
  Worker -->|insert results| PG
  Worker -->|insert artifacts metadata| PG

  Worker -->|provider=dummy| Dummy[Dummy Provider]
  Worker -->|provider=mock_portal_scrape| Scraper[Playwright Scraper]

  Scraper -->|GET /track + /api/track/{code}| Portal[Mock Portal Service]
  Scraper -->|capture response/HTML/screenshot| ArtifactStore[(./artifacts)]

  Dummy -->|payload.json| ArtifactStore
```

Notes:
- Redis stream: `tracking:jobs` with consumer group `tracking-workers`.
- Job lifecycle: `PENDING` → `RUNNING` → `DONE` or `FAILED`.
- Artifacts are saved locally and referenced by S3-ready keys in Postgres.

## Quickstart

```bash
docker compose up -d
make migrate
make api
```

In another terminal:

```bash
make worker
```

In another terminal (mock portal UI):

```bash
make mock-portal
```

Create a job:

```bash
curl -X POST http://localhost:8080/v1/tracking/jobs \
  -H 'Content-Type: application/json' \
  -d '{"provider":"mock_portal_scrape","tracking_code":"AA123"}'
```

Check job status:

```bash
curl http://localhost:8080/v1/jobs/<job_id>
```

Fetch tracking result:

```bash
curl http://localhost:8080/v1/tracking/results/<job_id>
```

Mock portal UI:

- http://localhost:8090/track

## Playwright setup

```bash
make playwright-install
```

If you want to see the browser window:

```bash
PLAYWRIGHT_HEADLESS=false make worker
```

## Configuration

The app loads environment variables from `.env` if present.

Environment variables:

- `DB_URL` (default `postgres://postgres:postgres@localhost:5432/logisync?sslmode=disable`)
- `REDIS_ADDR` (default `localhost:6379`)
- `REDIS_STREAM` (default `tracking:jobs`)
- `REDIS_GROUP` (default `tracking-workers`)
- `REDIS_CONSUMER` (default `worker-1`)
- `HTTP_ADDR` (default `:8080`)
- `OP_TIMEOUT` (default `5s`)
- `ARTIFACTS_ROOT` (default `./artifacts`)
- `MOCK_PORTAL_URL` (default `http://localhost:8090`)
- `PLAYWRIGHT_HEADLESS` (default `true`)

## Testing

```bash
make test
```

View test coverage:

```bash
# Show coverage percentages per package
make test-coverage

# Generate detailed HTML coverage report
make test-coverage-html
# Then open coverage.html in your browser
```

## API

- `POST /v1/tracking/jobs`
- `GET /v1/jobs/{jobId}`
- `GET /v1/tracking/results/{jobId}`

## Artifacts

Artifacts are stored under `./artifacts` using keys like:

```
provider=<provider>/yyyy=YYYY/mm=MM/dd=DD/job=<jobId>/step=<step>/file=<filename>
```
