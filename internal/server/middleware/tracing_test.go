package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTracingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	middleware := TracingMiddleware("test-service")
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestTracingMiddleware_WithDifferentMethods(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"GET request", "GET", "/api/users"},
		{"POST request", "POST", "/api/users"},
		{"PUT request", "PUT", "/api/users/123"},
		{"DELETE request", "DELETE", "/api/users/123"},
		{"PATCH request", "PATCH", "/api/users/123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.method, r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				w.WriteHeader(http.StatusOK)
			})

			middleware := TracingMiddleware("test-service")
			wrappedHandler := middleware(handler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestTracingMiddleware_PreservesHeaders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		w.Header().Set("X-Custom-Header", "value")
		w.WriteHeader(http.StatusOK)
	})

	middleware := TracingMiddleware("test-service")
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token123")
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "value", rec.Header().Get("X-Custom-Header"))
}

func TestTracingMiddleware_HandlesErrors(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	})

	middleware := TracingMiddleware("test-service")
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/error", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Internal Server Error")
}

func TestTracingMiddleware_MultipleServices(t *testing.T) {
	services := []string{"service-a", "service-b", "service-c"}

	for _, serviceName := range services {
		t.Run(serviceName, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := TracingMiddleware(serviceName)
			require.NotNil(t, middleware)

			wrappedHandler := middleware(handler)

			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestTracingMiddleware_ChainedMiddleware(t *testing.T) {
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Chained"))
	})

	// Apply multiple middleware layers
	tracing1 := TracingMiddleware("service-1")
	tracing2 := TracingMiddleware("service-2")

	wrappedHandler := tracing1(tracing2(finalHandler))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Chained", rec.Body.String())
}

func TestTracingMiddleware_WithQueryParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value1", r.URL.Query().Get("param1"))
		assert.Equal(t, "value2", r.URL.Query().Get("param2"))
		w.WriteHeader(http.StatusOK)
	})

	middleware := TracingMiddleware("test-service")
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test?param1=value1&param2=value2", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
