package api

import (
	"context"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repo"
	"github.com/arttet/reddit-feed-api/internal/broker"
	"github.com/arttet/reddit-feed-api/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
	"github.com/arttet/reddit-feed-api/pkg/transform"
)

var (
	tracer trace.Tracer

	totalFeedNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "reddit_feed_api_feed_not_found_total",
		Help: "The total number of feeds that were not found",
	})
)

func init() {
	tracer = otel.Tracer("github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/api")
}

type api struct {
	pb.UnimplementedRedditFeedAPIServiceServer
	repository    repo.Repo
	producer      broker.Producer
	logger        *zap.Logger
	maxCountPosts uint64
}

func NewRedditFeedAPI(
	repository repo.Repo,
	producer broker.Producer,
	logger *zap.Logger,
) pb.RedditFeedAPIServiceServer {

	return &api{
		repository:    repository,
		producer:      producer,
		logger:        logger,
		maxCountPosts: 27,
	}
}

func (a *api) CreatePostsV1(
	ctx context.Context,
	request *pb.CreatePostsV1Request,
) (
	*pb.CreatePostsV1Response,
	error,
) {

	posts := transform.PbToPostPtrList(request.Posts)

	numberOfCreatedPosts, err := a.repository.CreatePosts(ctx, posts)
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

	posts, err := a.repository.ListPosts(ctx, a.maxCountPosts, a.maxCountPosts*(request.PageId-1))
	if err != nil {
		a.logger.Error("failed to list the data", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(posts) == 0 {
		a.logger.Warn("a feed not found")
		totalFeedNotFound.Inc()
		return nil, status.Error(codes.NotFound, "a feed not found")
	}

	response := &pb.GenerateFeedV1Response{
		Posts: a.filterPosts(ctx, posts),
	}

	return response, nil
}

func (a *api) filterPosts(
	ctx context.Context,
	posts model.Posts,
) (
	list []*pb.Post,
) {

	ctx, span := tracer.Start(ctx, "filterPosts")
	defer span.End()

	n := len(posts)

	var err error
	var promotedPost *model.Post

	if n >= 3 {
		promotedPost, err = a.repository.GetPromotedPost(ctx)
		if err != nil {
			a.logger.Warn("failed to find a promoted post")
		}
	}

	list = make([]*pb.Post, 0, n)
	counter := 0

	pushBack := func(post *model.Post) {
		result := &pb.Post{
			Title:          post.Title,
			Author:         post.Author,
			Subreddit:      post.Subreddit,
			Score:          post.Score,
			Promoted:       post.Promoted,
			NotSafeForWork: post.NotSafeForWork,
		}
		if post.Link != "" {
			result.PostType = &pb.Post_Link{Link: post.Link}
		} else {
			result.PostType = &pb.Post_Content{Content: post.Content}
		}

		list = append(list, result)
		counter++
	}

	for i, post := range posts {
		if counter == n {
			break
		}

		if promotedPost != nil && (counter == 1 || counter == 15) && !posts[i].NotSafeForWork && !posts[i-1].NotSafeForWork {
			pushBack(promotedPost)
		}

		if post.Promoted && counter > 0 && list[counter-1].NotSafeForWork {
			continue
		}

		if post.Promoted && i < n-1 && posts[i+1].NotSafeForWork {
			continue
		}

		pushBack(post)
	}

	return list
}
