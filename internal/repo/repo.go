package repo

import (
	"context"
	"database/sql"

	"github.com/arttet/reddit-feed-api/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type Repo interface {
	CreatePosts(ctx context.Context, posts []model.Post) (int64, error)
	ListPosts(ctx context.Context, limit uint64, offset uint64) ([]model.Post, error)
	PromotedPost(ctx context.Context) (*model.Post, error)
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

const tableName = "post"

var insertColumns = []string{
	"title",
	"author",
	"link",
	"subreddit",
	"content",
	"score",
	"promoted",
	"not_safe_for_work",
}

var selectColumns = append([]string{"id"}, insertColumns...)

type repo struct {
	db *sqlx.DB
}

func (r *repo) CreatePosts(ctx context.Context, posts []model.Post) (int64, error) {
	query := squirrel.Insert(tableName).
		Columns(insertColumns...).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

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

func (r *repo) ListPosts(ctx context.Context, limit uint64, offset uint64) ([]model.Post, error) {
	query := squirrel.Select(selectColumns...).
		From(tableName).
		OrderBy("score DESC").
		Limit(limit).
		Offset(offset).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0, limit)
	for rows.Next() {
		var post model.Post
		if err := readRow(rows, &post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return nil, sql.ErrNoRows
	}

	return posts, nil
}

func (r *repo) PromotedPost(ctx context.Context) (*model.Post, error) {
	query := squirrel.Select(selectColumns...).
		From(tableName).
		Where(squirrel.Eq{"promoted": true}).
		RunWith(r.db).
		OrderBy("score DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var promotedPost model.Post
	if rows.Next() {
		if err := readRow(rows, &promotedPost); err != nil {
			return nil, err
		}
	} else {
		return nil, sql.ErrNoRows
	}

	return &promotedPost, nil
}

func readRow(rows *sql.Rows, post *model.Post) error {
	err := rows.Scan(
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
	return err
}
