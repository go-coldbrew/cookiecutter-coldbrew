# third_party

The `OpenAPI/` directory contains the generated OpenAPI specification JSON file
produced by `make generate` (via `buf generate` with the `openapiv2` plugin).

Swagger UI is served via the [swaggest/swgui](https://github.com/swaggest/swgui)
Go package — no static UI files are vendored here. Update the UI version with
`go get -u github.com/swaggest/swgui`.
