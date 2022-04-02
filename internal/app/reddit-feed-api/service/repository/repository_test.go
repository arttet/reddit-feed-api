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
		testData = test.LoadTestData("data/posts.yaml")

		err    error
		errRow = fmt.Errorf("row error")

		ctx context.Context

		mockSQL sqlmock.Sqlmock
		db      *sql.DB
		sqlxDB  *sqlx.DB

		repo repository.Repository
	)

	BeforeEach(func() {
		Expect(testData).ShouldNot(BeNil())

		ctx = context.Background()

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
			query  = fmt.Sprintf("INSERT INTO %s", repository.TableName)
			values = test.PostToValuePtrValList(testData.Posts)

			numberOfTheCreatedPosts int64
		)

		Context("when creates successfully", Label("CreatePosts"), func() {
			BeforeEach(func() {
				lastInsertID := int64(len(testData.Posts))
				rowsAffected := lastInsertID

				mockSQL.ExpectBegin()
				mockSQL.ExpectExec(query).WithArgs(values...).WillReturnResult(sqlmock.NewResult(lastInsertID, rowsAffected))
				mockSQL.ExpectCommit()
				mockSQL.ExpectClose()

				numberOfTheCreatedPosts, err = repo.CreatePosts(ctx, testData.Posts)
			})

			It("should return the number of the created posts correctly", func() {
				Expect(numberOfTheCreatedPosts).Should(BeEquivalentTo(len(testData.Posts)))
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to create because of a connection done error", Label("CreatePosts"), func() {
			BeforeEach(func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectExec(query).WithArgs(values...).WillReturnError(sql.ErrConnDone)
				mockSQL.ExpectRollback()
				mockSQL.ExpectClose()

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
			query = fmt.Sprintf("SELECT %s FROM %s ORDER BY score DESC LIMIT %d OFFSET %d",
				strings.Join(repository.SelectColumns, ", "),
				repository.TableName,
				limit,
				offset,
			)

			rows   *sqlmock.Rows
			result model.Posts
		)

		BeforeEach(func() {
			rows = sqlmock.NewRows(repository.SelectColumns)
			rows = test.PostToRowPtrList(testData.Posts, rows, limit)
		})

		Context("when lists successfully", Label("ListPosts"), func() {
			BeforeEach(func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectQuery(query).WillReturnRows(rows)
				mockSQL.ExpectCommit()
				mockSQL.ExpectClose()

				result, err = repo.ListPosts(ctx, limit, offset)
			})

			It("should populate the slice correctly", func() {
				Expect(result).ShouldNot(BeEmpty())
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to list because of a connection done error", Label("ListPosts"), func() {
			BeforeEach(func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectQuery(query).WillReturnError(sql.ErrConnDone)
				mockSQL.ExpectRollback()
				mockSQL.ExpectClose()

				result, err = repo.ListPosts(ctx, limit, offset)
			})

			It("should return an empty list of the posts", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrConnDone))
			})
		})

		Context("when fails to list because of a no rows error", Label("ListPosts"), func() {
			BeforeEach(func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectQuery(query).WillReturnRows(sqlmock.NewRows(repository.SelectColumns))
				mockSQL.ExpectCommit()
				mockSQL.ExpectClose()

				result, err = repo.ListPosts(ctx, limit, offset)
			})

			It("should return an empty list of the posts", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrNoRows))
			})
		})

		Context("when fails to get because of a row error", Label("ListPosts"), func() {
			BeforeEach(func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectQuery(query).WillReturnRows(rows.RowError(0, errRow))
				mockSQL.ExpectRollback()
				mockSQL.ExpectClose()

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
			query = fmt.Sprintf("SELECT %s FROM %s WHERE promoted",
				strings.Join(repository.SelectColumns, ", "),
				repository.TableName,
			)

			rows   *sqlmock.Rows
			result *model.Post
		)

		BeforeEach(func() {
			rows = sqlmock.NewRows(repository.SelectColumns)
			rows = test.PostToRowPtrList(testData.PromotedPosts, rows, limit)
		})

		Context("when gets successfully", Label("GetPromotedPost"), func() {
			BeforeEach(func() {

				mockSQL.ExpectBegin()
				mockSQL.ExpectQuery(query).WithArgs(true).WillReturnRows(rows)
				mockSQL.ExpectCommit()
				mockSQL.ExpectClose()

				result, err = repo.GetPromotedPost(ctx)
			})

			It("should populate the slice correctly", func() {
				Expect(result).ShouldNot(BeNil())
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to get because of a connection done error", Label("GetPromotedPost"), func() {
			BeforeEach(func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectQuery(query).WithArgs(true).WillReturnError(sql.ErrConnDone)
				mockSQL.ExpectRollback()
				mockSQL.ExpectClose()

				result, err = repo.GetPromotedPost(ctx)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrConnDone))
			})
		})

		Context("when fails to get because of a no rows error", Label("GetPromotedPost"), func() {
			BeforeEach(func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectQuery(query).WillReturnRows(sqlmock.NewRows(repository.SelectColumns))
				mockSQL.ExpectCommit()
				mockSQL.ExpectClose()

				result, err = repo.GetPromotedPost(ctx)
			})

			It("should return an empty promoted post", func() {
				Expect(result).Should(BeNil())
				Expect(err).Should(MatchError(sql.ErrNoRows))
			})
		})

		Context("when fails to get because of a row error", Label("GetPromotedPost"), func() {
			BeforeEach(func() {
				mockSQL.ExpectBegin()
				mockSQL.ExpectQuery(query).WithArgs(true).WillReturnRows(rows.RowError(0, errRow))
				mockSQL.ExpectRollback()
				mockSQL.ExpectClose()

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
