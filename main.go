package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FACorreiaa/glasses-management-platform/app"
	"github.com/FACorreiaa/glasses-management-platform/config"
	"github.com/FACorreiaa/glasses-management-platform/db"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"golang.org/x/crypto/bcrypt"

	// Tempo

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	// gRPC exporter
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" // HTTP exporter
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0" // Use an appropriate semantic convention version
	// For gRPC exporter options
)

var tracer = otel.Tracer("main")

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func handleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func hasPassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password during startup", "error", err)
		panic(err)
	}
	fmt.Printf("Hashed Password: %s\n", hash)
}

func initTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// Read Tempo endpoint and protocol from environment variables
	// Example: OTEL_EXPORTER_OTLP_ENDPOINT=tempo.internal:4317
	// Example: OTEL_EXPORTER_OTLP_PROTOCOL=grpc (or http/protobuf)
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	if endpoint == "" {
		slog.Warn("OTEL_EXPORTER_OTLP_ENDPOINT not set, tracing will be disabled.")
		// Return a no-op provider if tracing is disabled
		return sdktrace.NewTracerProvider(), nil // No-op provider
	}

	slog.Info("Initializing OTLP exporter", "endpoint", endpoint)

	// Create resource identifier
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(os.Getenv("glasses-management-platform")), // Use Fly app name if available
			semconv.ServiceVersionKey.String("1.0.0"),                               // TODO: Set your app version
			// Add other relevant resource attributes (deployment env, host, etc.)
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel resource: %w", err)
	}

	// Create OTLP Exporter based on protocol
	var exporter sdktrace.SpanExporter
	var creationErr error // Use a different name to avoid shadowing the 'err' from resource.New

	// Conditional exporter creation
	if os.Getenv("MODE") == "development" {
		slog.Info("Development mode detected, creating insecure OTLP HTTP exporter.")
		// Use '=' to assign to the outer 'exporter' and 'creationErr'
		exporter, creationErr = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithInsecure(), // Use insecure for development
		)
	} else {
		slog.Info("Production/Non-Development mode detected, creating OTLP HTTP exporter (scheme determines security).")
		// Use '=' to assign to the outer 'exporter' and 'creationErr'
		// Rely on the scheme (http:// or https://) in the endpoint URL for security
		exporter, creationErr = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
			// No WithInsecure() here
		)
	}
	if creationErr != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP trace exporter for endpoint %s: %w", endpoint, creationErr)
	}

	// Batch Span Processor
	bsp := sdktrace.NewBatchSpanProcessor(exporter,
		// Adjust batching options if needed
		sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay),
		sdktrace.WithMaxQueueSize(sdktrace.DefaultMaxQueueSize),
	)

	// Create Tracer Provider
	// Consider using ParentBased(AlwaysSample()) or TraceIDRatioBased() for production sampling
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Sample all traces for now
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Set the global TracerProvider
	otel.SetTracerProvider(tp)

	// Set the global TextMapPropagator to use W3C Trace Context
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	slog.Info("OpenTelemetry TracerProvider initialized successfully.")
	return tp, nil
}

