package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/api"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repo"
	"github.com/arttet/reddit-feed-api/internal/broker"
	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/telemetry"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc_otel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
)

type Server struct {
	producer   broker.Producer
	repository repo.Repo
	logger     *zap.Logger
}

func NewServer(
	producer broker.Producer,
	repository repo.Repo,
	logger *zap.Logger,
) *Server {

	return &Server{
		producer:   producer,
		repository: repository,
		logger:     logger,
	}
}

func (s *Server) Start(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)
	statusAdrr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.REST.Host, cfg.REST.Port)
	grpcAddr := fmt.Sprintf("%s:%v", cfg.GRPC.Host, cfg.GRPC.Port)

	logger := s.logger

	/**
	 *  OpenTelemetry Tracer
	 **/

	tp, err := telemetry.NewTracer(cfg)
	if err != nil {
		logger.Error("tracing initialization", zap.Error(err))
		return err
	}

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("tracer shut down", zap.Error(err))
		} else {
			logger.Info("tracer shut down correctly")
		}
	}()

	logger.Info("tracer is running", zap.String("address", cfg.Jaeger.URL))

	/**
	 * Metric Server
	 **/

	metricsServer := telemetry.CreateMetricsServer(cfg)

	go func() {
		logger.Info("metrics server is running", zap.String("address", metricAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed running metrics server", zap.Error(err))
			cancel()
		}
	}()

	/**
	 * Status Server
	 **/

	isReady := &atomic.Value{}
	isReady.Store(false)
	statusServer := telemetry.NewStatusServer(cfg, isReady)

	go func() {
		logger.Info("status server is running", zap.String("address", statusAdrr))
		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed running status server", zap.Error(err))
			// cancel()
		}
	}()

	/**
	 * Gateway
	 **/

	gatewayServer := createGatewayServer(grpcAddr, gatewayAddr)

	go func() {
		logger.Info("gateway server is running", zap.String("address", gatewayAddr))
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed running gateway server", zap.Error(err))
			cancel()
		}
	}()

	/**
	 * gRPC
	 **/

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.GRPC.MaxConnectionIdle,
			Timeout:           cfg.GRPC.Timeout,
			MaxConnectionAge:  cfg.GRPC.MaxConnectionAge,
			Time:              cfg.GRPC.Timeout,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_otel.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, opts...),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_otel.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger, opts...),
			grpc_recovery.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
		)),
	)

	pb.RegisterRedditFeedAPIServiceServer(
		grpcServer,
		api.NewRedditFeedAPI(s.repository, s.producer, logger),
	)

	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	go func() {
		logger.Info("gRPC server is running", zap.String("address", grpcAddr))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("failed running gRPC server", zap.Error(err))
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		logger.Info("the service is ready to accept requests")
	}()

	if cfg.Project.Debug {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		logger.Info("", zap.String("signal.Notify", fmt.Sprintf("%v", v)))
	case done := <-ctx.Done():
		logger.Info("", zap.String("ctx.Done", fmt.Sprintf("%v", done)))
	}

	isReady.Store(false)

	if err := gatewayServer.Shutdown(ctx); err != nil {
		logger.Error("gateway server shut down", zap.Error(err))
	} else {
		logger.Info("gateway server shut down correctly")
	}

	if err := statusServer.Shutdown(ctx); err != nil {
		logger.Error("status server shut down", zap.Error(err))
	} else {
		logger.Info("status server shut down correctly")
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Error("metrics server shut down", zap.Error(err))
	} else {
		logger.Info("metrics server shut down correctly")
	}

	grpcServer.GracefulStop()
	logger.Info("gRPC server shut down correctly")

	return nil
}
