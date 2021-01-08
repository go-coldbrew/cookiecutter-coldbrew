package main

import (
	"context"
	"flag"
	"fmt"
	"mime"
	"net/http"

	"github.com/ankurs/ExampleProject/config"
	exampleproject "github.com/ankurs/ExampleProject/proto"
	"github.com/ankurs/ExampleProject/service"
	"github.com/ankurs/ExampleProject/version"
	"github.com/go-coldbrew/core"
	"github.com/go-coldbrew/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"

	_ "github.com/ankurs/ExampleProject/statik"
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

func (s *svc) GetOpenAPIHandler(ctx context.Context) http.Handler {
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

	versionFlag := flag.Bool("version", false, "Version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
		return
	}

	cb := core.New(config.GetColdBrewConfig())

	cb.SetService(&svc{})

	log.Error(context.Background(), cb.Run())
}
