package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RateLimiter manages rate limiting for different types of requests.
type RateLimiter struct {
	// Per-IP rate limiters for login endpoints
	loginLimiters     map[string]*rate.Limiter
	loginLimiterMutex sync.RWMutex
	loginLimit        rate.Limit

	// Per-IP rate limiters for registration endpoints (stricter)
	registerLimiters     map[string]*rate.Limiter
	registerLimiterMutex sync.RWMutex
	registerLimit        rate.Limit

	// Per-user rate limiters for DE executions
	userDELimiters     map[string]*rateLimiterWithConcurrency
	userDELimiterMutex sync.RWMutex
	deLimit            rate.Limit
	maxConcurrentDE    int

	// Global rate limiter
	globalLimiter *rate.Limiter
}

// rateLimiterWithConcurrency combines rate limiting with concurrency control.
type rateLimiterWithConcurrency struct {
	limiter     *rate.Limiter
	concurrent  int
	maxConcur   int
	concurMutex sync.Mutex
}

// NewRateLimiter creates a new rate limiter with the specified limits.
func NewRateLimiter(loginRequestsPerMinute, registerRequestsPerMinute, deExecutionsPerUser, maxConcurrentDE, maxRequestsPerSecond int) *RateLimiter {
	return &RateLimiter{
		loginLimiters:    make(map[string]*rate.Limiter),
		loginLimit:       rate.Limit(float64(loginRequestsPerMinute) / 60.0), // Convert to per-second
		registerLimiters: make(map[string]*rate.Limiter),
		registerLimit:    rate.Limit(float64(registerRequestsPerMinute) / 60.0), // Convert to per-second
		userDELimiters:   make(map[string]*rateLimiterWithConcurrency),
		deLimit:          rate.Limit(float64(deExecutionsPerUser) / 60.0), // Convert to per-second
		maxConcurrentDE:  maxConcurrentDE,
		globalLimiter:    rate.NewLimiter(rate.Limit(maxRequestsPerSecond), maxRequestsPerSecond),
	}
}

// getIPFromContext extracts the client IP from the gRPC context.
func getIPFromContext(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "unknown"
	}
	return p.Addr.String()
}

// getUsernameFromContext extracts the username from the gRPC context metadata.
func getUsernameFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	usernames := md.Get("username")
	if len(usernames) == 0 {
		return ""
	}
	return usernames[0]
}

// getLoginLimiter gets or creates a rate limiter for login requests from the given IP.
func (rl *RateLimiter) getLoginLimiter(ip string) *rate.Limiter {
	rl.loginLimiterMutex.RLock()
	limiter, exists := rl.loginLimiters[ip]
	rl.loginLimiterMutex.RUnlock()

	if exists {
		return limiter
	}

	rl.loginLimiterMutex.Lock()
	defer rl.loginLimiterMutex.Unlock()

	// Double-check after acquiring write lock
	limiter, exists = rl.loginLimiters[ip]
	if exists {
		return limiter
	}

	// Create new limiter with burst of 2 to allow occasional bursts
	limiter = rate.NewLimiter(rl.loginLimit, 2)
	rl.loginLimiters[ip] = limiter

	return limiter
}

// getRegisterLimiter gets or creates a rate limiter for registration requests from the given IP.
func (rl *RateLimiter) getRegisterLimiter(ip string) *rate.Limiter {
	rl.registerLimiterMutex.RLock()
	limiter, exists := rl.registerLimiters[ip]
	rl.registerLimiterMutex.RUnlock()

	if exists {
		return limiter
	}

	rl.registerLimiterMutex.Lock()
	defer rl.registerLimiterMutex.Unlock()

	// Double-check after acquiring write lock
	limiter, exists = rl.registerLimiters[ip]
	if exists {
		return limiter
	}

	// Create new limiter with burst of 1 (stricter for registration)
	limiter = rate.NewLimiter(rl.registerLimit, 1)
	rl.registerLimiters[ip] = limiter

	return limiter
}

// getUserDELimiter gets or creates a rate limiter for the given username.
func (rl *RateLimiter) getUserDELimiter(username string) *rateLimiterWithConcurrency {
	rl.userDELimiterMutex.RLock()
	limiter, exists := rl.userDELimiters[username]
	rl.userDELimiterMutex.RUnlock()

	if exists {
		return limiter
	}

	rl.userDELimiterMutex.Lock()
	defer rl.userDELimiterMutex.Unlock()

	// Double-check after acquiring write lock
	limiter, exists = rl.userDELimiters[username]
	if exists {
		return limiter
	}

	// Create new limiter with burst of maxConcurrentDE to allow concurrent requests
	limiter = &rateLimiterWithConcurrency{
		limiter:    rate.NewLimiter(rl.deLimit, rl.maxConcurrentDE),
		concurrent: 0,
		maxConcur:  rl.maxConcurrentDE,
	}
	rl.userDELimiters[username] = limiter

	return limiter
}

