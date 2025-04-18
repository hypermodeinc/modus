package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Dummy http.Handler to track invocation
func dummyNextHandler(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

func generateRSAKey() (*rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func generateSignedJWT(priv *rsa.PrivateKey, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(priv)
}

func TestHandleJWT(t *testing.T) {
	var generatedJWT string

	tests := []struct {
		name           string
		headerFunc     func() string
		setupFunc      func()
		expectStatus   int
		expectNextCall bool
	}{
		{
			name:       "Valid JWT",
			headerFunc: func() string { return "Bearer " + generatedJWT },
			setupFunc: func() {
				priv, err := generateRSAKey()
				if err != nil {
					panic(err)
				}
				claims := jwt.MapClaims{"sub": "test", "exp": time.Now().Add(time.Hour).Unix()}
				signedJWT, err := generateSignedJWT(priv, claims)
				if err != nil {
					panic(err)
				}
				globalAuthKeys = &AuthKeys{}
				globalAuthKeys.setPemPublicKeys(map[string]any{"testkey": &priv.PublicKey})
				generatedJWT = signedJWT
			},
			expectStatus:   http.StatusOK,
			expectNextCall: true,
		},
		{
			name:       "Expired JWT",
			headerFunc: func() string { return "Bearer " + generatedJWT },
			setupFunc: func() {
				priv, err := generateRSAKey()
				if err != nil {
					panic(err)
				}
				claims := jwt.MapClaims{"sub": "test", "exp": time.Now().Add(-time.Hour).Unix()}
				signedJWT, err := generateSignedJWT(priv, claims)
				if err != nil {
					panic(err)
				}
				globalAuthKeys = &AuthKeys{}
				globalAuthKeys.setPemPublicKeys(map[string]any{"testkey": &priv.PublicKey})
				generatedJWT = signedJWT
			},
			expectStatus:   http.StatusUnauthorized,
			expectNextCall: false,
		},
		{
			name:       "No Authorization header",
			headerFunc: func() string { return "" },
			setupFunc: func() {
				globalAuthKeys = &AuthKeys{}
				globalAuthKeys.setPemPublicKeys(map[string]any{"dummy": []byte("dummy")})
			},
			expectStatus:   http.StatusUnauthorized,
			expectNextCall: false,
		},
		{
			name:       "Malformed JWT",
			headerFunc: func() string { return "Bearer bad.token" },
			setupFunc: func() {
				globalAuthKeys = &AuthKeys{}
				globalAuthKeys.setPemPublicKeys(map[string]any{"dummy": []byte("dummy")})
			},
			expectStatus:   http.StatusUnauthorized,
			expectNextCall: false,
		},
		{
			name:       "Bad Bearer Prefix",
			headerFunc: func() string { return "Token sometoken" },
			setupFunc: func() {
				globalAuthKeys = &AuthKeys{}
				globalAuthKeys.setPemPublicKeys(map[string]any{"dummy": []byte("dummy")})
			},
			expectStatus:   http.StatusBadRequest,
			expectNextCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			tt.setupFunc()
			req.Header.Set("Authorization", tt.headerFunc())
			mw := HandleJWT(dummyNextHandler(&called))
			mw.ServeHTTP(rec, req)

			if rec.Code != tt.expectStatus {
				t.Errorf("expected status %d, got %d", tt.expectStatus, rec.Code)
			}
			if called != tt.expectNextCall {
				t.Errorf("expected next handler called: %v, got %v", tt.expectNextCall, called)
			}
		})
	}
}
