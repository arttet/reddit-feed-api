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

	"github.com/arttet/reddit-feed-api/internal/api"
	"github.com/arttet/reddit-feed-api/internal/broker"
	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/repo"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api"
)

type Server struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func NewServer(logger *zap.Logger, db *sqlx.DB) *Server {
	return &Server{
		logger: logger,
		db:     db,
	}
}

func (s *Server) Start(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gatewayAddr := fmt.Sprintf("%s:%v", cfg.REST.Host, cfg.REST.Port)
	grpcAddr := fmt.Sprintf("%s:%v", cfg.GRPC.Host, cfg.GRPC.Port)
	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)
	statusAdrr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)

	gatewayServer := createGatewayServer(grpcAddr, gatewayAddr)

	logger := s.logger

	go func() {
		logger.Info("gateway server is running", zap.String("address", gatewayAddr))
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed running gateway server", zap.Error(err))
			cancel()
		}
	}()

	metricsServer := createMetricsServer(cfg)

	go func() {
		logger.Info("metrics server is running", zap.String("address", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed running metrics server", zap.Error(err))
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)
	statusServer := createStatusServer(cfg, isReady)

	go func() {
		logger.Info("status server is running", zap.String("address", statusAdrr))
		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed running status server", zap.Error(err))
			// cancel()
		}
	}()

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.GRPC.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.GRPC.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.GRPC.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.GRPC.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
		)),
	)

	producer, err := broker.NewProducer(ctx, &cfg.Kafka, logger)
	if err != nil {
		return fmt.Errorf("failed to create a producer: %w", err)
	}
	logger.Info("the Kafka producer is running", zap.Strings("brokers", cfg.Kafka.Brokers))

	repository := repo.NewRepo(s.db)
	pb.RegisterRedditFeedAPIServiceServer(
		grpcServer,
		api.NewRedditFeedAPI(repository, producer, logger),
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