func run(ctx context.Context) error {
	hasPassword(os.Getenv("ADMIN_PASSWORD"))
	cfg, err := config.NewConfig()

	if err != nil {
		log.Printf("Failed to load config: %v", err)

		return err
	}

	//pprofAddr := os.Getenv("PPROF_ADDR")
	//pprofPort := os.Getenv("PPROF_PORT")
	var logHandler slog.Handler

	logHandlerOptions := slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.Log.Level,
	}
	if cfg.Log.Format == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, &logHandlerOptions)
		log.Println("Configured JSON logging for stdout.") // Use standard log before slog is default
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &logHandlerOptions)
		log.Println("Configured Text logging for stdout.") // Use standard log before slog is default
		// Warn if not JSON, as it's preferred for Loki
		log.Println("WARNING: Text logging format selected. JSON format is recommended for Loki integration.")
	}
	slog.SetDefault(slog.New(logHandler))
	slog.Info("Structured logger initialized", "format", cfg.Log.Format, "level", cfg.Log.Level.String())

	tp, err := initTracerProvider(ctx)
	if err != nil {
		// Log error and decide whether to continue without tracing or exit
		slog.Error("Failed to initialize OpenTelemetry TracerProvider", "error", err)
		// return err // Option: Exit if tracing is critical
		slog.Warn("Continuing without distributed tracing enabled.")
	}

	defer func() {
		if tp != nil { // Check if tp was successfully initialized
			slog.Info("Shutting down OpenTelemetry TracerProvider...")
			// Use a background context with timeout for shutdown
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := tp.Shutdown(shutdownCtx); err != nil {
				slog.Error("Error shutting down OpenTelemetry TracerProvider", "error", err)
			} else {
				slog.Info("OpenTelemetry TracerProvider shut down successfully.")
			}
		}
	}()
	pool, err := db.Init()
	if err != nil {
		log.Println(err)
	}
	defer pool.Close()

	db.WaitForDB(pool)

	if err = db.Migrate(pool); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)

	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewBuildInfoCollector())

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"code", "method"}, // Labels: HTTP status code, HTTP method
	)
	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets, // Default buckets
		},
		[]string{"code", "method"}, // Labels: HTTP status code, HTTP method
	)

	reg.MustRegister(httpRequestsTotal)
	reg.MustRegister(httpRequestDuration)
	slog.Info("Registered Prometheus HTTP request metrics with custom registry.")

	baseRouter := app.Router(pool, []byte(cfg.Server.SessionKey))

	otelHandler := otelhttp.NewHandler(baseRouter, "http.server", // Span name prefix
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	)

	slog.Info("Exposing /metrics endpoint for Prometheus.")

	// 4. Wrap the OTel handler with Prometheus middleware
	//    Order: OTel -> Prometheus counter -> Prometheus duration
	promCounterHandler := promhttp.InstrumentHandlerCounter(httpRequestsTotal, otelHandler)
	promDurationHandler := promhttp.InstrumentHandlerDuration(httpRequestDuration, promCounterHandler)

	mux := http.NewServeMux()
	// Attach the /metrics handler (using the custom registry `reg`)
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	slog.Info("Exposing /metrics endpoint for Prometheus.")

	// Attach the instrumented application handler chain to the root
	mux.Handle("/", promDurationHandler)
	slog.Info("Mounted instrumented application handler at /.")

	// Server
	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      mux,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Starting main HTTP server", "address", cfg.Server.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// --- Graceful Shutdown ---
	// Wait for interrupt signal or error
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err)
			return err // Return error to main
		}
		slog.Info("HTTP server closed.")
	case <-ctx.Done():
		slog.Info("Shutdown signal received.")
		// Trigger graceful shutdown
		ctxShutdown, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
		defer cancel()

		slog.Info("Attempting graceful shutdown of HTTP server...", "timeout", cfg.Server.GracefulTimeout)
		if err := srv.Shutdown(ctxShutdown); err != nil {
			handleError(err, "Graceful server shutdown failed")
			// Force close if graceful shutdown fails
			if closeErr := srv.Close(); closeErr != nil {
				handleError(closeErr, "Force server close failed")
			}
		} else {
			slog.Info("HTTP server shutdown gracefully.")
		}
	}

	slog.Info("Application shut down complete.")
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill) // Handle Interrupt and Kill
	slog.Info("Application starting...")                            // Initial log message

	if err := run(ctx); err != nil {
		slog.Error("Application run failed", "error", err)
		// Use fmt.Fprintf for final stderr output before exit
		fmt.Fprintf(os.Stderr, "Application run failed: %s\n", err)
		cancel() // Ensure context is cancelled on error exit
		os.Exit(1)
	}

	slog.Info("Application exited successfully.")
}
