package main

import (
	"context"
	"net/http"

	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	{{cookiecutter.app_name|lower}} "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service/auth"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version"
	"github.com/go-coldbrew/core"
	"github.com/go-coldbrew/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"
	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

	openapi "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/third_party/OpenAPI"
)

// Compile-time interface assertions.
var (
	_ core.CBService          = (*cbSvc)(nil)
	_ core.CBStopper          = (*cbSvc)(nil)
	_ core.CBGracefulStopper  = (*cbSvc)(nil)
)

// cbSvc is the service implementation of ColdBrew service
type cbSvc struct {
	stopper core.CBStopper
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
	s.stopper.Stop()

	// Add your additional cleanup code here if needed
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

	// register stopper
	s.stopper = impl
	return nil
}

// getOpenAPIHandler returns the OpenAPI UI handler powered by swgui (Swagger UI v5).
// ColdBrew mounts this at /swagger/ with StripPrefix — we prepend the prefix back
// so swgui can route its assets correctly.
func getOpenAPIHandler() http.Handler {
	const prefix = "/swagger/"
	specFile := "{{cookiecutter.app_name|lower}}.swagger.json"
	specHandler := http.FileServerFS(openapi.SpecFS)
	uiHandler := v5emb.NewWithConfig(swgui.Config{
		SettingsUI: map[string]string{
			// Auto-prepend "Bearer " to the BearerJWT auth value in Swagger UI
			"requestInterceptor": `function(req) {
				var auth = req.headers["Authorization"];
				if (auth && !auth.startsWith("Bearer ")) {
					req.headers["Authorization"] = "Bearer " + auth;
				}
				return req;
			}`,
		},
	})("{{cookiecutter.service_name}}", prefix+specFile, prefix)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the generated spec JSON
		if r.URL.Path == "/"+specFile || r.URL.Path == specFile {
			r.URL.Path = "/" + specFile
			specHandler.ServeHTTP(w, r)
			return
		}
		// Restore the prefix that ColdBrew stripped so swgui can route assets
		r.URL.Path = prefix + r.URL.Path
		r.RequestURI = prefix + r.RequestURI
		uiHandler.ServeHTTP(w, r)
	})
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

	// Register auth interceptors if JWT_SECRET or API_KEYS env vars are set.
	// See service/auth/auth.go and https://docs.coldbrew.cloud/howto/auth/
	auth.Setup(context.Background(), config.Get().AuthConfig)

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
