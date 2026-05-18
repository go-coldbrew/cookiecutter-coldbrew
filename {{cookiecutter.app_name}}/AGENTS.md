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
│       ├── auth.go      # JWT + API-key auth interceptors (enabled when JWT_SECRET/API_KEYS are set)
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
- **Context propagation**: `context.Context` is the first parameter everywhere. Interceptors propagate trace IDs, log fields, and options through it. Service code uses `slog.LogAttrs(ctx, ...)` for logging; ColdBrew's Handler automatically injects context fields. Use `github.com/go-coldbrew/log.AddAttrsToContext` (imported as `cblog` in service.go) to add typed context fields.
- **Configuration**: All config via environment variables using `envconfig`. Add fields to `config/config.go` with struct tags. See [ColdBrew config docs](https://pkg.go.dev/github.com/go-coldbrew/core/config#Config) for framework options.
- **Authentication**: JWT and API key auth are built in via `service/auth/`. Config-controlled — set `JWT_SECRET` or `API_KEYS` env vars to enable. Health/ready/reflection RPCs bypass auth automatically. See [Authentication docs](https://docs.coldbrew.cloud/howto/auth/).
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

### Adding a streaming endpoint (SSE / AI/LLM tokens)

Server-streaming RPCs are exposed by the HTTP gateway as newline-delimited JSON by default, or Server-Sent Events when the client sends `Accept: text/event-stream`. SSE support is enabled by ColdBrew out of the box (set `DISABLE_SSE_MARSHALER=true` to opt out), and the compression wrapper auto-excludes `text/event-stream` so frames are not buffered by gzip.

1. Define a server-streaming RPC:
   ```protobuf
   rpc MyStream(MyRequest) returns (stream MyEvent) {
     option (google.api.http) = {
       post: "/api/v1/my-stream"
       body: "*"
     };
   }
   ```
2. Implement using `grpc.ServerStreamingServer[MyEvent]`:
   ```go
   func (s *svc) MyStream(req *proto.MyRequest, stream grpc.ServerStreamingServer[proto.MyEvent]) error {
       ctx := stream.Context()
       for event := range produce(ctx, req) {
           if err := ctx.Err(); err != nil {
               return errors.Wrap(err, "my_stream canceled")
           }
           if err := stream.Send(event); err != nil {
               return errors.Wrap(err, "my_stream send")
           }
       }
       return nil
   }
   ```
3. **Always check `stream.Context().Err()` before each Send.** Client disconnect cancels the context — for AI/LLM workloads this is the signal to stop generating (and stop paying for) tokens. Pass the context into your LLM SDK call so cancellation propagates to the upstream provider.
4. Track time-to-first-token (TTFT) as a separate metric from total stream duration — TTFT surfaces upstream latency independently of generation throughput. See `IncStreamEchoTotal` / `ObserveStreamEchoTTFT` in `service/metrics/` for the pattern.
5. The `StreamEcho` RPC in the generated proto is a synthetic example — replace its body with your real streaming source (LLM SDK, change-data-capture feed, progress events, etc.) or delete it once you have your own streams.

> **Gateway message wrapping.** grpc-gateway wraps each streamed message in `{"result": <message>}` over HTTP — this is the gateway's documented convention for server-streaming responses. Native gRPC clients see the unwrapped message. Use `google.api.HttpBody` as the stream response type if you need full control over the wire bytes for SSE consumers.

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
