# Go caching-rung starter

A plain `net/http` service that is **slow on purpose**. It is the starting point for the Go
briefs on the DayZer0 caching rung: measure how slow it is first, then make it fast — a cache,
a better read path, a pool, whatever your brief asks for — and prove the difference with the
same load test.

What ships here, deliberately:

- A seeded catalog of **50,000 products**, generated deterministically at startup, so every
  clone measures the same data.
- A store (`store.go`) where **every query pays a simulated ~30ms database round-trip and then
  scans the whole table**. No index, no cache. Do not delete the delay — your job is to stop
  paying it on every request, not to pretend it never existed.
- A **k6 baseline** (`loadtest/baseline.js`) whose threshold the scaffold fails. Run it before
  you change anything and keep the numbers.
- **No caching anywhere.** If your brief asks you to add a cache, you are adding the first one.

## Run it

```bash
go run .
```

The server listens on `:8080` (override with `PORT`). Confirm it is up:

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

Feel the problem: `curl http://localhost:8080/products/42` — every call pays the full round
trip. `SEED_ROWS` and `DB_DELAY_MS` tune the dataset size and the simulated round-trip.

## Test it

```bash
go test ./...
```

## Measure it

```bash
go run .                        # terminal 1
k6 run loadtest/baseline.js     # terminal 2
```

## Infrastructure

`docker-compose.yml` starts Memcached, Redis and Postgres. The scaffold uses none of them —
wiring one in is the work. `docker compose up -d memcached` (or `redis`, or `postgres`).

## Layout

- `main.go` — entry point and routes. Add yours beside `/health` and `/products`.
- `store.go` — the deliberately slow data layer. Read its top comment before touching it.
- `main_test.go` — the health and product tests, a pattern to copy for your own handlers.
- `loadtest/baseline.js` — the k6 baseline your pull requests quote.

## Working the tickets

Each ticket is a GitHub issue on this repo. For a work ticket: branch `issue-<N>-<slug>` off
`main`, commit your work, open a pull request with `Closes #<N>` in its body, and Raj reviews it.
An approved pull request merges and marks the ticket submitted.
