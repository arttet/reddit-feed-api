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

	"github.com/golang/mock/gomock"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
	"github.com/arttet/reddit-feed-api/pkg/transform"
)

var _ = Describe("Reddit Feed API Server", func() {
	var (
		err      error
		ctx      context.Context
		testData = test.LoadTestData("data/posts.yaml")

		ctrl         *gomock.Controller
		mockProducer *mock.MockProducer
		mockSQL      sqlmock.Sqlmock

		db     *sql.DB
		sqlxDB *sqlx.DB

		server pb.RedditFeedAPIServiceServer
	)

	BeforeEach(func() {
		ctx = context.Background()
		Expect(testData).ShouldNot(BeNil())

		ctrl = gomock.NewController(GinkgoT())
		Expect(ctrl).ShouldNot(BeNil())

		db, mockSQL, err = sqlmock.New()
		Expect(err).Should(BeNil())
		sqlxDB = sqlx.NewDb(db, "sqlmock")
		Expect(sqlxDB).ShouldNot(BeNil())

		mockProducer = mock.NewMockProducer(ctrl)
		Expect(mockProducer).ShouldNot(BeNil())
		mockProducer.EXPECT().CreatePosts(gomock.Any()).AnyTimes()

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
				exec *sqlmock.ExpectedExec
			)

			BeforeEach(func() {
				values := test.PostToValuePtrValList(testData.Posts)
				exec = mockSQL.ExpectExec(fmt.Sprintf("INSERT INTO %s", repository.TableName)).
					WithArgs(values...)
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

					request = &pb.CreatePostsV1Request{
						Posts: transform.PostToPbPtrList(testData.Posts),
					}

					response, err = server.CreatePostsV1(ctx, request)
				})

				It("should return a number of the created posts correctly", func() {
					Expect(response.NumberOfCreatedPosts).Should(BeEquivalentTo(len(testData.Posts)))
					Expect(err).Should(BeNil())
				})
			})

			Context("when fails to create because of a database connection error", func() {
				BeforeEach(func() {
					exec.WillReturnError(sql.ErrConnDone)

					request = &pb.CreatePostsV1Request{
						Posts: transform.PostToPbPtrList(testData.Posts),
					}

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
		)

		var (
			request  *pb.GenerateFeedV1Request
			response *pb.GenerateFeedV1Response
		)

		Describe("using a database", func() {
			Context("when generates successfully", func() {
				filenames := []string{"feed1.yaml", "feed2.yaml"}
				for _, feed := range filenames {
					feed := feed
					var (
						rows     *sqlmock.Rows
						testData *test.TestData
						query    string
					)

					BeforeEach(func() {
						filename := fmt.Sprintf("data/%s", feed)
						testData = test.LoadTestData(filename)
						Expect(testData).ShouldNot(BeNil())

						rows = sqlmock.NewRows(repository.SelectColumns)
						rows = test.PostToRowPtrList(testData.Posts, rows, feedLimit)

						query = fmt.Sprintf("SELECT %s FROM %s ORDER BY score DESC LIMIT %d OFFSET %d",
							strings.Join(repository.SelectColumns, ", "),
							repository.TableName,
							feedLimit,
							offset,
						)
					})

					Context("with a promoted post", func() {
						BeforeEach(func() {
							mockSQL.ExpectQuery(query).WillReturnRows(rows)

							promotedQuery := fmt.Sprintf("SELECT %s FROM %s WHERE promoted",
								strings.Join(repository.SelectColumns, ", "),
								repository.TableName,
							)

							mockSQL.ExpectQuery(promotedQuery).
								WithArgs(true).
								WillReturnRows(rows)

							request = &pb.GenerateFeedV1Request{
								PageId: 1,
							}

							response, err = server.GenerateFeedV1(ctx, request)
						})

						It("should return a number of the created posts correctly", func() {
							Expect(response.Posts).ShouldNot(BeEmpty())
							Expect(err).Should(BeNil())
						})
					})

					Context("without a promoted post", func() {
						BeforeEach(func() {
							mockSQL.ExpectQuery(query).WillReturnRows(rows)

							query = fmt.Sprintf("SELECT %s FROM %s WHERE promoted",
								strings.Join(repository.SelectColumns, ", "),
								repository.TableName,
							)
							mockSQL.ExpectQuery(query).WithArgs(true).
								WillReturnError(sql.ErrNoRows)

							request = &pb.GenerateFeedV1Request{
								PageId: 1,
							}

							response, err = server.GenerateFeedV1(ctx, request)
						})

						It("should return a number of the created posts correctly", func() {
							Expect(response.Posts).ShouldNot(BeEmpty())
							Expect(err).Should(BeNil())
						})
					})
				}
			})

			Context("when generates unsuccessfully", func() {
				var (
					exec *sqlmock.ExpectedQuery
				)

				BeforeEach(func() {
					query := fmt.Sprintf("SELECT %s FROM %s ORDER BY score DESC LIMIT %d OFFSET %d",
						strings.Join(repository.SelectColumns, ", "),
						repository.TableName,
						feedLimit,
						offset,
					)
					exec = mockSQL.ExpectQuery(query)
				})

				Context("when fails to generate because of a database connection error", func() {
					BeforeEach(func() {
						exec.WillReturnError(sql.ErrConnDone)

						request = &pb.GenerateFeedV1Request{
							PageId: 1,
						}

						response, err = server.GenerateFeedV1(ctx, request)
					})

					It("should return an empty response", func() {
						Expect(response).Should(BeNil())
						Expect(status.Convert(err).Code()).Should(Equal(codes.Internal))
					})
				})

				Context("when fails to generate because of a not found error", func() {
					BeforeEach(func() {
						exec.WillReturnRows(sqlmock.NewRows(repository.SelectColumns))

						request = &pb.GenerateFeedV1Request{
							PageId: 1,
						}

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
