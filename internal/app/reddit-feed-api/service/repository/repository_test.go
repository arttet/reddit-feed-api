package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repository"
	"github.com/arttet/reddit-feed-api/internal/model"
	"github.com/arttet/reddit-feed-api/internal/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	var (
		err      error
		errRow   = fmt.Errorf("row error")
		testData = test.LoadTestData("data/posts.yaml")

		ctx context.Context

		mockSQL sqlmock.Sqlmock

		db     *sql.DB
		sqlxDB *sqlx.DB

		repo repository.Repository
	)

	BeforeEach(func() {
		ctx = context.Background()

		Expect(testData).ShouldNot(BeNil())

		db, mockSQL, err = sqlmock.New()
		Expect(err).Should(BeNil())
		sqlxDB = sqlx.NewDb(db, "sqlmock")
		Expect(sqlxDB).ShouldNot(BeNil())

		repo = repository.NewRepository(sqlxDB)
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
			values := test.PostToValuePtrValList(testData.Posts)
			exec = mockSQL.ExpectExec(fmt.Sprintf("INSERT INTO %s", repository.TableName)).WithArgs(values...)
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
				numberOfTheCreatedPosts, err = repo.CreatePosts(ctx, testData.Posts)
			})

			It("should return the number of the created posts correctly", func() {
				Expect(numberOfTheCreatedPosts).To(Equal(rowsAffected))
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to create because of a connection done error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrConnDone)
				numberOfTheCreatedPosts, err = repo.CreatePosts(ctx, testData.Posts)
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
			result model.Posts
		)

		BeforeEach(func() {
			query := fmt.Sprintf("SELECT %s FROM %s ORDER BY score DESC LIMIT %d OFFSET %d",
				strings.Join(repository.SelectColumns, ", "),
				repository.TableName,
				limit,
				offset,
			)
			exec = mockSQL.ExpectQuery(query)

			rows = sqlmock.NewRows(repository.SelectColumns)
			rows = test.PostToRowPtrList(testData.Posts, rows, limit)
		})

		Context("when lists successfully", func() {
			BeforeEach(func() {
				exec.WillReturnRows(rows)
				result, err = repo.ListPosts(ctx, limit, offset)
			})

			It("should populate the slice correctly", func() {
				Expect(result).ShouldNot(BeEmpty())
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to list because of a connection done error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrConnDone)
				result, err = repo.ListPosts(ctx, limit, offset)
			})

			It("should return an empty list of the posts", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrConnDone))
			})
		})

		Context("when fails to list because of a no rows error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrNoRows)
				result, err = repo.ListPosts(ctx, limit, offset)
			})

			It("should return an empty list of the posts", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrNoRows))
			})
		})

		Context("when fails to get because of a row error", func() {
			BeforeEach(func() {
				exec.WillReturnRows(rows.RowError(0, errRow))
				result, err = repo.ListPosts(ctx, limit, offset)
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
				strings.Join(repository.SelectColumns, ", "),
				repository.TableName,
			)
			exec = mockSQL.ExpectQuery(query).WithArgs(true)

			rows = sqlmock.NewRows(repository.SelectColumns)
			rows = test.PostToRowPtrList(testData.PromotedPosts, rows, limit)
		})

		Context("when gets successfully", func() {
			BeforeEach(func() {
				exec.WillReturnRows(rows)
				result, err = repo.GetPromotedPost(ctx)
			})

			It("should populate the slice correctly", func() {
				Expect(result).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to get because of a connection done error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrConnDone)
				result, err = repo.GetPromotedPost(ctx)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrConnDone))
			})
		})

		Context("when fails to get because of a no rows error", func() {
			BeforeEach(func() {
				exec.WillReturnError(sql.ErrNoRows)
				result, err = repo.GetPromotedPost(ctx)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrNoRows))
			})
		})

		Context("when fails to get because of a row error", func() {
			BeforeEach(func() {
				exec.WillReturnRows(rows.RowError(0, errRow))
				result, err = repo.GetPromotedPost(ctx)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(errRow))
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////
})
