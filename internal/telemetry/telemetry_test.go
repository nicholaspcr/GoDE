package telemetry

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTracerProvider_Stdout(t *testing.T) {
	cfg := TracingConfig{
		ExporterType: TracingExporterStdout,
		SampleRatio:  1.0,
	}

	tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
	require.NoError(t, err)
	require.NotNil(t, tp)

	err = tp.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestNewTracerProvider_OTLP(t *testing.T) {
	t.Run("insecure", func(t *testing.T) {
		cfg := TracingConfig{
			ExporterType: TracingExporterOTLP,
			OTLPEndpoint: "localhost:4317",
			SampleRatio:  1.0,
		}

		tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
		require.NoError(t, err)
		require.NotNil(t, tp)

		err = tp.Shutdown(context.Background())
		assert.NoError(t, err)
	})

	t.Run("with TLS", func(t *testing.T) {
		tmpDir := t.TempDir()
		certPath, keyPath := generateTestCert(t, tmpDir)

		cfg := TracingConfig{
			ExporterType: TracingExporterOTLP,
			OTLPEndpoint: "localhost:4317",
			SampleRatio:  1.0,
			TLS: TLSConfig{
				Enabled:    true,
				CertFile:   certPath,
				KeyFile:    keyPath,
				SkipVerify: true,
			},
		}

		tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
		require.NoError(t, err)
		require.NotNil(t, tp)

		err = tp.Shutdown(context.Background())
		assert.NoError(t, err)
	})

	t.Run("with invalid TLS cert", func(t *testing.T) {
		cfg := TracingConfig{
			ExporterType: TracingExporterOTLP,
			OTLPEndpoint: "localhost:4317",
			SampleRatio:  1.0,
			TLS: TLSConfig{
				Enabled:  true,
				CertFile: "/nonexistent/cert.pem",
				KeyFile:  "/nonexistent/key.pem",
			},
		}

		tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
		assert.Error(t, err)
		assert.Nil(t, tp)
		assert.Contains(t, err.Error(), "failed to build TLS credentials")
	})

	t.Run("without endpoint", func(t *testing.T) {
		cfg := TracingConfig{
			ExporterType: TracingExporterOTLP,
			SampleRatio:  0.5,
		}

		tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
		require.NoError(t, err)
		require.NotNil(t, tp)

		err = tp.Shutdown(context.Background())
		assert.NoError(t, err)
	})

	t.Run("with never sample", func(t *testing.T) {
		cfg := TracingConfig{
			ExporterType: TracingExporterOTLP,
			OTLPEndpoint: "localhost:4317",
			SampleRatio:  0.0,
		}

		tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
		require.NoError(t, err)
		require.NotNil(t, tp)

		err = tp.Shutdown(context.Background())
		assert.NoError(t, err)
	})
}

func TestNewTracerProvider_SampleRatio(t *testing.T) {
	tests := []struct {
		name        string
		sampleRatio float64
	}{
		{"always sample", 1.0},
		{"never sample", 0.0},
		{"half sample", 0.5},
		{"tenth sample", 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := TracingConfig{
				ExporterType: TracingExporterStdout,
				SampleRatio:  tt.sampleRatio,
			}

			tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
			require.NoError(t, err)
			require.NotNil(t, tp)

			err = tp.Shutdown(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestNewTracerProvider_InvalidExporter(t *testing.T) {
	cfg := TracingConfig{
		ExporterType: "invalid",
		SampleRatio:  1.0,
	}

	tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
	assert.Error(t, err)
	assert.Nil(t, tp)
	assert.Contains(t, err.Error(), "unknown tracing exporter type")
}

func TestNewTracerProvider_ServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
	}{
		{"simple name", "myservice"},
		{"hyphenated name", "my-service"},
		{"underscored name", "my_service"},
		{"with numbers", "service-v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := TracingConfig{
				ExporterType: TracingExporterStdout,
				SampleRatio:  1.0,
			}

			tp, err := NewTracerProvider(context.Background(), tt.serviceName, cfg)
			require.NoError(t, err)
			require.NotNil(t, tp)

			err = tp.Shutdown(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestNewMeterProvider_Stdout(t *testing.T) {
	mp, err := NewMeterProvider("test-service", MetricsExporterStdout)
	require.NoError(t, err)
	require.NotNil(t, mp)

	err = mp.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestNewMeterProvider_Prometheus(t *testing.T) {
	mp, err := NewMeterProvider("test-service", MetricsExporterPrometheus)
	require.NoError(t, err)
	require.NotNil(t, mp)

	err = mp.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestNewMeterProvider_InvalidExporter(t *testing.T) {
	mp, err := NewMeterProvider("test-service", "invalid")
	assert.Error(t, err)
	assert.Nil(t, mp)
	assert.Contains(t, err.Error(), "unknown metrics exporter type")
}

func TestNewMeterProvider_ServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
	}{
		{"simple name", "myservice"},
		{"hyphenated name", "my-service"},
		{"underscored name", "my_service"},
		{"with numbers", "service-v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp, err := NewMeterProvider(tt.serviceName, MetricsExporterStdout)
			require.NoError(t, err)
			require.NotNil(t, mp)

			err = mp.Shutdown(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestInitMetrics(t *testing.T) {
	// Setup meter provider first
	mp, err := NewMeterProvider("test-service", MetricsExporterStdout)
	require.NoError(t, err)
	require.NotNil(t, mp)
	defer func() {
		_ = mp.Shutdown(context.Background())
	}()

	m, err := InitMetrics(context.Background(), "test-service")
	require.NoError(t, err)
	require.NotNil(t, m)

	// Verify all API metrics are initialized
	assert.NotNil(t, m.APIRequestsTotal)
	assert.NotNil(t, m.APIRequestDuration)
	assert.NotNil(t, m.APIRequestsInFlight)
	assert.NotNil(t, m.APIErrorsTotal)

	// Verify all DE execution metrics are initialized
	assert.NotNil(t, m.DEExecutionsTotal)
	assert.NotNil(t, m.DEExecutionDuration)
	assert.NotNil(t, m.DEExecutionsInFlight)
	assert.NotNil(t, m.DEGenerationsTotal)
	assert.NotNil(t, m.ParetoSetSize)

	// Verify executor worker pool metrics are initialized
	assert.NotNil(t, m.ExecutorWorkersActive)
	assert.NotNil(t, m.ExecutorWorkersTotal)
	assert.NotNil(t, m.ExecutorQueueWaitDuration)
	assert.NotNil(t, m.ExecutorUtilizationPercent)

	// Verify auth metrics are initialized
	assert.NotNil(t, m.AuthAttemptsTotal)
	assert.NotNil(t, m.AuthSuccessTotal)
	assert.NotNil(t, m.AuthFailuresTotal)

	// Verify rate limiting metrics are initialized
	assert.NotNil(t, m.RateLimitExceeded)

	// Verify panic metrics are initialized
	assert.NotNil(t, m.PanicsTotal)
}

func TestNewTracerProvider_File(t *testing.T) {
	t.Run("writes traces to file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "traces.json")

		cfg := TracingConfig{
			ExporterType: TracingExporterFile,
			SampleRatio:  1.0,
			FilePath:     filePath,
		}

		tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
		require.NoError(t, err)
		require.NotNil(t, tp)

		err = tp.Shutdown(context.Background())
		assert.NoError(t, err)

		_, err = os.Stat(filePath)
		assert.NoError(t, err, "trace file should exist")
	})

	t.Run("uses default filename when empty", func(t *testing.T) {
		// Change to temp dir to avoid polluting working directory
		origDir, _ := os.Getwd()
		tmpDir := t.TempDir()
		_ = os.Chdir(tmpDir)
		defer func() { _ = os.Chdir(origDir) }()

		cfg := TracingConfig{
			ExporterType: TracingExporterFile,
			SampleRatio:  1.0,
			FilePath:     "",
		}

		tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
		require.NoError(t, err)
		require.NotNil(t, tp)

		err = tp.Shutdown(context.Background())
		assert.NoError(t, err)

		_, err = os.Stat(filepath.Join(tmpDir, "traces.json"))
		assert.NoError(t, err, "default traces.json should be created")
	})

	t.Run("returns error for invalid file path", func(t *testing.T) {
		cfg := TracingConfig{
			ExporterType: TracingExporterFile,
			SampleRatio:  1.0,
			FilePath:     "/nonexistent/dir/traces.json",
		}

		tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
		assert.Error(t, err)
		assert.Nil(t, tp)
		assert.Contains(t, err.Error(), "failed to open trace file")
	})
}

func TestNewTracerProvider_None(t *testing.T) {
	tests := []struct {
		name        string
		sampleRatio float64
	}{
		{"always sample", 1.0},
		{"never sample", 0.0},
		{"ratio sample", 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := TracingConfig{
				ExporterType: TracingExporterNone,
				SampleRatio:  tt.sampleRatio,
			}

			tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
			require.NoError(t, err)
			require.NotNil(t, tp)

			err = tp.Shutdown(context.Background())
			assert.NoError(t, err)
		})
	}
}

// generateTestCert creates a self-signed cert and key in tmpDir, returning their paths.
func generateTestCert(t *testing.T, tmpDir string) (certPath, keyPath string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{Organization: []string{"Test"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		IsCA:         true,
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	certPath = filepath.Join(tmpDir, "cert.pem")
	certFile, err := os.Create(certPath)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}))
	certFile.Close()

	keyDER, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)

	keyPath = filepath.Join(tmpDir, "key.pem")
	keyFile, err := os.Create(keyPath)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))
	keyFile.Close()

	return certPath, keyPath
}

func TestBuildTLSCredentials(t *testing.T) {
	t.Run("minimal config with skip verify", func(t *testing.T) {
		creds, err := buildTLSCredentials(TLSConfig{
			Enabled:    true,
			SkipVerify: true,
		})
		require.NoError(t, err)
		assert.NotNil(t, creds)
	})

	t.Run("with client certificate", func(t *testing.T) {
		tmpDir := t.TempDir()
		certPath, keyPath := generateTestCert(t, tmpDir)

		creds, err := buildTLSCredentials(TLSConfig{
			Enabled:    true,
			CertFile:   certPath,
			KeyFile:    keyPath,
			SkipVerify: true,
		})
		require.NoError(t, err)
		assert.NotNil(t, creds)
	})

	t.Run("with CA certificate", func(t *testing.T) {
		tmpDir := t.TempDir()
		certPath, _ := generateTestCert(t, tmpDir)

		creds, err := buildTLSCredentials(TLSConfig{
			Enabled:    true,
			CAFile:     certPath, // self-signed cert also works as CA
			SkipVerify: true,
		})
		require.NoError(t, err)
		assert.NotNil(t, creds)
	})

	t.Run("returns error for invalid client cert", func(t *testing.T) {
		_, err := buildTLSCredentials(TLSConfig{
			Enabled:  true,
			CertFile: "/nonexistent/cert.pem",
			KeyFile:  "/nonexistent/key.pem",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load client certificate")
	})

	t.Run("returns error for invalid CA file", func(t *testing.T) {
		_, err := buildTLSCredentials(TLSConfig{
			Enabled: true,
			CAFile:  "/nonexistent/ca.pem",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read CA certificate")
	})

	t.Run("returns error for unparseable CA", func(t *testing.T) {
		tmpDir := t.TempDir()
		badCAPath := filepath.Join(tmpDir, "bad-ca.pem")
		require.NoError(t, os.WriteFile(badCAPath, []byte("not a certificate"), 0600))

		_, err := buildTLSCredentials(TLSConfig{
			Enabled: true,
			CAFile:  badCAPath,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse CA certificate")
	})
}

func TestInitMetrics_DEAlgorithmMetrics(t *testing.T) {
	mp, err := NewMeterProvider("test-service", MetricsExporterStdout)
	require.NoError(t, err)
	defer func() { _ = mp.Shutdown(context.Background()) }()

	m, err := InitMetrics(context.Background(), "test-service")
	require.NoError(t, err)

	// Verify DE algorithm-specific metrics
	assert.NotNil(t, m.DEObjectiveEvaluations)
	assert.NotNil(t, m.DEMutationsTotal)
	assert.NotNil(t, m.DECrossoverTotal)
	assert.NotNil(t, m.DEPopulationDiversity)
	assert.NotNil(t, m.DEConvergenceRate)
	assert.NotNil(t, m.DENonDominatedCount)
	assert.NotNil(t, m.DERankZeroSize)
	assert.NotNil(t, m.DECrowdingDistanceAvg)
	assert.NotNil(t, m.DEVariantPerformance)
	assert.NotNil(t, m.DEProblemComplexity)
}

func TestInitMetrics_WithPrometheus(t *testing.T) {
	// Setup Prometheus meter provider
	mp, err := NewMeterProvider("test-service", MetricsExporterPrometheus)
	require.NoError(t, err)
	require.NotNil(t, mp)
	defer func() {
		_ = mp.Shutdown(context.Background())
	}()

	m, err := InitMetrics(context.Background(), "test-service")
	require.NoError(t, err)
	require.NotNil(t, m)

	// Verify metrics can be recorded (basic smoke test)
	ctx := context.Background()
	m.APIRequestsTotal.Add(ctx, 1)
	m.DEExecutionsTotal.Add(ctx, 1)
	m.AuthAttemptsTotal.Add(ctx, 1)
}
