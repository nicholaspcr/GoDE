package server

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/slo"
	"github.com/nicholaspcr/GoDE/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
)

// ExecutorShutdowner is an interface for shutting down the executor.
type ExecutorShutdowner interface {
	Shutdown(ctx context.Context) error
}

// lifecycle manages the complete server lifecycle: setup, runtime, and shutdown
type lifecycle struct {
	cfg    Config
	server *server

	// Setup phase resources
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider

	// Runtime phase resources
	grpcServer  *grpc.Server
	httpServer  *http.Server
	pprofServer *http.Server
	rateLimiter *middleware.RateLimiter
	cleanupDone chan struct{}
	healthSrv   *health.Server
	executor    ExecutorShutdowner
}

// newLifecycle creates a new lifecycle manager
func newLifecycle(cfg Config, srv *server, exec ExecutorShutdowner) *lifecycle {
	return &lifecycle{
		cfg:      cfg,
		server:   srv,
		executor: exec,
	}
}

// setup initializes all server components
func (l *lifecycle) setup(ctx context.Context) error {
	slog.Info("Setting up server components")

	// Setup telemetry
	if err := l.setupTelemetry(ctx); err != nil {
		return err
	}

	// Setup pprof if enabled
	if err := l.setupPprof(); err != nil {
		return err
	}

	// Setup rate limiter
	l.setupRateLimiter(ctx)

	// Setup gRPC server
	if err := l.setupGRPCServer(); err != nil {
		return err
	}

	// Setup HTTP gateway
	if err := l.setupHTTPGateway(ctx); err != nil {
		return err
	}

	return nil
}

// setupTelemetry initializes tracing and metrics
func (l *lifecycle) setupTelemetry(ctx context.Context) error {
	var err error

	// Initialize tracing if enabled
	if l.cfg.TracingEnabled {
		slog.Info("Initializing distributed tracing",
			slog.String("exporter", string(l.cfg.TracingConfig.ExporterType)),
			slog.Float64("sample_ratio", l.cfg.TracingConfig.SampleRatio),
		)

		l.tracerProvider, err = telemetry.NewTracerProvider(ctx, "deserver", l.cfg.TracingConfig)
		if err != nil {
			return err
		}

		if l.cfg.TracingConfig.ExporterType == telemetry.TracingExporterOTLP {
			slog.Info("OTLP trace exporter configured",
				slog.String("endpoint", l.cfg.TracingConfig.OTLPEndpoint),
			)
		}

		slog.Info("Distributed tracing initialized successfully")
	} else {
		slog.Info("Distributed tracing is disabled")
	}

	// Initialize metrics if enabled
	if l.cfg.MetricsEnabled {
		slog.Info("Initializing metrics", slog.String("type", string(l.cfg.MetricsType)))

		l.meterProvider, err = telemetry.NewMeterProvider("deserver", l.cfg.MetricsType)
		if err != nil {
			return err
		}

		l.server.metrics, err = telemetry.InitMetrics(ctx, "deserver")
		if err != nil {
			return err
		}
		slog.Info("Metrics initialized successfully")
	}

	// Initialize SLO tracking if enabled
	if l.cfg.SLOEnabled {
		slog.Info("Initializing SLO tracking with default objectives")

		l.server.sloTracker, err = slo.NewTracker(ctx, slo.DefaultObjectives())
		if err != nil {
			return err
		}

		slog.Info("SLO tracking initialized successfully",
			slog.Int("objectives", len(slo.DefaultObjectives())),
		)
	} else {
		slog.Info("SLO tracking is disabled")
	}

	return nil
}

// setupPprof starts the pprof server if enabled
func (l *lifecycle) setupPprof() error {
	if !l.cfg.PprofEnabled {
		return nil
	}

	slog.Info("Starting pprof server", slog.String("port", l.cfg.PprofPort))
	slog.Warn("pprof is enabled - this should only be used in development or with proper security")

	l.pprofServer = &http.Server{
		Addr:              l.cfg.PprofPort,
		Handler:           nil, // Use default http.DefaultServeMux which pprof registers to
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := l.pprofServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("pprof server error", slog.String("error", err.Error()))
		}
	}()

	return nil
}

