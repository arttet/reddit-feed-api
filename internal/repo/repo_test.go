package repo_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/arttet/reddit-feed-api/internal/model"
	"github.com/arttet/reddit-feed-api/internal/repo"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	var (
		err error

		mock sqlmock.Sqlmock

		db     *sql.DB
		sqlxDB *sqlx.DB

		ctx        context.Context
		repository repo.Repo
	)

	BeforeEach(func() {
		db, mock, err = sqlmock.New()
		Expect(err).Should(BeNil())
		sqlxDB = sqlx.NewDb(db, "sqlmock")

		ctx = context.Background()
		repository = repo.NewRepo(sqlxDB)
	})

	JustBeforeEach(func() {
	})

	AfterEach(func() {
		mock.ExpectClose()
		err = db.Close()
		Expect(err).Should(BeNil())
	})

	// ////////////////////////////////////////////////////////////////////////

	Describe("creating new posts", func() {
		var (
			exec                    *sqlmock.ExpectedExec
			numberOfTheCreatedPosts int64
		)

		BeforeEach(func() {
			values := make([]driver.Value, 0, reflect.TypeOf(model.Post{}).NumField()*len(testData.Posts))
			for _, post := range testData.Posts {
				values = append(values,
					post.Title,
					post.Author,
					post.Link,
					post.Subreddit,
					post.Content,
					post.Score,
					post.Promoted,
					post.NotSafeForWork,
				)
			}

			exec = mock.ExpectExec("INSERT INTO post").WithArgs(values...)
		})

		Context("when creates successfully", func() {
			var (
				lastInsertID int64
				rowsAffected int64
			)

			BeforeEach(func() {
				lastInsertID = int64(len(testData.Posts))
				rowsAffected = lastInsertID

				exec.WillReturnResult(sqlmock.NewResult(lastInsertID, rowsAffected))
				numberOfTheCreatedPosts, err = repository.CreatePosts(ctx, testData.Posts)
			})

			It("should return the number of the created posts correctly", func() {
				Expect(numberOfTheCreatedPosts).To(Equal(rowsAffected))
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to create", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrConnDone)
				numberOfTheCreatedPosts, err = repository.CreatePosts(ctx, testData.Posts)
			})

			It("should return the zero-value", func() {
				Expect(numberOfTheCreatedPosts).To(BeZero())
				Expect(err).Should(MatchError(sql.ErrConnDone))
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////

	Describe("lists posts", func() {
		const (
			limit  = 10
			offset = 0
		)

		var (
			exec   *sqlmock.ExpectedQuery
			result []model.Post
			rows   *sqlmock.Rows
		)

		BeforeEach(func() {
			query := fmt.Sprintf("SELECT (.+) FROM post ORDER BY %s DESC LIMIT %d OFFSET %d",
				"score",
				limit,
				offset,
			)
			exec = mock.ExpectQuery(query)

			rows = sqlmock.NewRows([]string{
				"id",
				"title",
				"author",
				"link",
				"subreddit",
				"content",
				"score",
				"promoted",
				"not_safe_for_work",
			})

			for i, post := range testData.Posts {
				if i == limit {
					break
				}

				rows.AddRow(
					post.ID,
					post.Title,
					post.Author,
					post.Link,
					post.Subreddit,
					post.Content,
					post.Score,
					post.Promoted,
					post.NotSafeForWork,
				)
			}
		})

		Context("when lists successfully", func() {
			BeforeEach(func() {
				exec.WillReturnRows(rows)
				result, err = repository.ListPosts(ctx, limit, offset)
			})

			It("should populate the slice correctly", func() {
				Expect(result).ShouldNot(BeEmpty())
				Expect(err).Should(BeNil())
			})
		})

		Context("when a database is empty", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrNoRows)
				result, err = repository.ListPosts(ctx, limit, offset)
			})

			It("should return an empty list of the posts", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrNoRows))
			})
		})

		Context("when fails to list", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrConnDone)
				result, err = repository.ListPosts(ctx, limit, offset)
			})

			It("should return an empty list of the posts", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrConnDone))
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////
})
