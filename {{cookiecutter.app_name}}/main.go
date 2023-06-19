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
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

	_ "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/statik"
)

// cbSvc is the service implementation of ColdBrew service
type cbSvc struct {
}

// FailCheck allows graceful termination of the service
// This is called by the health check endpoint to determine if the service is ready to serve requests or not
func (s *cbSvc) FailCheck(fail bool) {
	if fail {
		service.SetNotReady()
	} else {
		service.SetReady()
	}
}

// Stop is called when the service is being stopped by the ColdBrew framework
// This is a good place to clean up resources and gracefully shutdown the service if needed before the process exits completely
func (s *cbSvc) Stop() {
	//  TODO: Add your cleanup code here
}

// InitHTTP is called by the ColdBrew framework to initialize the HTTP server and register the HTTP handlers
// This is a good place to register your HTTP handlers if you have any custom handlers that you want to register with the HTTP server
// If you are using the grpc-gateway, you can use the RegisterMySvcHandlerFromEndpoint function to register the HTTP handlers
// The endpoint is the address of the gRPC server
// The opts are the grpc.DialOptions that are used to connect to the gRPC server
func (s *cbSvc) InitHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return {{cookiecutter.app_name|lower}}.Register{{cookiecutter.service_name}}HandlerFromEndpoint(ctx, mux, endpoint, opts)
}

// InitGRPC is called by the ColdBrew framework to initialize the gRPC server and register the gRPC handlers
// This is a good place to register your gRPC handlers if you have any custom handlers that you want to register with the gRPC server
// If you are using the grpc-gateway, you can use the RegisterMySvcHandlerFromEndpoint function to register the HTTP handlers
func (s *cbSvc) InitGRPC(ctx context.Context, server *grpc.Server) error {
	// Create the service implementation
	impl, err := service.New(config.Get())
	if err != nil {
		return err
	}
	// Register the service implementation with the gRPC server
	{{cookiecutter.app_name|lower}}.Register{{cookiecutter.service_name}}Server(server, impl)

	// Register the health check service implementation with the gRPC server so that the gRPC health check endpoint is available
	healthgrpc.RegisterHealthServer(server, impl)
	return nil
}

// getOpenAPIHandler returns the OpenAPI UI handler that is used by the ColdBrew framework to serve the OpenAPI UI
func getOpenAPIHandler() http.Handler {
	// getOpenAPIHandler serves an OpenAPI UI.
	// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		log.Error(context.Background(), "msg", "error adding mime type", "err", err)
	}

	statikFS, err := fs.New()
	if err != nil {
		panic("creating OpenAPI filesystem: " + err.Error())
	}
	return http.FileServer(statikFS)
}

// main is the entry point of the service
// This is where the ColdBrew framework is initialized and the service is started
func main() {
	// Initialize the ColdBrew framework configuration from the environment variables
	cfg := config.GetColdBrewConfig()
	if cfg.AppName == "" {
		// If the app name is not set in the environment variables, use the app name from the version package
		cfg.AppName = version.AppName
	}
	// Set the release name to the git commit hash from the version package
	cfg.ReleaseName = version.GitCommit

	// Initialize the ColdBrew framework with the given configuration
	// This is a good place to customise the ColdBrew framework configuration if needed
	cb := core.New(cfg)
	// Set the OpenAPI handler that is used by the ColdBrew framework to serve the OpenAPI UI
	cb.SetOpenAPIHandler(getOpenAPIHandler())
	// Register the service implementation with the ColdBrew framework
	err := cb.SetService(&cbSvc{})
	if err != nil {
		// If there is an error registering the service implementation, panic and exit
		panic(err)
	}

	// Start the service and wait for it to exit
	// This is a blocking call and will not return until the service exits completely
	log.Error(context.Background(), cb.Run())
}