// setupRateLimiter initializes the rate limiter and cleanup goroutine
func (l *lifecycle) setupRateLimiter(ctx context.Context) {
	l.rateLimiter = middleware.NewRateLimiter(
		l.cfg.RateLimit.LoginRequestsPerMinute,
		l.cfg.RateLimit.RegisterRequestsPerMinute,
		l.cfg.RateLimit.DEExecutionsPerUser,
		l.cfg.RateLimit.MaxConcurrentDEPerUser,
		l.cfg.RateLimit.MaxRequestsPerSecond,
	)

	// Start periodic cleanup
	l.cleanupDone = make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				l.rateLimiter.Cleanup(1 * time.Hour)
				slog.Debug("Rate limiter cleanup completed")
			case <-ctx.Done():
				close(l.cleanupDone)
				return
			}
		}
	}()
}

// setupGRPCServer creates and configures the gRPC server
func (l *lifecycle) setupGRPCServer() error {
	logger := slog.Default()

	// Prepare gRPC server options
	grpcOpts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.MaxRecvMsgSize(l.cfg.RateLimit.MaxMessageSizeBytes),
		grpc.MaxSendMsgSize(l.cfg.RateLimit.MaxMessageSizeBytes),
	}

	// Build unary interceptor chain
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		middleware.UnaryPanicRecoveryMiddleware(),
	}

	// Add metrics and SLO middleware
	if l.cfg.SLOEnabled && l.server.sloTracker != nil {
		unaryInterceptors = append(unaryInterceptors,
			middleware.UnaryMetricsAndSLOMiddleware(l.server.metrics, l.server.sloTracker),
		)
	} else {
		unaryInterceptors = append(unaryInterceptors,
			middleware.UnaryMetricsMiddleware(l.server.metrics),
		)
	}

	// Add remaining interceptors
	unaryInterceptors = append(unaryInterceptors,
		l.rateLimiter.UnaryGlobalRateLimitMiddleware(),
		l.rateLimiter.UnaryAuthRateLimitMiddleware(),        // IP-based, runs before auth
		middleware.UnaryAuthMiddleware(l.server.jwtService, l.server.revoker), // Auth BEFORE DE limiter
		l.rateLimiter.UnaryDERateLimitMiddleware(),          // User-based, needs auth
		logging.UnaryServerInterceptor(InterceptorLogger(logger)),
	)

	grpcOpts = append(grpcOpts,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(
			middleware.StreamPanicRecoveryMiddleware(),
			middleware.StreamMetricsMiddleware(l.server.metrics),
			middleware.StreamAuthMiddleware(l.server.jwtService, l.server.revoker),
			logging.StreamServerInterceptor(InterceptorLogger(logger)),
		),
	)

	// Add TLS credentials if enabled
	if l.cfg.TLS.Enabled {
		slog.Info("TLS enabled, loading certificates",
			slog.String("cert", l.cfg.TLS.CertFile),
			slog.String("key", l.cfg.TLS.KeyFile),
		)
		cert, err := tls.LoadX509KeyPair(l.cfg.TLS.CertFile, l.cfg.TLS.KeyFile)
		if err != nil {
			return err
		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		})
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	} else {
		slog.Warn("TLS is disabled - this is insecure for production use")
	}

	l.grpcServer = grpc.NewServer(grpcOpts...)

	// Register health check service
	l.healthSrv = l.server.setupHealthCheck(l.grpcServer)

	// Register handlers
	slog.Info("Registering RPC services")
	for _, handler := range l.server.handlers {
		handler.RegisterService(l.grpcServer)
	}

	return nil
}

