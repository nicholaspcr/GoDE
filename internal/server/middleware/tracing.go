package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TracingMiddleware wraps an HTTP handler with OpenTelemetry tracing.
// It automatically creates spans for each HTTP request and propagates trace context.
func TracingMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, serviceName,
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return r.Method + " " + r.URL.Path
			}),
		)
	}
}
