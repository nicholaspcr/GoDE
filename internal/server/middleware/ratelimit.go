package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// maxLimiterEntries is the maximum number of per-IP/user rate limiter entries
// to prevent memory exhaustion from a large number of unique IPs/users.
const maxLimiterEntries = 10000

// RateLimiter manages rate limiting for different types of requests.
type RateLimiter struct {
	// Per-IP rate limiters for login endpoints
	loginLimiters     map[string]*limiterWithTimestamp
	loginLimiterMutex sync.RWMutex
	loginLimit        rate.Limit

	// Per-IP rate limiters for registration endpoints (stricter)
	registerLimiters     map[string]*limiterWithTimestamp
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

// limiterWithTimestamp wraps a rate.Limiter with last access tracking for LRU eviction.
type limiterWithTimestamp struct {
	limiter    *rate.Limiter
	lastAccess time.Time
	mu         sync.Mutex
}

// rateLimiterWithConcurrency combines rate limiting with concurrency control.
type rateLimiterWithConcurrency struct {
	limiter     *rate.Limiter
	concurrent  int
	maxConcur   int
	concurMutex sync.Mutex
	lastAccess  time.Time
}

// NewRateLimiter creates a new rate limiter with the specified limits.
func NewRateLimiter(loginRequestsPerMinute, registerRequestsPerMinute, deExecutionsPerUser, maxConcurrentDE, maxRequestsPerSecond int) *RateLimiter {
	return &RateLimiter{
		loginLimiters:    make(map[string]*limiterWithTimestamp),
		loginLimit:       rate.Limit(float64(loginRequestsPerMinute) / 60.0), // Convert to per-second
		registerLimiters: make(map[string]*limiterWithTimestamp),
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

// getUsernameFromContext extracts the username from the authenticated context.
// This reads from the context values set by the auth middleware (JWT claims),
// NOT from client-supplied metadata which could be spoofed.
func getUsernameFromContext(ctx context.Context) string {
	return UsernameFromContext(ctx)
}

// getLoginLimiter gets or creates a rate limiter for login requests from the given IP.
func (rl *RateLimiter) getLoginLimiter(ip string) *rate.Limiter {
	rl.loginLimiterMutex.RLock()
	lwt, exists := rl.loginLimiters[ip]
	rl.loginLimiterMutex.RUnlock()

	if exists {
		lwt.mu.Lock()
		lwt.lastAccess = time.Now()
		lwt.mu.Unlock()
		return lwt.limiter
	}

	rl.loginLimiterMutex.Lock()
	defer rl.loginLimiterMutex.Unlock()

	// Double-check after acquiring write lock
	lwt, exists = rl.loginLimiters[ip]
	if exists {
		lwt.mu.Lock()
		lwt.lastAccess = time.Now()
		lwt.mu.Unlock()
		return lwt.limiter
	}

	// Prevent unbounded map growth
	if len(rl.loginLimiters) >= maxLimiterEntries {
		return rate.NewLimiter(rl.loginLimit, 1) // Restrictive fallback
	}

	// Create new limiter with burst of 2 to allow occasional bursts
	lwt = &limiterWithTimestamp{
		limiter:    rate.NewLimiter(rl.loginLimit, 2),
		lastAccess: time.Now(),
	}
	rl.loginLimiters[ip] = lwt

	return lwt.limiter
}

// getRegisterLimiter gets or creates a rate limiter for registration requests from the given IP.
func (rl *RateLimiter) getRegisterLimiter(ip string) *rate.Limiter {
	rl.registerLimiterMutex.RLock()
	lwt, exists := rl.registerLimiters[ip]
	rl.registerLimiterMutex.RUnlock()

	if exists {
		lwt.mu.Lock()
		lwt.lastAccess = time.Now()
		lwt.mu.Unlock()
		return lwt.limiter
	}

	rl.registerLimiterMutex.Lock()
	defer rl.registerLimiterMutex.Unlock()

	// Double-check after acquiring write lock
	lwt, exists = rl.registerLimiters[ip]
	if exists {
		lwt.mu.Lock()
		lwt.lastAccess = time.Now()
		lwt.mu.Unlock()
		return lwt.limiter
	}

	// Prevent unbounded map growth
	if len(rl.registerLimiters) >= maxLimiterEntries {
		return rate.NewLimiter(rl.registerLimit, 1) // Restrictive fallback
	}

	// Create new limiter with burst of 1 (stricter for registration)
	lwt = &limiterWithTimestamp{
		limiter:    rate.NewLimiter(rl.registerLimit, 1),
		lastAccess: time.Now(),
	}
	rl.registerLimiters[ip] = lwt

	return lwt.limiter
}

// getUserDELimiter gets or creates a rate limiter for the given username.
func (rl *RateLimiter) getUserDELimiter(username string) *rateLimiterWithConcurrency {
	rl.userDELimiterMutex.RLock()
	limiter, exists := rl.userDELimiters[username]
	rl.userDELimiterMutex.RUnlock()

	if exists {
		limiter.concurMutex.Lock()
		limiter.lastAccess = time.Now()
		limiter.concurMutex.Unlock()
		return limiter
	}

	rl.userDELimiterMutex.Lock()
	defer rl.userDELimiterMutex.Unlock()

	// Double-check after acquiring write lock
	limiter, exists = rl.userDELimiters[username]
	if exists {
		limiter.concurMutex.Lock()
		limiter.lastAccess = time.Now()
		limiter.concurMutex.Unlock()
		return limiter
	}

	// Prevent unbounded map growth
	if len(rl.userDELimiters) >= maxLimiterEntries {
		return &rateLimiterWithConcurrency{
			limiter:    rate.NewLimiter(rl.deLimit, 1), // Restrictive fallback
			concurrent: 0,
			maxConcur:  1,
			lastAccess: time.Now(),
		}
	}

	// Create new limiter with burst of maxConcurrentDE to allow concurrent requests
	limiter = &rateLimiterWithConcurrency{
		limiter:    rate.NewLimiter(rl.deLimit, rl.maxConcurrentDE),
		concurrent: 0,
		maxConcur:  rl.maxConcurrentDE,
		lastAccess: time.Now(),
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
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
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
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
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
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Only apply to DE execution endpoint
		if info.FullMethod != "/api.v1.DifferentialEvolutionService/RunAsync" {
			return handler(ctx, req)
		}

		username := getUsernameFromContext(ctx)
		if username == "" {
			// Return auth error instead of bypassing rate limiting
			return nil, status.Error(codes.Unauthenticated, "username not found in context for DE rate limiting")
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
// Only entries older than maxAge will be removed.
func (rl *RateLimiter) Cleanup(maxAge time.Duration) {
	now := time.Now()
	cutoff := now.Add(-maxAge)

	rl.cleanupLoginLimiters(cutoff)
	rl.cleanupRegisterLimiters(cutoff)
	rl.cleanupDELimiters(cutoff)
}

// cleanupLoginLimiters removes stale login limiters.
// Uses a write lock throughout to avoid TOCTOU race conditions.
// This is acceptable since cleanup runs infrequently.
func (rl *RateLimiter) cleanupLoginLimiters(cutoff time.Time) {
	rl.loginLimiterMutex.Lock()
	defer rl.loginLimiterMutex.Unlock()

	for ip, lwt := range rl.loginLimiters {
		lwt.mu.Lock()
		isStale := lwt.lastAccess.Before(cutoff)
		lwt.mu.Unlock()

		if isStale {
			delete(rl.loginLimiters, ip)
		}
	}
}

// cleanupRegisterLimiters removes stale register limiters.
// Uses a write lock throughout to avoid TOCTOU race conditions.
func (rl *RateLimiter) cleanupRegisterLimiters(cutoff time.Time) {
	rl.registerLimiterMutex.Lock()
	defer rl.registerLimiterMutex.Unlock()

	for ip, lwt := range rl.registerLimiters {
		lwt.mu.Lock()
		isStale := lwt.lastAccess.Before(cutoff)
		lwt.mu.Unlock()

		if isStale {
			delete(rl.registerLimiters, ip)
		}
	}
}

// cleanupDELimiters removes stale DE limiters.
// Uses a write lock throughout to avoid TOCTOU race conditions.
func (rl *RateLimiter) cleanupDELimiters(cutoff time.Time) {
	rl.userDELimiterMutex.Lock()
	defer rl.userDELimiterMutex.Unlock()

	for username, limiter := range rl.userDELimiters {
		limiter.concurMutex.Lock()
		// Only delete if no active executions and entry is old enough
		isStale := limiter.concurrent == 0 && limiter.lastAccess.Before(cutoff)
		limiter.concurMutex.Unlock()

		if isStale {
			delete(rl.userDELimiters, username)
		}
	}
}

// StartCleanupRoutine starts a background goroutine that periodically cleans up stale limiters.
// Returns a cancel function to stop the cleanup routine.
func (rl *RateLimiter) StartCleanupRoutine(ctx context.Context, interval, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				rl.Cleanup(maxAge)
			}
		}
	}()
}
