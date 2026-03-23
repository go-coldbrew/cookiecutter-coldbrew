# CLAUDE.md

## Project Overview

{{cookiecutter.project_short_description}}

This is a ColdBrew gRPC microservice built with the [ColdBrew framework](https://docs.coldbrew.cloud).

## Build & Test Commands

```bash
make build           # Compile the project
make test            # Run tests with race detection and coverage
make lint            # Run golangci-lint
make bench           # Run benchmarks
make generate        # Generate code from proto files
make mock            # Generate mocks for interfaces
make fmt             # Format Go source files
make run             # Build and run locally (Swagger UI at http://localhost:9091/swagger/)
make build-docker    # Build Docker image
make run-docker      # Run in Docker container
```

## Architecture

- `main.go` — Entry point, initializes ColdBrew framework and registers services
- `service/` — Service implementation (business logic, health checks)
- `config/` — Configuration via environment variables (envconfig)
- `proto/` — Protocol buffer definitions and generated code
- `version/` — Build-time version info injected via ldflags

## Key Patterns

- **gRPC-first**: Define endpoints in `proto/{{cookiecutter.app_name|lower}}.proto`, run `make generate` to create gRPC + HTTP handlers
- **Health checks**: Kubernetes liveness (`/healthcheck`) and readiness (`/readycheck`) are built-in
- **Configuration**: All config via environment variables; see `config/config.go` and [ColdBrew config docs](https://pkg.go.dev/github.com/go-coldbrew/core/config#Config)
- **Observability**: Prometheus metrics at `/metrics`, pprof at `/debug/pprof/`, OpenAPI at `/swagger/`
