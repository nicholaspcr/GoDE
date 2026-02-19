package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	handler := SecurityHeadersMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "0", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
	assert.Contains(t, w.Header().Get("Content-Security-Policy"), "default-src 'none'")
	assert.Contains(t, w.Header().Get("Permissions-Policy"), "camera=()")
}

func TestSecurityHeadersMiddleware_PassesThrough(t *testing.T) {
	called := false
	handler := SecurityHeadersMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Header().Set("X-Custom", "test")
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "test", w.Header().Get("X-Custom"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
}
