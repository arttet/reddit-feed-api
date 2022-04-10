package feed

import (
	"context"
	"sort"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repository"
	"github.com/arttet/reddit-feed-api/internal/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

type Feed interface {
	CreatePosts(ctx context.Context, posts model.Posts) (int64, error)
	GenerateFeed(ctx context.Context, pageID uint64) (model.Posts, error)
	FilterPosts(ctx context.Context, posts model.Posts) (model.Posts, error)
}

// NewFeed creates a new instance of Feed Service.
func NewFeed(repo repository.Repository, logger *zap.Logger) Feed {
	return &feed{
		repo:   repo,
		logger: logger,
	}
}

const maximumPosts = 27

var (
	tracer trace.Tracer
)

func init() {
	tracer = otel.Tracer("github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/feed")
}

type feed struct {
	repo   repository.Repository
	logger *zap.Logger
}

func (f *feed) CreatePosts(ctx context.Context, posts model.Posts) (int64, error) {
	return f.repo.CreatePosts(ctx, posts)
}

func (f *feed) GenerateFeed(ctx context.Context, pageID uint64) (model.Posts, error) {
	ctx, span := tracer.Start(ctx, "GenerateFeed")
	defer span.End()

	offset := maximumPosts * (pageID - 1)
	posts, err := f.repo.ListPosts(ctx, maximumPosts, offset)
	if err != nil {
		return nil, err
	}

	comparator := func(p, q int) bool {
		return posts[p].Score > posts[q].Score
	}

	if !sort.SliceIsSorted(posts, comparator) {
		f.logger.Warn("posts should be sorted")
		sort.SliceStable(posts, comparator)
	}

	posts, err = f.FilterPosts(ctx, posts)
	return posts, err
}

func (f *feed) FilterPosts(ctx context.Context, posts model.Posts) (model.Posts, error) {
	size := len(posts)
	if size < 3 {
		return posts, nil
	}

	promotedPost, err := f.repo.GetPromotedPost(ctx)
	if promotedPost == nil || err != nil {
		f.logger.Warn("failed to find a promoted post")
	}

	filteredPost := make(model.Posts, 0, size)
	for i, nextPost := range posts {
		n := len(filteredPost)
		if n == maximumPosts {
			break
		}

		if nth := n + 1; promotedPost != nil && (nth == 2 && size >= 3 || nth == 16 && size > 16) {
			previousPost := posts[i-1]
			if !previousPost.NotSafeForWork && !nextPost.NotSafeForWork {
				filteredPost = append(filteredPost, promotedPost)
				size++
			}
		}

		filteredPost = append(filteredPost, nextPost)
	}

	return filteredPost, nil
}
