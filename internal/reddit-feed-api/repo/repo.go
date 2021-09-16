package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/arttet/reddit-feed-api/internal/reddit-feed-api/model"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const tableName = "post"

var ErrPostNotFound = errors.New("post not found")

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

type repo struct {
	db *sqlx.DB
}

func (r *repo) CreatePosts(ctx context.Context, posts []model.Post) (int64, error) {
	query := squirrel.Insert(tableName).
		Columns(
			"title",
			"author",
			"link",
			"subreddit",
			"content",
			"score",
			"promoted",
			"not_safe_for_work",
		).
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
	query := squirrel.Select(
		"id",
		"title",
		"author",
		"link",
		"subreddit",
		"content",
		"score",
		"promoted",
		"not_safe_for_work",
	).
		From(tableName).
		OrderBy("score DESC").
		Limit(limit).
		Offset(offset).
		RunWith(r.db).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	posts := make([]model.Post, 0, 16)
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Author,
			&post.Link,
			&post.Subreddit,
			&post.Content,
			&post.Score,
			&post.Promoted,
			&post.NotSafeForWork,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *repo) PromotedPost(ctx context.Context) (*model.Post, error) {
	query := squirrel.Select(
		"id",
		"title",
		"author",
		"link",
		"subreddit",
		"content",
		"score",
		"promoted",
		"not_safe_for_work",
	).
		From(tableName).
		Where(squirrel.Eq{"promoted": true}).
		RunWith(r.db).
		OrderBy("score DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var promotedPost model.Post
	if rows.Next() {
		if err := rows.Scan(
			&promotedPost.ID,
			&promotedPost.Title,
			&promotedPost.Author,
			&promotedPost.Link,
			&promotedPost.Subreddit,
			&promotedPost.Content,
			&promotedPost.Score,
			&promotedPost.Promoted,
			&promotedPost.NotSafeForWork,
		); err != nil {
			return nil, err
		}
	} else {
		return nil, ErrPostNotFound
	}

	return &promotedPost, nil
}
