package tracer

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go"

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

	tracer, closer, err := cfgTracer.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		log.Err(err).Msgf("Failed to init Jaeger: %v", err)
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	log.Info().Msg("Traces started")

	return closer, nil
}
