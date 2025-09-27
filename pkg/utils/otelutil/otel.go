package otelutil

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/grpc/credentials"
)

type Options struct {
	ServiceName            string
	TraceExporterType      string
	TraceExporterEndpoint  string
	TraceExportTimeout     time.Duration
	MetricExporterType     string
	MetricExporterEndpoint string
	MetricExportInterval   time.Duration
	LoggerExporterType     string
	LoggerExporterEndpoint string
	LoggerExportTimeout    time.Duration
	TlsCertPath            string
}

const (
	ExporterTypeGRPC    = "grpc"
	ExporterTypeConsole = "console"
)

var options Options
var tlsCreds credentials.TransportCredentials

func SetupOTelSDK(ctx context.Context, opts Options) (shutdown func(context.Context) error, err error) {
	options = opts
	if options.TraceExportTimeout == 0 {
		options.TraceExportTimeout = time.Minute
	}
	if options.MetricExportInterval == 0 {
		options.MetricExportInterval = time.Minute
	}
	if options.LoggerExportTimeout == 0 {
		options.LoggerExportTimeout = time.Minute
	}
	if strings.TrimSpace(opts.TlsCertPath) != "" {
		err = InitTLSConfig(opts)
		if err != nil {
			return nil, err
		}
	}
	return setupOTelSDK(ctx)
}

func InitTLSConfig(opts Options) error {
	// 加载 CA 证书，用于验证服务器证书
	caCert, err := os.ReadFile(opts.TlsCertPath)
	if err != nil {
		return err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return errors.New("failed to append CA cert")
	}

	// 创建 TLS 凭证
	tlsCreds = credentials.NewTLS(&tls.Config{
		RootCAs: certPool,
	})
	return nil
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider() (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter
	var err error
	switch options.TraceExporterType {
	case ExporterTypeGRPC:
		{
			var grpcOpts []otlptracegrpc.Option
			grpcOpts = append(grpcOpts, otlptracegrpc.WithEndpoint(options.TraceExporterEndpoint))
			if tlsCreds != nil {
				grpcOpts = append(grpcOpts, otlptracegrpc.WithTLSCredentials(tlsCreds))
			} else {
				grpcOpts = append(grpcOpts, otlptracegrpc.WithInsecure())
			}
			exporter, err = otlptracegrpc.New(context.Background(), grpcOpts...)
			if err != nil {
				return nil, err
			}
		}
	default:
		{
			exporter, err = stdouttrace.New()
			if err != nil {
				return nil, err
			}
		}
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(options.TraceExportTimeout)),
		trace.WithResource(
			resource.NewWithAttributes(
				//固定写法
				semconv.SchemaURL,
				//设置service
				semconv.ServiceNameKey.String(options.ServiceName),
				//设置Process键值对 可以让其他人员分析 全局的，设置到trace上的
				//attribute.Int("pid", os.Getpid()),
			),
		),
	)
	return traceProvider, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	var exporter metric.Exporter
	var err error
	switch options.MetricExporterType {
	case ExporterTypeGRPC:
		{
			var grpcOpts []otlpmetricgrpc.Option
			grpcOpts = append(grpcOpts, otlpmetricgrpc.WithEndpoint(options.MetricExporterEndpoint))
			if tlsCreds != nil {
				grpcOpts = append(grpcOpts, otlpmetricgrpc.WithTLSCredentials(tlsCreds))
			} else {
				grpcOpts = append(grpcOpts, otlpmetricgrpc.WithInsecure())
			}
			exporter, err = otlpmetricgrpc.New(context.Background(), grpcOpts...)
			if err != nil {
				return nil, err
			}
		}
	default:
		{
			exporter, err = stdoutmetric.New()
			if err != nil {
				return nil, err
			}
		}
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(options.MetricExportInterval))),
		metric.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(options.ServiceName),
			),
		),
	)
	return meterProvider, nil
}

func newLoggerProvider() (*log.LoggerProvider, error) {
	var exporter log.Exporter
	var err error
	switch options.LoggerExporterType {
	case ExporterTypeGRPC:
		{
			var grpcOpts []otlploggrpc.Option
			grpcOpts = append(grpcOpts, otlploggrpc.WithEndpoint(options.LoggerExporterEndpoint))
			if tlsCreds != nil {
				grpcOpts = append(grpcOpts, otlploggrpc.WithTLSCredentials(tlsCreds))
			} else {
				grpcOpts = append(grpcOpts, otlploggrpc.WithInsecure())
			}
			exporter, err = otlploggrpc.New(context.Background(), grpcOpts...)
			if err != nil {
				return nil, err
			}
		}
	default:
		{
			exporter, err = stdoutlog.New()
			if err != nil {
				return nil, err
			}
		}
	}
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter, log.WithExportTimeout(options.LoggerExportTimeout))),
		log.WithResource(
			resource.NewWithAttributes(
				//固定写法
				semconv.SchemaURL,
				//设置service
				semconv.ServiceNameKey.String(options.ServiceName),
			),
		),
	)
	return loggerProvider, nil
}
