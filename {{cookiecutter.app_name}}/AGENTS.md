# AGENTS.md

## Project Overview

{{cookiecutter.project_short_description}}

This is a ColdBrew gRPC microservice built with the [ColdBrew framework](https://docs.coldbrew.cloud). It uses gRPC as the primary protocol with an HTTP/JSON gateway auto-generated from protobuf definitions.

## Build & Test Commands

```bash
make build           # Compile the project
make test            # Run tests with race detection and coverage
make lint            # Run golangci-lint and govulncheck
make vulncheck       # Run Go vulnerability check only
make bench           # Run benchmarks
make generate        # Generate gRPC/HTTP code from proto files
make mock            # Generate mocks for interfaces
make fmt             # Format Go source files
make run             # Build and run locally (Swagger UI at http://localhost:9091/swagger/)
make build-docker    # Build Docker image
make run-docker      # Run in Docker container
```

## Architecture

```
.
├── main.go              # Entry point: initializes ColdBrew, registers services
├── config/
│   └── config.go        # App configuration via environment variables (envconfig)
├── service/
│   ├── service.go       # Business logic: implements gRPC service interface
│   ├── healthcheck.go   # Kubernetes liveness/readiness probes
│   ├── service_test.go  # Unit tests and benchmarks
│   ├── healthcheck_test.go
│   └── auth/
│       ├── auth.go      # JWT + API-key auth interceptor examples (uncomment in main.go to enable)
│       └── auth_test.go
├── proto/
│   └── *.proto          # Protobuf definitions (source of truth for API)
│   └── *.pb.go          # GENERATED — do not edit
├── version/
│   └── version.go       # Build-time version info (injected via ldflags)
├── third_party/
│   └── OpenAPI/         # Swagger UI assets (embedded via go:embed)
├── .github/workflows/
│   └── go.yml           # GitHub Actions CI (build, test, bench, lint)
├── .gitlab-ci.yml       # GitLab CI (unit-test, lint, benchmark)
├── Makefile             # Build automation
├── Dockerfile           # Multi-stage Docker build
├── .golangci.yml        # Linter configuration
├── buf.yaml             # Protobuf linting rules
├── buf.gen.yaml         # Code generation plugins
└── local.env.example    # Example environment variables (copy to local.env)
```

## Key Patterns

- **gRPC-first**: All endpoints are defined in `proto/{{cookiecutter.app_name|lower}}.proto`. HTTP/JSON routes are auto-generated via grpc-gateway annotations. Never create HTTP handlers manually.
- **Context propagation**: `context.Context` is the first parameter everywhere. Interceptors propagate trace IDs, log fields, and options through it.
- **Configuration**: All config via environment variables using `envconfig`. Add fields to `config/config.go` with struct tags. See [ColdBrew config docs](https://pkg.go.dev/github.com/go-coldbrew/core/config#Config) for framework options.
- **Health checks**: Kubernetes liveness (`/healthcheck`) and readiness (`/readycheck`) are built-in. Service starts as NOT_SERVING until `SetReady()` is called.
- **Observability**: Prometheus metrics at `/metrics`, pprof at `/debug/pprof/`, OpenAPI/Swagger at `/swagger/`.
- **Graceful shutdown**: ColdBrew handles SIGINT/SIGTERM. The `Stop()` method on your service is called for cleanup.

## Development Workflows

### Adding a new endpoint

1. Define the RPC in `proto/{{cookiecutter.app_name|lower}}.proto` with HTTP annotations:
   ```protobuf
   rpc MyMethod(MyRequest) returns (MyResponse) {
     option (google.api.http) = {
       post: "/api/v1/my-endpoint"
       body: "*"
     };
   }
   ```
2. Run `make generate` to regenerate Go code
3. Implement the method in `service/service.go`
4. Add tests in `service/service_test.go`
5. Run `make test` and `make lint`

### Adding configuration

1. Add a field to the `Config` struct in `config/config.go` with an `envconfig` tag:
   ```go
   MyField string `envconfig:"MY_FIELD" default:"value"`
   ```
2. Access it via `config.Get().MyField` in your service code
3. Add the env var to `local.env.example` for documentation

### Adding tests

- Unit tests go in `service/service_test.go` alongside the code
- Use `testify/assert` for assertions
- Create the service with `New(config.Get())` to test with real config
- Benchmarks use `func BenchmarkX(b *testing.B)` with `b.ResetTimer()` before the hot loop

### Private modules

GOPRIVATE is pre-configured in Makefile, Dockerfile, and CI workflows. For private repos:
- **Local dev**: `git config --global url."git@github.com:".insteadOf "https://github.com/"` (SSH) or add a `.netrc` with a PAT
- **Docker**: uncomment the auth section in `Dockerfile`
- **CI**: uncomment the auth steps in `.github/workflows/go.yml` or `.gitlab-ci.yml`
- See [Private Modules guide](https://docs.coldbrew.cloud/howto/private-modules/) for details

## Local Development Stack

Start infrastructure with docker-compose, then run the app locally with `make run`:

```bash
make local-stack                               # start default services (selected during generation)
make local-stack PROFILES="postgres kafka obs"  # override with specific services
make run                                        # run the app (fast native build, no Docker)
make local-stack-down                           # stop infra
make local-exec SVC=postgres CMD="psql -U postgres"  # exec into a service
make local-exec SVC=redis CMD="redis-cli"            # works with any service
```

Available profiles:

| Category | Profiles |
|----------|----------|
| Databases | `postgres`, `mysql`, `cockroachdb`, `mongodb` |
| Cache | `redis`, `valkey`, `memcached` |
| Messaging | `kafka`, `nats` |
| Search | `elasticsearch` |
| AWS | `ministack`, `dynamodb` |
| GCP | `spanner`, `pubsub`, `bigtable`, `firestore`, `alloydb` |
| Tools | `adminer` |
| Observability | `obs` (Prometheus, Grafana, Jaeger) |

Service endpoints (via `make run`):
- HTTP/Swagger: http://localhost:9091/swagger/
- gRPC: localhost:9090

Obs endpoints (when running with `obs` profile):
- Grafana: http://localhost:3000 (admin/admin) — ColdBrew dashboard pre-loaded
- Jaeger: http://localhost:16686 — distributed traces
- Prometheus: http://localhost:9100

## Load Testing

Run gRPC load tests against a locally running service using [ghz](https://ghz.sh):

```bash
make run                    # start the app in one terminal
make loadtest               # run load test in another terminal
```

The default config (`misc/loadtest/echo.json`) sends requests for 10 seconds at concurrency 10 to the Echo RPC via gRPC reflection. Edit the file to adjust duration, concurrency, or target a different RPC.

With the observability stack running (`make local-stack-obs`), load test results are visible in the Grafana dashboard in real-time.

## Rules

- **Never edit generated files** — files in `proto/*.pb.go`, `proto/*_grpc.pb.go`, `proto/*.gw.go` are generated. Edit the `.proto` file and run `make generate`.
- **Always run `make generate` after proto changes** — both gRPC stubs and HTTP gateway code must be regenerated.
- **Always run `make test` with race detection** — `make test` includes `-race` by default.
- **Always run `make lint` before committing** — includes golangci-lint and govulncheck.
- **Don't add `replace` directives to go.mod** — unless doing local cross-package development, and remove them before committing.
- **Keep config in environment variables** — never hardcode secrets or environment-specific values.
- **gRPC status codes** — use `google.golang.org/grpc/codes` and `google.golang.org/grpc/status` for errors, not plain errors. Use `github.com/go-coldbrew/errors` for wrapping with stack traces.
