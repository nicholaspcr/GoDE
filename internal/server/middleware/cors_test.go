package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware_AllowAll(t *testing.T) {
	config := CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false, // When using *, credentials must be false
		MaxAge:           3600,
	}

	handler := CORSMiddleware(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "3600", w.Header().Get("Access-Control-Max-Age"))
}

func TestCORSMiddleware_SpecificOrigin(t *testing.T) {
	config := CORSConfig{
		AllowedOrigins: []string{"https://example.com", "https://app.example.com"},
		AllowedMethods: []string{"GET", "POST"},
	}

	handler := CORSMiddleware(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test allowed origin
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))

	// Test different allowed origin
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://app.example.com")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, "https://app.example.com", w.Header().Get("Access-Control-Allow-Origin"))

	// Test disallowed origin
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	config := CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false, // When using *, credentials must be false
	}

	handler := CORSMiddleware(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// OPTIONS request (preflight)
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, DELETE", w.Header().Get("Access-Control-Allow-Methods"))
}

func TestCORSMiddleware_WildcardSubdomain(t *testing.T) {
	config := CORSConfig{
		AllowedOrigins: []string{"*.example.com"},
		AllowedMethods: []string{"GET"},
	}

	handler := CORSMiddleware(config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test subdomain match
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://app.example.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, "https://app.example.com", w.Header().Get("Access-Control-Allow-Origin"))

	// Test another subdomain
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://api.example.com")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Equal(t, "https://api.example.com", w.Header().Get("Access-Control-Allow-Origin"))

	// Test non-matching domain
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.org")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestDefaultCORSConfig(t *testing.T) {
	config := DefaultCORSConfig()

	assert.Equal(t, []string{"*"}, config.AllowedOrigins)
	assert.Contains(t, config.AllowedMethods, "GET")
	assert.Contains(t, config.AllowedMethods, "POST")
	assert.Contains(t, config.AllowedHeaders, "Content-Type")
	assert.Contains(t, config.AllowedHeaders, "Authorization")
	assert.False(t, config.AllowCredentials) // Must be false when using "*"
	assert.Equal(t, 86400, config.MaxAge)
}

func TestProductionCORSConfig(t *testing.T) {
	origins := []string{"https://app.example.com", "https://www.example.com"}
	config := ProductionCORSConfig(origins)

	assert.Equal(t, origins, config.AllowedOrigins)
	assert.NotContains(t, config.AllowedMethods, "PATCH") // More restrictive
	assert.Contains(t, config.AllowedMethods, "GET")
	assert.Contains(t, config.AllowedMethods, "POST")
	assert.True(t, config.AllowCredentials)
	assert.Equal(t, 3600, config.MaxAge) // Shorter than default
}
