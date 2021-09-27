package api

import (
	"context"

	"github.com/arttet/reddit-feed-api/internal/broker"
	"github.com/arttet/reddit-feed-api/internal/model"
	"github.com/arttet/reddit-feed-api/internal/repo"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"go.uber.org/zap"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api"
)

var (
	totalFeedNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "reddit_feed_api_feed_not_found_total",
		Help: "The total number of feeds that were not found",
	})
)

type api struct {
	pb.UnimplementedRedditFeedAPIServiceServer
	repository    repo.Repo
	producer      broker.Producer
	logger        *zap.Logger
	maxCountPosts uint64
}

func NewRedditFeedAPI(
	r repo.Repo,
	p broker.Producer,
	logger *zap.Logger,
) pb.RedditFeedAPIServiceServer {

	return &api{
		repository:    r,
		producer:      p,
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

	span, ctx := opentracing.StartSpanFromContext(ctx, "CreatePostsV1")
	defer span.Finish()
	span.LogFields(
		log.Int("len", len(request.Posts)),
	)

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var numberOfCreatedPosts int64
	posts := make([]model.Post, 0, len(request.Posts))
	for i := range request.Posts {
		post := model.Post{
			Title:          request.Posts[i].Title,
			Author:         request.Posts[i].Author,
			Link:           request.Posts[i].GetLink(),
			Subreddit:      request.Posts[i].Subreddit,
			Content:        request.Posts[i].GetContent(),
			Score:          request.Posts[i].Score,
			Promoted:       request.Posts[i].Promoted,
			NotSafeForWork: request.Posts[i].NotSafeForWork,
		}
		posts = append(posts, post)
	}

	numberOfCreatedPosts, err := a.repository.CreatePosts(ctx, span, posts)
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

	span, ctx := opentracing.StartSpanFromContext(ctx, "GenerateFeedV1")
	defer span.Finish()
	span.LogFields(
		log.Uint64("PageId", request.PageId),
	)

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	posts, err := a.repository.ListPosts(ctx, span, a.maxCountPosts, a.maxCountPosts*(request.PageId-1))
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
		Posts: a.filterPosts(ctx, span, posts),
	}

	return response, nil
}

func (a *api) filterPosts(
	ctx context.Context,
	parentSpan opentracing.Span,
	posts []model.Post,
) (
	list []*pb.Post,
) {

	span := opentracing.StartSpan(
		"filterPosts",
		opentracing.ChildOf(parentSpan.Context()))
	defer span.Finish()

	n := len(posts)

	var err error
	var promotedPost *model.Post

	if n >= 3 {
		promotedPost, err = a.repository.GetPromotedPost(ctx, span)
		if err != nil {
			a.logger.Warn("failed to find a promoted post")
		}
	}

	list = make([]*pb.Post, 0, n)
	counter := 0

	pushBack := func(post model.Post) {
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
			pushBack(*promotedPost)
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
