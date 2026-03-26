# cookiecutter-coldbrew

Powered by [Cookiecutter](https://github.com/cookiecutter/cookiecutter), Cookiecutter Coldbrew is a template for jumpstarting production-ready Go gRPC microservices with the [ColdBrew framework](https://docs.coldbrew.cloud).

## Features

- Complete gRPC service with HTTP/JSON gateway (grpc-gateway)
- Kubernetes health checks (liveness + readiness probes)
- Prometheus metrics, distributed tracing, structured logging
- Swagger UI for interactive API documentation
- Multi-stage Docker build for minimal production images
- CI/CD pipelines for GitHub Actions and GitLab CI
- golangci-lint v2 configuration with govulncheck
- Makefile with build, test, lint, benchmark, and run targets
- Build-time version injection (git commit, branch, date)

## Prerequisites

Install [Cookiecutter](https://cookiecutter.readthedocs.io/):

```shell
brew install cookiecutter
```

Or via pip:

```shell
pip install cookiecutter
```

## Usage

```shell
cookiecutter gh:go-coldbrew/cookiecutter-coldbrew
```

Answer the prompts:

```shell
source_path [github.com/ankurs]: github.com/yourname
app_name [MyApp]: EchoServer
grpc_package [com.github.ankurs]: com.github.yourname
service_name [MySvc]: EchoSvc
project_short_description [A Golang project.]: My first ColdBrew service
docker_image [alpine:latest]:
docker_build_image [golang]:
Select docker_build_image_version:
1 - 1.26
2 - 1.25
Choose from 1, 2 [1]: 1
```

Then build and run:

```shell
cd EchoServer/
make run
```

Your service starts on `:9090` (gRPC) and `:9091` (HTTP/Swagger).

For a full walkthrough, see the [Getting Started](https://docs.coldbrew.cloud/getting-started/) guide.

## CI/CD

The generated project includes ready-to-use CI pipelines for both platforms. Delete whichever you don't need.

### GitHub Actions (`.github/workflows/go.yml`)

Runs on push to `main`/`master` and on pull requests. Four parallel jobs: **build**, **test** (race detector + coverage), **benchmark**, and **lint** (govulncheck + golangci-lint). Each job has concurrency control to cancel duplicate runs.

### GitLab CI (`.gitlab-ci.yml`)

Three jobs in a single `test` stage: **unit-test** (with Cobertura coverage report), **lint** (golangci-lint + govulncheck), and **benchmark**. Go module caching is enabled.

## Docker

Uses multi-stage builds: compiles a static Go binary in the builder stage, then copies it to a minimal Alpine image. Ports 9090 (gRPC) and 9091 (HTTP) are exposed.

```shell
make build-docker
make run-docker
```

## Documentation

- [Getting Started](https://docs.coldbrew.cloud/getting-started/) — Full walkthrough
- [ColdBrew Documentation](https://docs.coldbrew.cloud) — Framework reference
- [How-To Guides](https://docs.coldbrew.cloud/howto/) — Tracing, logging, metrics, and more
