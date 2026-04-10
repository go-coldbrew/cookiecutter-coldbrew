// Package auth provides authentication interceptors for ColdBrew gRPC services.
//
// Auth is config-controlled: set JWT_SECRET or API_KEYS environment variables to enable.
// When neither is set, auth is a no-op.
//
// The AuthConfig struct is embedded in config.Config (same pattern as cbConfig.Config)
// and Setup() is called from main() to register interceptors.
//
// References:
//   - go-grpc-middleware auth: https://github.com/grpc-ecosystem/go-grpc-middleware/tree/main/interceptors/auth
//   - grpc-go authz (policy-based authorization): https://github.com/grpc/grpc-go/tree/master/authz
//   - golang-jwt/jwt: https://github.com/golang-jwt/jwt
package auth

import (
	"context"
	"fmt"

	"github.com/go-coldbrew/interceptors"
	"github.com/golang-jwt/jwt/v5"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const apiKeyHeader = "x-api-key"

// AuthConfig holds authentication configuration loaded from environment variables.
// Embedded in config.Config (same pattern as cbConfig.Config).
type AuthConfig struct {
	JWTSecret string   `envconfig:"JWT_SECRET"`
	APIKeys   []string `envconfig:"API_KEYS"`
}

// Setup registers auth interceptors based on the loaded config.
// Called from main() after config is loaded. If neither JWTSecret nor APIKeys
// are set, this is a no-op.
func Setup(cfg AuthConfig) {
	if cfg.JWTSecret != "" {
		jwtAuth := JWTAuthFunc(cfg.JWTSecret)
		interceptors.AddUnaryServerInterceptor(context.Background(),
			grpcauth.UnaryServerInterceptor(jwtAuth))
		interceptors.AddStreamServerInterceptor(context.Background(),
			grpcauth.StreamServerInterceptor(jwtAuth))
	}
	if len(cfg.APIKeys) > 0 {
		apiKeyAuth := APIKeyAuthFunc(cfg.APIKeys)
		interceptors.AddUnaryServerInterceptor(context.Background(),
			grpcauth.UnaryServerInterceptor(apiKeyAuth))
		interceptors.AddStreamServerInterceptor(context.Background(),
			grpcauth.StreamServerInterceptor(apiKeyAuth))
	}
}

type contextKey struct{}

// Claims holds the parsed JWT claims, accessible in handlers via ClaimsFromContext.
// Subject, Issuer, ExpiresAt, etc. are available via the embedded RegisteredClaims.
type Claims struct {
	jwt.RegisteredClaims
}

// ClaimsFromContext returns the JWT claims from the context, or nil if not present.
func ClaimsFromContext(ctx context.Context) *Claims {
	c, _ := ctx.Value(contextKey{}).(*Claims)
	return c
}

// JWTAuthFunc returns an [grpcauth.AuthFunc] that validates Bearer JWT tokens
// using HMAC-SHA256. The secret is the shared signing key.
//
// To use a different signing method (RSA, ECDSA), replace jwt.SigningMethodHS256
// with the appropriate method and change the keyFunc to return your public key.
// See https://github.com/golang-jwt/jwt for details.
func JWTAuthFunc(secret string) grpcauth.AuthFunc {
	secretBytes := []byte(secret)
	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretBytes, nil
	}
	validMethods := jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})
	return func(ctx context.Context) (context.Context, error) {
		tokenStr, err := grpcauth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc, validMethods)
		if err != nil || !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		return context.WithValue(ctx, contextKey{}, claims), nil
	}
}

// APIKeyAuthFunc returns an [grpcauth.AuthFunc] that validates API keys from the
// "x-api-key" gRPC metadata header. validKeys is the set of accepted keys.
func APIKeyAuthFunc(validKeys []string) grpcauth.AuthFunc {
	keySet := make(map[string]struct{}, len(validKeys))
	for _, k := range validKeys {
		keySet[k] = struct{}{}
	}
	return func(ctx context.Context) (context.Context, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		keys := md.Get(apiKeyHeader)
		if len(keys) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing %s header", apiKeyHeader)
		}
		if _, valid := keySet[keys[0]]; !valid {
			return nil, status.Error(codes.Unauthenticated, "invalid API key")
		}
		return ctx, nil
	}
}