// acquireConcurrent attempts to acquire a concurrency slot, returns error if max reached.
func (rlc *rateLimiterWithConcurrency) acquireConcurrent() error {
	rlc.concurMutex.Lock()
	defer rlc.concurMutex.Unlock()

	if rlc.concurrent >= rlc.maxConcur {
		return fmt.Errorf("maximum concurrent executions reached (%d)", rlc.maxConcur)
	}

	rlc.concurrent++
	return nil
}

// releaseConcurrent releases a concurrency slot.
func (rlc *rateLimiterWithConcurrency) releaseConcurrent() {
	rlc.concurMutex.Lock()
	defer rlc.concurMutex.Unlock()

	if rlc.concurrent > 0 {
		rlc.concurrent--
	}
}

// UnaryGlobalRateLimitMiddleware applies global rate limiting to all unary requests.
func (rl *RateLimiter) UnaryGlobalRateLimitMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check global rate limit
		if !rl.globalLimiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted,
				"global rate limit exceeded, please try again later")
		}

		return handler(ctx, req)
	}
}

// UnaryAuthRateLimitMiddleware applies rate limiting to authentication endpoints.
func (rl *RateLimiter) UnaryAuthRateLimitMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Only apply to auth endpoints
		if info.FullMethod != "/api.v1.AuthService/Login" &&
			info.FullMethod != "/api.v1.AuthService/Register" {
			return handler(ctx, req)
		}

		ip := getIPFromContext(ctx)
		var limiter *rate.Limiter
		var errorMsg string

		// Use different rate limiters for login vs registration
		if info.FullMethod == "/api.v1.AuthService/Register" {
			limiter = rl.getRegisterLimiter(ip)
			errorMsg = "too many registration attempts, please try again later"
		} else {
			limiter = rl.getLoginLimiter(ip)
			errorMsg = "too many login attempts, please try again later"
		}

		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "%s", errorMsg)
		}

		return handler(ctx, req)
	}
}

// UnaryDERateLimitMiddleware applies rate limiting to DE execution endpoints.
func (rl *RateLimiter) UnaryDERateLimitMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Only apply to DE execution endpoint
		if info.FullMethod != "/api.v1.DifferentialEvolutionService/Run" {
			return handler(ctx, req)
		}

		username := getUsernameFromContext(ctx)
		if username == "" {
			// This shouldn't happen if auth middleware is working correctly
			return handler(ctx, req)
		}

		limiter := rl.getUserDELimiter(username)

		// Check and acquire concurrency slot first (more important than rate limiting)
		if err := limiter.acquireConcurrent(); err != nil {
			return nil, status.Errorf(codes.ResourceExhausted,
				"maximum concurrent DE executions reached, please wait for existing executions to complete")
		}

		// Check rate limit after acquiring concurrency slot
		if !limiter.limiter.Allow() {
			// Release concurrency slot since we're not proceeding
			limiter.releaseConcurrent()
			return nil, status.Errorf(codes.ResourceExhausted,
				"too many DE execution requests, please try again later")
		}

		// Ensure we release the concurrency slot when done
		defer limiter.releaseConcurrent()

		return handler(ctx, req)
	}
}

// Cleanup removes old limiters that haven't been used recently.
// Should be called periodically (e.g., every hour) to prevent memory leaks.
func (rl *RateLimiter) Cleanup(maxAge time.Duration) {
	// For simplicity, we'll clear all limiters.
	// In a production system, you'd track last access time and only remove old ones.
	rl.loginLimiterMutex.Lock()
	rl.loginLimiters = make(map[string]*rate.Limiter)
	rl.loginLimiterMutex.Unlock()

	rl.registerLimiterMutex.Lock()
	rl.registerLimiters = make(map[string]*rate.Limiter)
	rl.registerLimiterMutex.Unlock()

	rl.userDELimiterMutex.Lock()
	rl.userDELimiters = make(map[string]*rateLimiterWithConcurrency)
	rl.userDELimiterMutex.Unlock()
}
