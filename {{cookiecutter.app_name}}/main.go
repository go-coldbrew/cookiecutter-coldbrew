package main

import (
	"context"
	"mime"
	"net/http"

	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	{{cookiecutter.app_name|lower}} "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version"
	"github.com/go-coldbrew/core"
	"github.com/go-coldbrew/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"

	_ "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/statik"
)

type svc struct {
}

func (s *svc) InitHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return {{cookiecutter.app_name|lower}}.Register{{cookiecutter.service_name}}HandlerFromEndpoint(ctx, mux, endpoint, opts)
}

func (s *svc) InitGRPC(ctx context.Context, server *grpc.Server) error {
	{{cookiecutter.app_name|lower}}.Register{{cookiecutter.service_name}}Server(server, service.New(config.Get()))
	return nil
}

func getOpenAPIHandler() http.Handler {
	// getOpenAPIHandler serves an OpenAPI UI.
	// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		panic("creating OpenAPI filesystem: " + err.Error())
	}
	return http.FileServer(statikFS)
}

func main() {
	cfg := config.GetColdBrewConfig()
	if cfg.AppName == "" {
		cfg.AppName = "{{cookiecutter.app_name}}"
	}
	cfg.ReleaseName = version.GitCommit

	cb := core.New(cfg)
	cb.SetOpenAPIHandler(getOpenAPIHandler())
	cb.SetService(&svc{})

	log.Error(context.Background(), cb.Run())
}
