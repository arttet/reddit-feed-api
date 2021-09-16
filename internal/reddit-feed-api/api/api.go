package api

import (
	"context"

	"github.com/arttet/reddit-feed-api/internal/reddit-feed-api/model"
	"github.com/arttet/reddit-feed-api/internal/reddit-feed-api/repo"
	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type api struct {
	pb.UnimplementedRedditFeedAPIServiceServer
	repo      repo.Repo
	chunkSize uint
}

func NewRedditFeedAPI(r repo.Repo, chunkSize uint) pb.RedditFeedAPIServiceServer {
	return &api{
		repo:      r,
		chunkSize: chunkSize,
	}
}

func (a *api) CreatePostsV1(
	ctx context.Context,
	request *pb.CreatePostsV1Request,
) (
	*pb.CreatePostsV1Response,
	error,
) {

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	for i := range request.Posts {
		emptyLink := request.Posts[i].Link == ""
		emptyContent := request.Posts[i].Content == ""

		if (emptyLink && emptyContent) || (!emptyLink && !emptyContent) {
			return nil, status.Error(codes.InvalidArgument, "a post cannot have both a link and content populated")
		}
	}

	var numberOfCreatedPosts int64
	posts := make([]model.Post, a.chunkSize)

	for i, n := uint(0), uint(len(request.Posts)); i < n; i += a.chunkSize {
		end := i + a.chunkSize
		if end > n {
			end = n
		}

		number, err := func() (int64, error) {
			var j int
			for i < end {
				posts[j].Title = request.Posts[i].Title
				posts[j].Author = request.Posts[i].Author
				posts[j].Link = request.Posts[i].Link
				posts[j].Subreddit = request.Posts[i].Subreddit
				posts[j].Content = request.Posts[i].Content
				posts[j].Score = request.Posts[i].Score
				posts[j].Promoted = request.Posts[i].Promoted
				posts[j].NotSafeForWork = request.Posts[i].NotSafeForWork

				i++
				j++
			}

			return a.repo.CreatePosts(ctx, posts[:j])
		}()

		if err != nil {
			log.Error().Err(err).Msg("Failed to insert the data")
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		}

		numberOfCreatedPosts += number
	}

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

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}