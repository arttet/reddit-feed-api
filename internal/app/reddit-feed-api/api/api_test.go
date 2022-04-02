package api_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/api"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/feed"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repository"
	"github.com/arttet/reddit-feed-api/internal/mock"
	"github.com/arttet/reddit-feed-api/internal/test"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/mock/gomock"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
	"github.com/arttet/reddit-feed-api/pkg/transform"
)

var _ = Describe("Reddit Feed API Server", func() {
	var (
		testData = test.LoadTestData("data/posts.yaml")

		err error

		ctx context.Context

		ctrl         *gomock.Controller
		mockProducer *mock.MockProducer
		mockSQL      sqlmock.Sqlmock
		db           *sql.DB
		sqlxDB       *sqlx.DB

		server pb.RedditFeedAPIServiceServer
	)

	BeforeEach(func() {
		Expect(testData).ShouldNot(BeNil())

		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		Expect(ctrl).ShouldNot(BeNil())

		mockProducer = mock.NewMockProducer(ctrl)
		Expect(mockProducer).ShouldNot(BeNil())
		mockProducer.EXPECT().CreatePosts(gomock.Any()).AnyTimes()

		db, mockSQL, err = sqlmock.New()
		Expect(err).Should(BeNil())
		sqlxDB = sqlx.NewDb(db, "sqlmock")
		Expect(sqlxDB).ShouldNot(BeNil())

		config := zap.NewDevelopmentConfig()
		config.Level.SetLevel(zapcore.PanicLevel)
		config.DisableCaller = true
		config.DisableStacktrace = true
		logger, _ := config.Build()
		defer logger.Sync()

		repo := repository.NewRepository(sqlxDB)
		feed := feed.NewFeed(repo, logger)

		server = api.NewRedditFeedAPIServiceServer(feed, mockProducer, logger)
	})

	// ////////////////////////////////////////////////////////////////////////

	Describe("creating new posts", func() {
		var (
			request  *pb.CreatePostsV1Request
			response *pb.CreatePostsV1Response
		)

		Describe("using a database", func() {
			var (
				query  = fmt.Sprintf("INSERT INTO %s", repository.TableName)
				values = test.PostToValuePtrValList(testData.Posts)
			)

			Context("when creates successfully", Label("CreatePostsV1"), func() {
				BeforeEach(func() {
					lastInsertID := int64(len(testData.Posts))
					rowsAffected := lastInsertID

					mockSQL.ExpectBegin()
					mockSQL.ExpectExec(query).WithArgs(values...).WillReturnResult(sqlmock.NewResult(lastInsertID, rowsAffected))
					mockSQL.ExpectCommit()
					mockSQL.ExpectClose()

					request = &pb.CreatePostsV1Request{Posts: transform.PostToPbPtrList(testData.Posts)}
					response, err = server.CreatePostsV1(ctx, request)
				})

				It("should return a number of the created posts correctly", func() {
					Expect(response.NumberOfCreatedPosts).Should(BeEquivalentTo(len(testData.Posts)))
					Expect(err).Should(BeNil())
				})
			})

			Context("when fails to create because of a database connection error", Label("CreatePostsV1"), func() {
				BeforeEach(func() {
					mockSQL.ExpectBegin()
					mockSQL.ExpectExec(query).WithArgs(values...).WillReturnError(sql.ErrConnDone)
					mockSQL.ExpectCommit()
					mockSQL.ExpectClose()

					request = &pb.CreatePostsV1Request{Posts: transform.PostToPbPtrList(testData.Posts)}
					response, err = server.CreatePostsV1(ctx, request)
				})

				It("should return an empty response", func() {
					Expect(response).Should(BeNil())
					Expect(status.Convert(err).Code()).Should(Equal(codes.Internal))
				})
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////

	Describe("generating a new feed", func() {
		const (
			feedLimit = 27
			offset    = 0
			pageId    = 1
		)

		var (
			request  *pb.GenerateFeedV1Request
			response *pb.GenerateFeedV1Response
		)

		Describe("using a database", func() {
			var (
				query = fmt.Sprintf("SELECT %s FROM %s ORDER BY score DESC LIMIT %d OFFSET %d",
					strings.Join(repository.SelectColumns, ", "),
					repository.TableName,
					feedLimit,
					offset,
				)
				promotedPostQuery = fmt.Sprintf("SELECT %s FROM %s WHERE promoted",
					strings.Join(repository.SelectColumns, ", "),
					repository.TableName,
				)
			)

			Context("when generates successfully", func() {
				filenames := []string{"feed1.yaml", "feed2.yaml"}
				for _, feed := range filenames {
					feed := feed
					var (
						rows            *sqlmock.Rows
						promotedPostRow *sqlmock.Rows
						testData        *test.TestData
					)

					BeforeEach(func() {
						filename := fmt.Sprintf("data/%s", feed)
						testData = test.LoadTestData(filename)
						Expect(testData).ShouldNot(BeNil())

						rows = sqlmock.NewRows(repository.SelectColumns)
						rows = test.PostToRowPtrList(testData.Posts, rows, feedLimit)

						promotedPostRow = sqlmock.NewRows(repository.SelectColumns)
						promotedPostRow = test.PostToRowPtrList(testData.PromotedPosts, promotedPostRow, feedLimit)
					})

					Context("with a promoted post", Label("GenerateFeedV1Request"), func() {
						BeforeEach(func() {
							mockSQL.ExpectBegin()
							mockSQL.ExpectQuery(query).WillReturnRows(rows)
							mockSQL.ExpectCommit()

							mockSQL.ExpectBegin()
							mockSQL.ExpectQuery(promotedPostQuery).WithArgs(true).WillReturnRows(promotedPostRow)
							mockSQL.ExpectCommit()

							mockSQL.ExpectClose()

							request = &pb.GenerateFeedV1Request{PageId: pageId}
							response, err = server.GenerateFeedV1(ctx, request)
						})

						It("should return a number of the created posts correctly", func() {
							Expect(response.Posts).ShouldNot(BeEmpty())
							Expect(err).Should(BeNil())

							posts := transform.PbToPostPtrList(response.Posts)
							Expect(testData.PromotedPosts[0]).Should(Equal(posts[1]))
						})
					})

					Context("without a promoted post", Label("GenerateFeedV1Request"), func() {
						BeforeEach(func() {
							mockSQL.ExpectBegin()
							mockSQL.ExpectQuery(query).WillReturnRows(rows)
							mockSQL.ExpectCommit()

							mockSQL.ExpectBegin()
							mockSQL.ExpectQuery(promotedPostQuery).WithArgs(true).WillReturnError(sql.ErrNoRows)
							mockSQL.ExpectRollback()

							mockSQL.ExpectClose()

							request = &pb.GenerateFeedV1Request{PageId: pageId}
							response, err = server.GenerateFeedV1(ctx, request)
						})

						It("should return a number of the created posts correctly", func() {
							Expect(response.Posts).ShouldNot(BeEmpty())
							Expect(err).Should(BeNil())
						})
					})
				}
			})

			Context("when generates unsuccessfully", Label("GenerateFeedV1Request"), func() {
				Context("when fails to generate because of a database connection error", Label("GenerateFeedV1Request"), func() {
					BeforeEach(func() {
						mockSQL.ExpectBegin()
						mockSQL.ExpectQuery(query).WillReturnError(sql.ErrConnDone)
						mockSQL.ExpectRollback()
						mockSQL.ExpectClose()

						request = &pb.GenerateFeedV1Request{PageId: pageId}
						response, err = server.GenerateFeedV1(ctx, request)
					})

					It("should return an empty response", func() {
						Expect(response).Should(BeNil())
						Expect(status.Convert(err).Code()).Should(Equal(codes.Internal))
					})
				})

				Context("when fails to generate because of a not found error", Label("GenerateFeedV1Request"), func() {
					BeforeEach(func() {
						mockSQL.ExpectBegin()
						mockSQL.ExpectQuery(query).WillReturnRows(sqlmock.NewRows(repository.SelectColumns))
						mockSQL.ExpectCommit()
						mockSQL.ExpectClose()

						request = &pb.GenerateFeedV1Request{PageId: pageId}
						response, err = server.GenerateFeedV1(ctx, request)
					})

					It("should return an empty response", func() {
						Expect(response).Should(BeNil())
						Expect(status.Convert(err).Code()).Should(Equal(codes.NotFound))
					})
				})
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////
})
