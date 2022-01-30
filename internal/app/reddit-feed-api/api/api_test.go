package api_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/api"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repo"
	"github.com/arttet/reddit-feed-api/internal/mock"
	"github.com/arttet/reddit-feed-api/internal/model"
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

		repository repo.Repo
		server     pb.RedditFeedAPIServiceServer
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

		repository = repo.NewRepo(sqlxDB)
		server = api.NewRedditFeedAPI(repository, mockProducer, logger)
	})

	AfterEach(func() {
		mockSQL.ExpectClose()
		err := db.Close()
		Expect(err).Should(BeNil())

		ctrl.Finish()
	})

	// ////////////////////////////////////////////////////////////////////////

	Describe("creating new posts", func() {
		var (
			request  *pb.CreatePostsV1Request
			response *pb.CreatePostsV1Response
		)

		Describe("using a database", func() {
			var (
				exec   *sqlmock.ExpectedExec
				values []driver.Value
				posts  []*pb.Post
			)

			BeforeEach(func() {
				values = make([]driver.Value, 0, reflect.TypeOf(model.Post{}).NumField()*len(testData.Posts))
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

				posts = make([]*pb.Post, 0, len(testData.Posts))
				for _, post := range testData.Posts {
					result := &pb.Post{
						Title:          post.Title,
						Author:         post.Author,
						Subreddit:      post.Subreddit,
						Score:          post.Score,
						Promoted:       post.Promoted,
						NotSafeForWork: post.NotSafeForWork,
					}
					if post.Link != "" {
						result.PostType = &pb.Post_Link{Link: post.Link}
					} else {
						result.PostType = &pb.Post_Content{Content: post.Content}
					}
					posts = append(posts, result)
				}

				exec = mockSQL.ExpectExec(fmt.Sprintf("INSERT INTO %s", repo.TableName)).
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
						Posts: posts,
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
						Posts: posts,
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
			feedLimit     = 27
			promotedLimit = 1
			offset        = 0
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

						rows = sqlmock.NewRows(repo.SelectColumns)
						for i, post := range testData.Posts {
							if i == feedLimit {
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
						query = fmt.Sprintf("SELECT %s FROM %s ORDER BY score DESC LIMIT %d OFFSET %d",
							strings.Join(repo.SelectColumns, ", "),
							repo.TableName,
							feedLimit,
							offset,
						)
					})

					Context("with a promoted post", func() {
						BeforeEach(func() {
							mockSQL.ExpectQuery(query).WillReturnRows(rows)

							promotedQuery := fmt.Sprintf("SELECT %s FROM %s WHERE promoted",
								strings.Join(repo.SelectColumns, ", "),
								repo.TableName,
							)

							promotedRows := sqlmock.NewRows(repo.SelectColumns)
							for i, post := range testData.PromotedPosts {
								if i == promotedLimit {
									break
								}

								promotedRows.AddRow(
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
								strings.Join(repo.SelectColumns, ", "),
								repo.TableName,
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
						strings.Join(repo.SelectColumns, ", "),
						repo.TableName,
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
						exec.WillReturnRows(sqlmock.NewRows(repo.SelectColumns))

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
