package middleware

import (
	"context"
	"time"

	"github.com/nicholaspcr/GoDE/internal/slo"
	"github.com/nicholaspcr/GoDE/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryMetricsMiddleware records metrics for unary RPC calls.
func UnaryMetricsMiddleware(metrics *telemetry.Metrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if metrics == nil {
			return handler(ctx, req)
		}

		start := time.Now()

		// Track in-flight requests
		metrics.APIRequestsInFlight.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method", info.FullMethod),
			),
		)
		defer metrics.APIRequestsInFlight.Add(ctx, -1,
			metric.WithAttributes(
				attribute.String("method", info.FullMethod),
			),
		)

		// Call the handler
		resp, err := handler(ctx, req)

		// Record duration
		duration := time.Since(start).Seconds()
		st, _ := status.FromError(err)
		code := st.Code()

		attrs := []attribute.KeyValue{
			attribute.String("method", info.FullMethod),
			attribute.String("status", code.String()),
		}

		metrics.APIRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		metrics.APIRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))

		// Record errors
		if err != nil && code != codes.OK {
			metrics.APIErrorsTotal.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("method", info.FullMethod),
					attribute.String("code", code.String()),
				),
			)
		}

		return resp, err
	}
}

// StreamMetricsMiddleware records metrics for stream RPC calls.
func StreamMetricsMiddleware(metrics *telemetry.Metrics) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if metrics == nil {
			return handler(srv, ss)
		}

		start := time.Now()
		ctx := ss.Context()

		// Track in-flight requests
		metrics.APIRequestsInFlight.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("method", info.FullMethod),
				attribute.String("type", "stream"),
			),
		)
		defer func() {
			metrics.APIRequestsInFlight.Add(ctx, -1,
				metric.WithAttributes(
					attribute.String("method", info.FullMethod),
					attribute.String("type", "stream"),
				),
			)
		}()

		// Call the handler
		err := handler(srv, ss)

		// Record duration
		duration := time.Since(start).Seconds()
		st, _ := status.FromError(err)
		code := st.Code()

		attrs := []attribute.KeyValue{
			attribute.String("method", info.FullMethod),
			attribute.String("type", "stream"),
			attribute.String("status", code.String()),
		}

		metrics.APIRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
		metrics.APIRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))

		// Record errors
		if err != nil && code != codes.OK {
			metrics.APIErrorsTotal.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("method", info.FullMethod),
					attribute.String("code", code.String()),
				),
			)
		}

		return err
	}
}

// RecordRateLimitExceeded records a rate limit exceeded event.
func RecordRateLimitExceeded(ctx context.Context, metrics *telemetry.Metrics, limitType string) {
	if metrics == nil {
		return
	}

	metrics.RateLimitExceeded.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("limit_type", limitType),
		),
	)
}

// RecordPanic records a panic event.
func RecordPanic(ctx context.Context, metrics *telemetry.Metrics, location string) {
	if metrics == nil {
		return
	}

	metrics.PanicsTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("location", location),
		),
	)
}

// UnaryMetricsAndSLOMiddleware records both metrics and SLO data for unary RPC calls.
func UnaryMetricsAndSLOMiddleware(metrics *telemetry.Metrics, sloTracker *slo.Tracker) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Track in-flight requests (metrics only)
		if metrics != nil {
			metrics.APIRequestsInFlight.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("method", info.FullMethod),
				),
			)
			defer metrics.APIRequestsInFlight.Add(ctx, -1,
				metric.WithAttributes(
					attribute.String("method", info.FullMethod),
				),
			)
		}

		// Call the handler
		resp, err := handler(ctx, req)

		// Record duration
		duration := time.Since(start).Seconds()
		st, _ := status.FromError(err)
		code := st.Code()
		success := (err == nil && code == codes.OK)

		// Record metrics
		if metrics != nil {
			attrs := []attribute.KeyValue{
				attribute.String("method", info.FullMethod),
				attribute.String("status", code.String()),
			}

			metrics.APIRequestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
			metrics.APIRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))

			// Record errors
			if err != nil && code != codes.OK {
				metrics.APIErrorsTotal.Add(ctx, 1,
					metric.WithAttributes(
						attribute.String("method", info.FullMethod),
						attribute.String("code", code.String()),
					),
				)
			}
		}

		// Record SLO
		if sloTracker != nil {
			sloTracker.RecordRequest(ctx, "deserver", success, duration)
		}

		return resp, err
	}
}
