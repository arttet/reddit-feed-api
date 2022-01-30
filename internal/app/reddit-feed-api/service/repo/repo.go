package repo

import (
	"context"
	"database/sql"

	"github.com/arttet/reddit-feed-api/internal/model"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repo")
}

type Repo interface {
	CreatePosts(ctx context.Context, posts model.Posts) (int64, error)
	ListPosts(ctx context.Context, limit uint64, offset uint64) (model.Posts, error)
	GetPromotedPost(ctx context.Context) (*model.Post, error)
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

const TableName = "post"

var InsertColumns = []string{
	"title",
	"author",
	"link",
	"subreddit",
	"content",
	"score",
	"promoted",
	"not_safe_for_work",
}

var SelectColumns = append([]string{"id"}, InsertColumns...)

type repo struct {
	db *sqlx.DB
}

func (r *repo) CreatePosts(
	ctx context.Context,
	posts model.Posts,
) (
	int64,
	error,
) {

	ctx, span := tracer.Start(ctx, "CreatePosts")
	defer span.End()

	query := squirrel.Insert(TableName).
		Columns(InsertColumns...).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db)

	for i := range posts {
		query = query.Values(
			posts[i].Title,
			posts[i].Author,
			posts[i].Link,
			posts[i].Subreddit,
			posts[i].Content,
			posts[i].Score,
			posts[i].Promoted,
			posts[i].NotSafeForWork,
		)
	}

	result, err := query.ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *repo) ListPosts(
	ctx context.Context,
	limit uint64,
	offset uint64,
) (
	model.Posts,
	error,
) {

	ctx, span := tracer.Start(ctx, "ListPosts")
	defer span.End()

	query := squirrel.Select(SelectColumns...).
		From(TableName).
		OrderBy("score DESC").
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make(model.Posts, 0, limit)
	for rows.Next() {
		post := &model.Post{}
		scanRow(rows, post)
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *repo) GetPromotedPost(
	ctx context.Context,
) (
	*model.Post,
	error,
) {

	ctx, span := tracer.Start(ctx, "GetPromotedPost")
	defer span.End()

	query := squirrel.Select(SelectColumns...).
		From(TableName).
		Where(squirrel.Eq{"promoted": true}).
		OrderBy("score DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var promotedPost model.Post
	if rows.Next() {
		scanRow(rows, &promotedPost)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &promotedPost, nil
}

func scanRow(rows *sql.Rows, post *model.Post) {
	// nolint:errcheck
	rows.Scan(
		&post.ID,
		&post.Title,
		&post.Author,
		&post.Link,
		&post.Subreddit,
		&post.Content,
		&post.Score,
		&post.Promoted,
		&post.NotSafeForWork,
	)
}
