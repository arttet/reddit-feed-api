package tracer

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"

	"github.com/arttet/reddit-feed-api/internal/config"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// NewTracer returns a new tracer.
func NewTracer(cfg *config.Config) (io.Closer, error) {
	tracerAddr := fmt.Sprintf("%s:%v", cfg.Jaeger.Host, cfg.Jaeger.Port)

	cfgTracer := &jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.Service,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: tracerAddr,
		},
	}

	logger := zap.L()

	tracer, closer, err := cfgTracer.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		logger.Fatal("failed to init Jaeger", zap.Error(err))
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	logger.Info("tracer started")

	return closer, nil
}
