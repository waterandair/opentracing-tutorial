package main

import (
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

	tracer, closer := initJaeger("lesson01")
	defer closer.Close()
	span := tracer.StartSpan("say-hello")
	defer span.Finish()

	helloTo := os.Args[1]
	span.SetTag("hello-to", helloTo)

	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

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
