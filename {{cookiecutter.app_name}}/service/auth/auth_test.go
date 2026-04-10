package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const testSecret = "test-secret-key"

func signToken(t *testing.T, claims jwt.Claims, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return s
}

func ctxWithBearer(token string) context.Context {
	md := metadata.Pairs("authorization", "bearer "+token)
	return metadata.NewIncomingContext(context.Background(), md)
}

func ctxWithAPIKey(key string) context.Context {
	md := metadata.Pairs(apiKeyHeader, key)
	return metadata.NewIncomingContext(context.Background(), md)
}

func TestJWTAuthFunc_ValidToken(t *testing.T) {
	authFunc := JWTAuthFunc(testSecret)

	tokenStr := signToken(t, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, testSecret)

	ctx, err := authFunc(ctxWithBearer(tokenStr))
	require.NoError(t, err)

	claims := ClaimsFromContext(ctx)
	require.NotNil(t, claims)
	assert.Equal(t, "user-123", claims.Subject)
}

func TestJWTAuthFunc_ExpiredToken(t *testing.T) {
	authFunc := JWTAuthFunc(testSecret)

	tokenStr := signToken(t, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}, testSecret)

	_, err := authFunc(ctxWithBearer(tokenStr))
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestJWTAuthFunc_WrongSecret(t *testing.T) {
	authFunc := JWTAuthFunc(testSecret)

	tokenStr := signToken(t, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, "wrong-secret")

	_, err := authFunc(ctxWithBearer(tokenStr))
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestJWTAuthFunc_MissingToken(t *testing.T) {
	authFunc := JWTAuthFunc(testSecret)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	_, err := authFunc(ctx)
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestJWTAuthFunc_InvalidSigningMethod(t *testing.T) {
	authFunc := JWTAuthFunc(testSecret)

	// Sign with HS384 but the auth func only accepts HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	tokenStr, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)

	_, err = authFunc(ctxWithBearer(tokenStr))
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestClaimsFromContext_NoClaims(t *testing.T) {
	claims := ClaimsFromContext(context.Background())
	assert.Nil(t, claims)
}

func TestAPIKeyAuthFunc_ValidKey(t *testing.T) {
	authFunc := APIKeyAuthFunc([]string{"key-1", "key-2"})

	ctx, err := authFunc(ctxWithAPIKey("key-1"))
	require.NoError(t, err)
	assert.NotNil(t, ctx)
}

func TestAPIKeyAuthFunc_InvalidKey(t *testing.T) {
	authFunc := APIKeyAuthFunc([]string{"key-1"})

	_, err := authFunc(ctxWithAPIKey("wrong-key"))
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAPIKeyAuthFunc_MissingHeader(t *testing.T) {
	authFunc := APIKeyAuthFunc([]string{"key-1"})

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	_, err := authFunc(ctx)
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAPIKeyAuthFunc_MissingMetadata(t *testing.T) {
	authFunc := APIKeyAuthFunc([]string{"key-1"})

	_, err := authFunc(context.Background())
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestSkipMethodsAuthFunc(t *testing.T) {
	// Auth func that always fails — should be skipped for health methods
	authFunc := func(ctx context.Context) (context.Context, error) {
		return nil, status.Error(codes.Unauthenticated, "should not be called")
	}
	wrapped := skipMethodsAuthFunc(authFunc, defaultSkipMethods)

	tests := []struct {
		method string
		skip   bool
	}{
		// Exact matches from defaultSkipMethods
		{"/{{cookiecutter.grpc_package}}.{{cookiecutter.service_name}}/HealthCheck", true},
		{"/{{cookiecutter.grpc_package}}.{{cookiecutter.service_name}}/ReadyCheck", true},
		{"/grpc.health.v1.Health/Check", true},
		{"/grpc.health.v1.Health/Watch", true},
		{"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo", true},
		{"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo", true},
		// Should NOT skip — exact match only, no substring
		{"/myservice.v1.Svc/GetHealthCheckStatus", false},
		{"/myservice.v1.Svc/Echo", false},
	}
	for _, tt := range tests {
		ctx := grpc.NewContextWithServerTransportStream(context.Background(), &fakeStream{method: tt.method})
		_, err := wrapped(ctx)
		if tt.skip {
			assert.NoError(t, err, "expected skip for %s", tt.method)
		} else {
			assert.Error(t, err, "expected auth for %s", tt.method)
		}
	}
}

// fakeStream implements grpc.ServerTransportStream for testing grpc.Method().
type fakeStream struct {
	method string
}

func (s *fakeStream) Method() string                  { return s.method }
func (s *fakeStream) SetHeader(metadata.MD) error     { return nil }
func (s *fakeStream) SendHeader(metadata.MD) error    { return nil }
func (s *fakeStream) SetTrailer(metadata.MD) error    { return nil }
