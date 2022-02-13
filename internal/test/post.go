package test

import (
	"database/sql/driver"
	"reflect"

	"github.com/arttet/reddit-feed-api/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func PostToRowPtr(src *model.Post, rows *sqlmock.Rows) *sqlmock.Rows {
	if src == nil {
		return rows
	}

	rows = rows.AddRow(
		src.ID,
		src.Title,
		src.Author,
		src.Link,
		src.Subreddit,
		src.Content,
		src.Score,
		src.Promoted,
		src.NotSafeForWork,
	)

	return rows
}

func PostToRowPtrList(src model.Posts, rows *sqlmock.Rows, limit int) *sqlmock.Rows {
	if src == nil {
		return rows
	}

	for i, s := range src {
		if i == limit {
			break
		}
		rows = PostToRowPtr(s, rows)
	}

	return rows
}

func PostToValuePtrVal(src *model.Post) []driver.Value {
	if src == nil {
		return nil
	}

	return []driver.Value{
		src.Title,
		src.Author,
		src.Link,
		src.Subreddit,
		src.Content,
		src.Score,
		src.Promoted,
		src.NotSafeForWork,
	}
}

func PostToValuePtrValList(src model.Posts) []driver.Value {
	resp := make([]driver.Value, 0, reflect.TypeOf(model.Post{}).NumField()*len(src))

	for _, s := range src {
		resp = append(resp, PostToValuePtrVal(s)...)
	}

	return resp
}
