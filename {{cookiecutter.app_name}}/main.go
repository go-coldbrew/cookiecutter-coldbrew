package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	{{cookiecutter.app_name|lower}} "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service/auth"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version"
	"github.com/go-coldbrew/core"
	"github.com/go-coldbrew/errors"
	"github.com/go-coldbrew/workers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5emb"
	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

	openapi "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/third_party/OpenAPI"
)

// Compile-time interface assertions — remove or adjust as you customize your service.
var (
	_ core.CBService         = (*cbSvc)(nil)
	_ core.CBStopper         = (*cbSvc)(nil)
	_ core.CBGracefulStopper = (*cbSvc)(nil)
	_ core.CBPreStarter      = (*cbSvc)(nil)
	_ core.CBWorkerProvider  = (*cbSvc)(nil)
)

// cbSvc is the ColdBrew service adapter. It delegates to the service
// implementation in service/service.go. Optional interfaces (CBPreStarter,
// CBWorkerProvider, etc.) are discovered automatically by ColdBrew's Run().
// impl stores the concrete *service.svc; typed as any so it can be passed
// to Register*Server and health server registration without exporting the type.
// Feel free to replace with a concrete or interface type if you prefer stronger typing.
type cbSvc struct {
	impl any
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
func (s *cbSvc) Stop() {
	if impl, ok := s.impl.(interface{ Stop() }); ok {
		impl.Stop()
	}
}

// PreStart is called before gRPC/HTTP servers start. Use this for initialization
// that must complete before accepting traffic: creating the service impl,
// auth interceptors, database connections, etc. Returning an error aborts startup.
func (s *cbSvc) PreStart(ctx context.Context) error {
	cfg := config.Get()

	impl, err := service.New(cfg)
	if err != nil {
		return err
	}
	s.impl = impl

	// Register auth interceptors (JWT_SECRET or API_KEYS env vars to enable).
	// See service/auth/auth.go and https://docs.coldbrew.cloud/howto/auth/
	auth.Setup(ctx, cfg.AuthConfig)
	return nil
}

// Workers delegates to the service implementation which owns its background workers.
func (s *cbSvc) Workers() []*workers.Worker {
	if impl, ok := s.impl.(interface{ Workers() []*workers.Worker }); ok {
		return impl.Workers()
	}
	return nil
}

// InitHTTP is called by the ColdBrew framework to initialize the HTTP server and register the HTTP handlers
// This is a good place to register your HTTP handlers if you have any custom handlers that you want to register with the HTTP server
// If you are using the grpc-gateway, you can use the RegisterMySvcHandlerFromEndpoint function to register the HTTP handlers
// The endpoint is the address of the gRPC server
// The opts are the grpc.DialOptions that are used to connect to the gRPC server
func (s *cbSvc) InitHTTP(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return {{cookiecutter.app_name|lower}}.Register{{cookiecutter.service_name}}HandlerFromEndpoint(ctx, mux, endpoint, opts)
}

// InitGRPC registers the service with the gRPC server.
// The service impl is created in PreStart — InitGRPC just registers it.
func (s *cbSvc) InitGRPC(ctx context.Context, server *grpc.Server) error {
	if s.impl == nil {
		return errors.New("nil service implementation; PreStart not run")
	}

	svcServer, ok := s.impl.({{cookiecutter.app_name|lower}}.{{cookiecutter.service_name}}Server)
	if !ok {
		return errors.Wrap(fmt.Errorf("expected {{cookiecutter.service_name}}Server, got %T", s.impl), "InitGRPC")
	}
	{{cookiecutter.app_name|lower}}.Register{{cookiecutter.service_name}}Server(server, svcServer)

	healthServer, ok := s.impl.(healthgrpc.HealthServer)
	if !ok {
		return errors.Wrap(fmt.Errorf("expected HealthServer, got %T", s.impl), "InitGRPC")
	}
	healthgrpc.RegisterHealthServer(server, healthServer)
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
		r.URL.Path = prefix + strings.TrimPrefix(r.URL.Path, "/")
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
	if err := cb.Run(); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "service exited", slog.Any("err", err))
	}
}
