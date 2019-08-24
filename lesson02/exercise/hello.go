package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}

	tracer, closer := initJaeger("hello-world")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	rootSpan := tracer.StartSpan("say-hello")
	defer rootSpan.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), rootSpan)

	helloTo := os.Args[1]
	rootSpan.SetTag("hello-to", helloTo)

	helloStr := formatString(ctx, helloTo)
	printHello(ctx, helloStr)
}

func formatString(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		"formatString",
	)
	defer span.Finish()

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(
		ctx,
		"printHello",
	)
	defer span.Finish()

	println(helloStr)
	span.LogKV("event", "println")
}

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	conf := &config.Configuration{
		ServiceName: "",
		Disabled:    false,
		RPCMetrics:  false,
		Tags:        nil,
		Sampler: &config.SamplerConfig{
			Type:                    "const",
			Param:                   1,
			SamplingServerURL:       "",
			MaxOperations:           0,
			SamplingRefreshInterval: 0,
		},
		Reporter: &config.ReporterConfig{
			QueueSize:           0,
			BufferFlushInterval: 0,
			LogSpans:            true,
			LocalAgentHostPort:  "",
			CollectorEndpoint:   "",
			User:                "",
			Password:            "",
		},
		Headers:             nil,
		BaggageRestrictions: nil,
		Throttler:           nil,
	}

	tracer, closer, err := conf.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}