// setupHTTPGateway creates and configures the HTTP gateway
func (l *lifecycle) setupHTTPGateway(ctx context.Context) error {
	mux := runtime.NewServeMux()

	// Prepare dial options for HTTP gateway
	var dialOpts []grpc.DialOption
	if l.cfg.TLS.Enabled {
		cert, err := tls.LoadX509KeyPair(l.cfg.TLS.CertFile, l.cfg.TLS.KeyFile)
		if err != nil {
			return err
		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		})
		dialOpts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		dialOpts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}

	// Register HTTP handlers
	slog.Info("Registering grpc-gateway handlers")
	for _, handler := range l.server.handlers {
		if err := handler.RegisterHTTPHandler(ctx, mux, l.cfg.LisAddr, dialOpts); err != nil {
			return err
		}
	}

	// Add health check HTTP endpoints
	_ = mux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"UP"}`))
	})

	_ = mux.HandlePath("GET", "/readiness", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		if !l.server.checkDatabaseHealth(ctx) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"DOWN","reason":"database unavailable"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"UP"}`))
	})

	slog.Info("Health check endpoints available at /health and /readiness")

	if l.cfg.MetricsEnabled && l.cfg.MetricsType == telemetry.MetricsExporterPrometheus {
		slog.Info("Prometheus metrics endpoint will be available at /metrics via Prometheus client")
	}

	// Wrap mux with middleware (tracing, then CORS)
	handler := http.Handler(mux)

	// Add tracing middleware if enabled
	if l.cfg.TracingEnabled {
		tracingMiddleware := middleware.TracingMiddleware("deserver-http")
		handler = tracingMiddleware(handler)
		slog.Info("HTTP tracing middleware enabled")
	}

	// Add security headers middleware
	handler = middleware.SecurityHeadersMiddleware()(handler)

	// Add CORS middleware
	corsMiddleware := middleware.CORSMiddleware(l.cfg.CORS)
	handler = corsMiddleware(handler)

	slog.Info("CORS middleware enabled",
		slog.Any("allowed_origins", l.cfg.CORS.AllowedOrigins),
		slog.Bool("allow_credentials", l.cfg.CORS.AllowCredentials),
	)

	l.httpServer = &http.Server{
		Addr:              l.cfg.HTTPPort,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	if l.cfg.TLS.Enabled {
		l.httpServer.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return nil
}

// run starts all servers and waits for shutdown signal
func (l *lifecycle) run(ctx context.Context) error {
	// Start gRPC server
	slog.Info("Creating listener")
	lis, err := net.Listen("tcp", l.cfg.LisAddr)
	if err != nil {
		return err
	}
	lisAddr := lis.Addr().String()

	slog.Info("RPC server listening on: ", slog.String("addr", lisAddr))
	go func() {
		if err := l.grpcServer.Serve(lis); err != nil {
			slog.Error("Unexpected error on the gRPC server", slog.String("error", err.Error()))
		}
	}()

	// Start HTTP server
	slog.Info("HTTP server listening on: ", slog.String("port", l.cfg.HTTPPort))
	go func() {
		var err error
		if l.cfg.TLS.Enabled {
			err = l.httpServer.ListenAndServeTLS(l.cfg.TLS.CertFile, l.cfg.TLS.KeyFile)
		} else {
			err = l.httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			slog.Error("Unexpected error on the HTTP server", slog.String("error", err.Error()))
		}
	}()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal or context cancellation
	select {
	case sig := <-sigChan:
		slog.Info("Received shutdown signal", slog.String("signal", sig.String()))
	case <-ctx.Done():
		slog.Info("Context cancelled, initiating shutdown")
	}

	return nil
}

// shutdown gracefully stops all servers and cleans up resources
func (l *lifecycle) shutdown(ctx context.Context) error {
	slog.Info("Starting graceful shutdown...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop accepting new connections
	if l.healthSrv != nil {
		l.healthSrv.Shutdown()
		slog.Info("Health check service stopped")
	}

	// Shutdown HTTP server
	if l.httpServer != nil {
		if err := l.httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error shutting down HTTP server", slog.String("error", err.Error()))
		} else {
			slog.Info("HTTP server shut down gracefully")
		}
	}

	// Gracefully stop gRPC server
	if l.grpcServer != nil {
		stopped := make(chan struct{})
		go func() {
			l.grpcServer.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			slog.Info("gRPC server shut down gracefully")
		case <-shutdownCtx.Done():
			slog.Warn("Graceful shutdown timeout, forcing stop")
			l.grpcServer.Stop()
		}
	}

	// Shutdown executor (cancels active DE executions and waits for workers)
	if l.executor != nil {
		if err := l.executor.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error shutting down executor", slog.String("error", err.Error()))
		} else {
			slog.Info("Executor shut down gracefully")
		}
	}

	// Wait for cleanup goroutine with timeout
	if l.cleanupDone != nil {
		select {
		case <-l.cleanupDone:
			slog.Info("Rate limiter cleanup goroutine stopped")
		case <-shutdownCtx.Done():
			slog.Warn("Timeout waiting for cleanup goroutine, continuing shutdown")
		}
	}

	// Shutdown pprof server
	if l.pprofServer != nil {
		pprofCtx, pprofCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer pprofCancel()
		if err := l.pprofServer.Shutdown(pprofCtx); err != nil {
			slog.Error("failed to shutdown pprof server", slog.String("error", err.Error()))
		}
	}

	// Shutdown telemetry
	if l.meterProvider != nil {
		if err := l.meterProvider.Shutdown(shutdownCtx); err != nil {
			slog.Error("failed to shutdown meter provider", slog.String("error", err.Error()))
		}
	}

	if l.tracerProvider != nil {
		if err := l.tracerProvider.Shutdown(shutdownCtx); err != nil {
			slog.Error("failed to shutdown tracer provider", slog.String("error", err.Error()))
		}
	}

	slog.Info("Server shutdown complete")
	return nil
}
