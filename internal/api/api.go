package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/arttet/reddit-feed-api/internal/model"
	"github.com/arttet/reddit-feed-api/internal/repo"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api"
)

type api struct {
	pb.UnimplementedRedditFeedAPIServiceServer
	repo          repo.Repo
	chunkSize     uint
	maxCountPosts uint64
}

func NewRedditFeedAPI(r repo.Repo, chunkSize uint) pb.RedditFeedAPIServiceServer {
	return &api{
		repo:          r,
		chunkSize:     chunkSize,
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

	if err := request.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
				posts[j].Link = request.Posts[i].GetLink()
				posts[j].Subreddit = request.Posts[i].Subreddit
				posts[j].Content = request.Posts[i].GetContent()
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
			if errors.Is(err, sql.ErrConnDone) {
				return nil, status.Error(codes.Unavailable, err.Error())
			}
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

	posts, err := a.repo.ListPosts(ctx, a.maxCountPosts, a.maxCountPosts*(request.PageId-1))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		log.Error().Err(err).Msg("Failed to get the data")
		if errors.Is(err, sql.ErrConnDone) {
			return nil, status.Error(codes.Unavailable, err.Error())
		}
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	response := &pb.GenerateFeedV1Response{
		Posts: a.filterPosts(ctx, posts),
	}

	return response, nil
}

func (a *api) filterPosts(
	ctx context.Context,
	posts []model.Post,
) (
	list []*pb.Post,
) {

	n := len(posts)

	var err error
	var promotedPost *model.Post

	if n >= 3 {
		promotedPost, err = a.repo.PromotedPost(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to find a promoted post")
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
