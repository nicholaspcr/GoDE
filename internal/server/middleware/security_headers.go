package middleware

import "net/http"

// SecurityHeadersMiddleware adds standard HTTP security headers to all responses.
func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "0") // Disabled per modern best practice
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
			w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

			next.ServeHTTP(w, r)
		})
	}
}
