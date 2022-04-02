package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/feed"
	"github.com/arttet/reddit-feed-api/internal/broker"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
	"github.com/arttet/reddit-feed-api/pkg/transform"
)

// NewRedditFeedAPIServiceServer creates a new instance of Reddit Feed API Service Server.
func NewRedditFeedAPIServiceServer(
	srv feed.Feed,
	producer broker.Producer,
	logger *zap.Logger,
) pb.RedditFeedAPIServiceServer {

	return &api{
		srv:      srv,
		producer: producer,
		logger:   logger,
	}
}

var (
	totalFeedNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "reddit_feed_api_feed_not_found_total",
		Help: "The total number of feeds that were not found",
	})
)

type api struct {
	pb.UnimplementedRedditFeedAPIServiceServer
	srv      feed.Feed
	producer broker.Producer
	logger   *zap.Logger
}

func (a *api) CreatePostsV1(
	ctx context.Context,
	request *pb.CreatePostsV1Request,
) (
	*pb.CreatePostsV1Response,
	error,
) {

	posts := transform.PbToPostPtrList(request.Posts)

	numberOfCreatedPosts, err := a.srv.CreatePosts(ctx, posts)
	if err != nil {
		a.logger.Error("failed to insert the data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	a.producer.CreatePosts(posts)

	response := &pb.CreatePostsV1Response{
		NumberOfCreatedPosts: numberOfCreatedPosts,
	}

	return response, nil
}

func (a *api) GenerateFeedV1(
	ctx context.Context,
	request *pb.GenerateFeedV1Request,
) (
	*pb.GenerateFeedV1Response,
	error,
) {

	posts, err := a.srv.GenerateFeed(ctx, request.PageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.Warn("a feed not found")
			totalFeedNotFound.Inc()
			return nil, status.Error(codes.NotFound, "a feed not found")
		}

		a.logger.Error("failed to list the data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.GenerateFeedV1Response{
		Posts: transform.PostToPbPtrList(posts),
	}

	return response, nil
}
