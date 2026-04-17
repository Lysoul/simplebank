# Simple Bank

A small banking service written in Go with Gin and PostgreSQL.

## Overview

Simple Bank provides basic account management and transfer functionality. The service supports:

- Create account
- Get account by ID
- List accounts
- Delete account
- Create money transfers between accounts

The application exposes a REST API under the path prefix `/simplebank` by default.

## Architecture

- `cmd/service/main.go` — application entrypoint
- `app` — startup and configuration logic
- `accounts` — account and transfer domain logic
- `adapters/ginserver` — Gin HTTP server setup
- `adapters/postgres` — PostgreSQL connection and configuration
- `docs/er-diagram.md` — entity relationship diagram for the database schema

## ER Diagram

The ER diagram is available in `./docs/er-diagram.md` and rendered as an image below.

![Simple Bank ER Diagram](docs/er-diagram.png)

It includes the following tables:

- `ACCOUNTS`
- `ENTRIES`
- `TRANSFERS`

And relationships:

- `ACCOUNTS` has many `ENTRIES`
- `ACCOUNTS` sends and receives many `TRANSFERS`

## OpenAPI Specification

The OpenAPI spec is available at `./docs/openapi.yaml`.

This spec defines the REST operations, request/response schemas, and standard error responses used by the service.

## Testing with rest.http

A ready-to-use request collection is available at `./docs/rest.http`.

To use it:

1. Start the service with Docker Compose or locally.
2. Open `docs/rest.http` in VS Code.
3. Make sure the `@base_url` variable at the top matches the running app URL, for example:

```http
@base_url=http://localhost:3002/simplebank
```

4. Use the VS Code REST Client extension or any HTTP client that supports `.http` files to send requests.

The file includes requests for:

- `POST /accounts` to create an account
- `GET /accounts/{id}` to fetch a specific account
- `GET /accounts` to list accounts
- `DELETE /accounts/{id}` to delete an account
- `POST /accounts/transfers` to create a transfer

## Prerequisites

- Docker
- Docker Compose
- Go 1.25 (for local build without Docker)

## Run with Docker Compose

This is the recommended setup for development and testing.

```bash
docker-compose up --build
```

Open the service at:

```text
http://localhost:3002/simplebank
```


## Build and run from the Dockerfile

If you want to build the app image directly, use:

```bash
docker build -t simplebank .
```

Then run it with:

```bash
docker run --rm -p 3002:3000 \
  -e POSTGRES_URL="postgres://simplebank_user:simplebank_password@host.docker.internal:5432/simplebank?sslmode=disable" \
  -e HTTP_PORT=3000 \
  -e HTTP_PATH_PREFIX=/simplebank \
  -e POSTGRES_MIGRATE=true \
  -e HTTP_MODE=release \
  simplebank
```

> Note: `host.docker.internal` is used for Docker Desktop on Windows/macOS. If you run PostgreSQL from another container, use the service hostname instead.

## Run locally without Docker

If you prefer to run the service directly from Go:

```bash
go run cmd/service/main.go app start
```

The application reads configuration from environment variables. Example values:

```bash
POSTGRES_URL="postgres://simplebank_user:simplebank_password@localhost:5432/simplebank?sslmode=disable"
HTTP_PORT=3000
HTTP_PATH_PREFIX=/simplebank
POSTGRES_MIGRATE=true
HTTP_MODE=release
```

## Notes

- The service uses environment variables via `envconfig`.
- Database migrations are handled at startup when `POSTGRES_MIGRATE=true`.
- The service binds internally to port `3000` and is exposed as `3002` in Docker Compose.
