package repository

import (
	"context"
	"database/sql"

	"github.com/arttet/reddit-feed-api/internal/database"
	"github.com/arttet/reddit-feed-api/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Repository interface {
	CreatePosts(ctx context.Context, posts model.Posts) (int64, error)
	ListPosts(ctx context.Context, limit uint64, offset uint64) (model.Posts, error)
	GetPromotedPost(ctx context.Context) (*model.Post, error)
}

// NewRepository creates a new instance of Repository Service.
func NewRepository(db *sqlx.DB) Repository {
	return &repo{
		db: db,
	}
}

const TableName = "post"

var (
	InsertColumns = []string{
		"title",
		"author",
		"link",
		"subreddit",
		"content",
		"score",
		"promoted",
		"not_safe_for_work",
	}

	SelectColumns = append([]string{"id"}, InsertColumns...)

	tracer trace.Tracer
)

func init() {
	tracer = otel.Tracer("github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repository")
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) CreatePosts(ctx context.Context, posts model.Posts) (int64, error) {
	ctx, span := tracer.Start(ctx, "CreatePosts")
	defer span.End()

	var (
		result sql.Result
		err    error
	)

	err = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.StatementBuilder.
			Insert(TableName).
			Columns(InsertColumns...).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		for _, p := range posts {
			sb = sb.Values(
				p.Title,
				p.Author,
				p.Link,
				p.Subreddit,
				p.Content,
				p.Score,
				p.Promoted,
				p.NotSafeForWork,
			)
		}

		result, err = sb.ExecContext(ctx)
		return err
	})

	if err != nil {
		return 0, errors.Wrap(err, "repo.Add")
	}

	return result.RowsAffected()
}

func (r *repo) ListPosts(ctx context.Context, limit uint64, offset uint64) (model.Posts, error) {
	ctx, span := tracer.Start(ctx, "ListPosts")
	defer span.End()

	var posts model.Posts

	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.Select(SelectColumns...).
			From(TableName).
			OrderBy("score DESC").
			Limit(limit).
			Offset(offset).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &posts, query, args...)
	})

	if err != nil {
		return nil, errors.Wrap(err, "repo.ListPosts")
	}

	if len(posts) == 0 {
		return nil, errors.Wrap(sql.ErrNoRows, "repo.ListPosts")
	}

	return posts, nil
}

func (r *repo) GetPromotedPost(ctx context.Context) (*model.Post, error) {
	ctx, span := tracer.Start(ctx, "GetPromotedPost")
	defer span.End()

	var promotedPost model.Posts

	err := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		sb := squirrel.Select(SelectColumns...).
			From(TableName).
			Where(squirrel.Eq{"promoted": true}).
			OrderBy("score DESC").
			Limit(1).
			PlaceholderFormat(squirrel.Dollar).
			RunWith(r.db)

		query, args, err := sb.ToSql()
		if err != nil {
			return err
		}

		return r.db.SelectContext(ctx, &promotedPost, query, args...)
	})

	if err != nil {
		return nil, errors.Wrap(err, "repo.GetPromotedPost")
	}

	if len(promotedPost) == 0 {
		return nil, errors.Wrap(sql.ErrNoRows, "repo.GetPromotedPost")
	}

	return promotedPost[0], nil
}
