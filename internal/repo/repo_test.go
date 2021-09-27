package repo_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/arttet/reddit-feed-api/internal/model"
	"github.com/arttet/reddit-feed-api/internal/repo"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	var (
		err    error
		errRow = fmt.Errorf("row error")

		ctx context.Context

		mockSQL sqlmock.Sqlmock

		db     *sql.DB
		sqlxDB *sqlx.DB

		repository repo.Repo
	)

	BeforeEach(func() {
		ctx = context.Background()

		db, mockSQL, err = sqlmock.New()
		Expect(err).Should(BeNil())
		sqlxDB = sqlx.NewDb(db, "sqlmock")
		Expect(sqlxDB).ShouldNot(BeNil())

		repository = repo.NewRepo(sqlxDB)
	})

	JustBeforeEach(func() {
	})

	AfterEach(func() {
		mockSQL.ExpectClose()
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

			exec = mockSQL.ExpectExec(fmt.Sprintf("INSERT INTO %s", repo.TableName)).WithArgs(values...)
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

		Context("when fails to create because of a connection done error", func() {
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
			limit  = 27
			offset = 0
		)

		var (
			exec   *sqlmock.ExpectedQuery
			rows   *sqlmock.Rows
			result []model.Post
		)

		BeforeEach(func() {
			query := fmt.Sprintf("SELECT %s FROM %s ORDER BY score DESC LIMIT %d OFFSET %d",
				strings.Join(repo.SelectColumns, ", "),
				repo.TableName,
				limit,
				offset,
			)
			exec = mockSQL.ExpectQuery(query)

			rows = sqlmock.NewRows(repo.SelectColumns)
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

		Context("when fails to list because of a connection done error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrConnDone)
				result, err = repository.ListPosts(ctx, limit, offset)
			})

			It("should return an empty list of the posts", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrConnDone))
			})
		})

		Context("when fails to list because of a no rows error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrNoRows)
				result, err = repository.ListPosts(ctx, limit, offset)
			})

			It("should return an empty list of the posts", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrNoRows))
			})
		})

		Context("when fails to get because of a row error", func() {
			BeforeEach(func() {
				exec.WillReturnRows(rows.RowError(0, errRow))
				result, err = repository.ListPosts(ctx, limit, offset)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(errRow))
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////

	Describe("getting a promoted post", func() {
		const (
			limit = 1
		)

		var (
			exec   *sqlmock.ExpectedQuery
			rows   *sqlmock.Rows
			result *model.Post
		)

		BeforeEach(func() {
			query := fmt.Sprintf("SELECT %s FROM %s WHERE promoted",
				strings.Join(repo.SelectColumns, ", "),
				repo.TableName,
			)
			exec = mockSQL.ExpectQuery(query).WithArgs(true)

			rows = sqlmock.NewRows(repo.SelectColumns)
			for i, post := range testData.PromotedPosts {
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

		Context("when gets successfully", func() {
			BeforeEach(func() {
				exec.WillReturnRows(rows)
				result, err = repository.GetPromotedPost(ctx)
			})

			It("should populate the slice correctly", func() {
				Expect(result).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to get because of a connection done error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrConnDone)
				result, err = repository.GetPromotedPost(ctx)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrConnDone))
			})
		})

		Context("when fails to get because of a no rows error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrNoRows)
				result, err = repository.GetPromotedPost(ctx)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrNoRows))
			})
		})

		Context("when fails to get because of a row error", func() {
			BeforeEach(func() {
				exec.WillReturnRows(rows.RowError(0, errRow))
				result, err = repository.GetPromotedPost(ctx)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(errRow))
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////
})
