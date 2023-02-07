package telemetry

import (
	"context"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"time"
)

type provider struct {
	ctx       context.Context
	closeFunc func(context.Context) error
	Close     func()
}

func (p *provider) configure() {
	shutdownTracer, err := configureTracer(p.ctx)
	if err != nil {
		panic(err)
	}

	shutdownMeterProvider, err := configureMetrics(p.ctx)
	if err != nil {
		panic(err)
	}

	p.Close = func() {
		if err := shutdownTracer(p.ctx); err != nil {
			panic(err)
		}

		if err := shutdownMeterProvider(p.ctx); err != nil {
			panic(err)
		}
	}
}

func configureTracer(ctx context.Context) (func(context.Context) error, error) {

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))

	sCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	exporter, err := otlptrace.New(sCtx, client)
	if err != nil {
		panic(fmt.Errorf("creating OTLP trace exporter: %w", err))
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("personal-fin-api"),
			semconv.ServiceVersionKey.String("0.0.1.Alpha"),
			attribute.String("library.language", "go"),
		)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider.Shutdown, nil

}

func configureMetrics(ctx context.Context) (func(ctx context.Context) error, error) {
	exporter, eErr := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if eErr != nil {
		return nil, eErr
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(time.Second*1))))
	global.SetMeterProvider(meterProvider)

	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		return nil, err
	}

	return meterProvider.Shutdown, nil
}

func Init(ctx context.Context) func() {
	telemetryProvider := provider{
		ctx: ctx,
	}

	telemetryProvider.configure()
	return telemetryProvider.Close
}
