package api_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"

	"github.com/arttet/reddit-feed-api/internal/reddit-feed-api/api"
	"github.com/arttet/reddit-feed-api/internal/reddit-feed-api/model"
	"github.com/arttet/reddit-feed-api/internal/reddit-feed-api/repo"
	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var errDatabaseConnection = errors.New("error establishing a database connection")

var _ = Describe("Reddit Feed API Server", func() {
	var (
		err error

		db     *sql.DB
		sqlxDB *sqlx.DB
		mock   sqlmock.Sqlmock

		ctx        context.Context
		repository repo.Repo
		server     pb.RedditFeedAPIServiceServer
	)

	BeforeEach(func() {
		db, mock, err = sqlmock.New()
		Expect(err).Should(BeNil())
		sqlxDB = sqlx.NewDb(db, "sqlmock")

		ctx = context.Background()
		repository = repo.NewRepo(sqlxDB)
		server = api.NewRedditFeedAPI(repository, 1024)
	})

	JustBeforeEach(func() {
	})

	AfterEach(func() {
		mock.ExpectClose()
		err := db.Close()
		Expect(err).Should(BeNil())
	})

	// ////////////////////////////////////////////////////////////////////////

	Describe("creating new posts", func() {
		var (
			values   []driver.Value
			posts    []*pb.Post
			request  *pb.CreatePostsV1Request
			response *pb.CreatePostsV1Response
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

			posts = make([]*pb.Post, len(testData.Posts))
			for i, post := range testData.Posts {
				posts[i] = &pb.Post{
					Title:          post.Title,
					Author:         post.Author,
					Link:           post.Link,
					Subreddit:      post.Subreddit,
					Content:        post.Content,
					Score:          post.Score,
					Promoted:       post.Promoted,
					NotSafeForWork: post.NotSafeForWork,
				}
			}
		})

		Context("when creates successfully", func() {
			var (
				lastInsertID int64
				rowsAffected int64
			)

			BeforeEach(func() {
				lastInsertID = int64(len(testData.Posts))
				rowsAffected = lastInsertID

				mock.ExpectExec("INSERT INTO post").
					WithArgs(values...).
					WillReturnResult(sqlmock.NewResult(lastInsertID, rowsAffected))

				request = &pb.CreatePostsV1Request{
					Posts: posts,
				}

				response, err = server.CreatePostsV1(ctx, request)
			})

			It("should return a number of the created posts correctly", func() {
				Expect(response.NumberOfCreatedPosts).Should(BeEquivalentTo(len(testData.Posts)))
			})

			It("should not be an error", func() {
				Expect(err).Should(BeNil())
			})
		})

		Context("when fails to create because of an empty argument", func() {
			BeforeEach(func() {
				request = &pb.CreatePostsV1Request{
					Posts: nil,
				}

				response, err = server.CreatePostsV1(ctx, request)
			})

			It("should return an empty response", func() {
				Expect(response).Should(BeNil())
			})

			It("should return an invalid argument error", func() {
				Expect(status.Convert(err).Code()).Should(Equal(codes.InvalidArgument))
			})
		})

		Context("when fails to create because of an invalid argument", func() {
			for _, post := range testData.WrongPosts {
				post := post

				BeforeEach(func() {
					request = &pb.CreatePostsV1Request{
						Posts: []*pb.Post{
							{
								Title:          post.Title,
								Author:         post.Author,
								Link:           post.Link,
								Subreddit:      post.Subreddit,
								Content:        post.Content,
								Score:          post.Score,
								Promoted:       post.Promoted,
								NotSafeForWork: post.NotSafeForWork,
							},
						},
					}

					response, err = server.CreatePostsV1(ctx, request)
				})

				It("should return an empty response", func() {
					Expect(response).Should(BeNil())
				})

				It("should return an invalid argument error", func() {
					Expect(status.Convert(err).Code()).Should(Equal(codes.InvalidArgument))
				})
			}
		})

		Context("when fails to create because of a database connection error", func() {
			BeforeEach(func() {
				mock.ExpectExec("INSERT INTO post").
					WithArgs(values...).
					WillReturnError(errDatabaseConnection)

				request = &pb.CreatePostsV1Request{
					Posts: posts,
				}

				response, err = server.CreatePostsV1(ctx, request)
			})

			It("should return an empty response", func() {
				Expect(response).Should(BeNil())
			})

			It("should return a resource exhausted error", func() {
				Expect(status.Convert(err).Code()).Should(Equal(codes.ResourceExhausted))
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////
})
