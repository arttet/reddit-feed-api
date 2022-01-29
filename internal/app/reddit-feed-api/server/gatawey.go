package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
)

var (
	httpTotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_microservice_requests_total",
		Help: "The total number of incoming HTTP requests",
	})

	grpcGatewayTag = opentracing.Tag{Key: string(ext.Component), Value: "grpc-gateway"}
)

func createGatewayServer(grpcAddr, gatewayAddr string) *http.Server {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		grpcAddr,
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.L().Fatal("failed to dial server", zap.Error(err))
	}

	mux := runtime.NewServeMux()
	if err := pb.RegisterRedditFeedAPIServiceHandler(context.Background(), mux, conn); err != nil {
		zap.L().Fatal("failed registration handler", zap.Error(err))
	}

	gatewayServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: tracingWrapper(mux),
	}

	return gatewayServer
}

func tracingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpTotalRequests.Inc()
		parentSpanContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil || errors.Is(err, opentracing.ErrSpanContextNotFound) {
			serverSpan := opentracing.GlobalTracer().StartSpan(
				"ServeHTTP",
				ext.RPCServerOption(parentSpanContext),
				grpcGatewayTag,
			)
			r = r.WithContext(opentracing.ContextWithSpan(r.Context(), serverSpan))
			defer serverSpan.Finish()
		}
		h.ServeHTTP(w, r)
	})
}
