package telemetry

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
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

func (p *provider) init() {
	closeFunc, err := configureTracer(p.ctx)
	if err != nil {
		panic(err)
	}

	p.Close = func() {
		if err := closeFunc(p.ctx); err != nil {
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
		sdktrace.WithResource(newResource()),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tracerProvider.Shutdown, nil

}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("bills-api"),
		semconv.ServiceVersionKey.String("0.0.1"),
	)
}

func Init(ctx context.Context) func() {
	telemetryProvider := provider{
		ctx: ctx,
	}

	telemetryProvider.init()
	return telemetryProvider.Close
}
