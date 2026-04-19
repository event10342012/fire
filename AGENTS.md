# Repository Guidelines

## Project Structure & Module Organization
- `main.go` is the application entrypoint; `wire.go`/`wire_gen.go` handle DI wiring.
- `internal/` contains core application code (domain, service, repository, web), plus `internal/integration/` for integration tests.
- `pkg/` hosts reusable packages (e.g., logger, limiter, ginx middleware); `model/` holds data models.
- `config/` and `config.yml` store configuration; `ioc/` provides infra setup (DB, Redis, logger, web).
- `script/` and `scripts/` hold DB init and load-testing assets; `k8s/` holds deployment manifests.

## Build, Test, and Development Commands
- `go run ./main.go -c config/dev.yml` starts the server with the dev config (override with `-c config.yml`).
- `go test ./...` runs all unit and integration tests.
- `go build ./...` builds the project binaries locally.
- `make mock` regenerates mocks via `mockgen` for service/repository interfaces.
- `make docker` builds a Linux/ARM binary and Docker image `event10342012/fire:v0.0.1`.
- `docker-compose up` uses `docker-compose.yaml` to start MySQL, Redis, and etcd.

## Coding Style & Naming Conventions
- Use `gofmt` (tabs, standard Go formatting).
- Package names are lowercase; exported identifiers use `PascalCase`, unexported use `camelCase`.
- Test files follow `*_test.go`; generated mocks use `*.mock.go` and should not be hand-edited.

## Testing Guidelines
- Tests use Go’s `testing` package and `testify/assert`.
- Unit tests live alongside code under `internal/`; integration tests are under `internal/integration/`.
- Run integration tests after bringing up dependencies (`docker-compose up`), then `go test ./internal/integration/...`.

## Commit & Pull Request Guidelines
- Commit messages are short, lowercase, and imperative (examples: `add viper config`, `update`).
- PRs should include: concise summary, testing notes (commands + results), and links to relevant issues.
- Include screenshots or request/response samples for web/API changes when behavior is user-facing.

## Configuration & Local Dependencies
- Config is loaded by Viper; default is `config/config.yml`, override with `-c`.
- MySQL init SQL lives in `script/mysql/`; Redis and MySQL are defined in `docker-compose.yaml`.

## Architecture Overview
- Request flow is `internal/web` handlers → `internal/service` business logic → `internal/repository` adapters for DB/cache.
- Persistence is split between `internal/repository/dao` (GORM + SQL) and `internal/repository/cache` (Redis).
- Infra bootstrapping lives under `ioc/`, with Wire (`wire.go`, `wire_gen.go`) assembling dependencies.
