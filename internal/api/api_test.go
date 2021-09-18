package api_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/arttet/reddit-feed-api/internal/api"
	"github.com/arttet/reddit-feed-api/internal/model"
	"github.com/arttet/reddit-feed-api/internal/repo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api"
)

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
					Subreddit:      post.Subreddit,
					Score:          post.Score,
					Promoted:       post.Promoted,
					NotSafeForWork: post.NotSafeForWork,
				}

				if post.Link != "" {
					posts[i].PostType = &pb.Post_Link{
						Link: post.Link,
					}
				} else {
					posts[i].PostType = &pb.Post_Content{
						Content: post.Content,
					}
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
				Expect(status.Convert(err).Code()).Should(Equal(codes.InvalidArgument))
			})
		})

		Context("when fails to create because of an invalid argument", func() {
			for _, post := range testData.WrongPosts {
				post := post

				BeforeEach(func() {
					result := pb.Post{
						Title:          post.Title,
						Author:         post.Author,
						Subreddit:      post.Subreddit,
						Score:          post.Score,
						Promoted:       post.Promoted,
						NotSafeForWork: post.NotSafeForWork,
					}

					if post.Link != "" {
						result.PostType = &pb.Post_Link{
							Link: post.Link,
						}
					} else {
						result.PostType = &pb.Post_Content{
							Content: post.Content,
						}
					}

					request = &pb.CreatePostsV1Request{
						Posts: []*pb.Post{&result},
					}

					response, err = server.CreatePostsV1(ctx, request)
				})

				It("should return an empty response", func() {
					Expect(response).Should(BeNil())
					Expect(status.Convert(err).Code()).Should(Equal(codes.InvalidArgument))
				})
			}
		})

		Context("when fails to create because of a database connection error", func() {
			BeforeEach(func() {
				mock.ExpectExec("INSERT INTO post").
					WithArgs(values...).
					WillReturnError(sql.ErrConnDone)

				request = &pb.CreatePostsV1Request{
					Posts: posts,
				}

				response, err = server.CreatePostsV1(ctx, request)
			})

			It("should return an empty response", func() {
				Expect(response).Should(BeNil())
				Expect(status.Convert(err).Code()).Should(Equal(codes.Unavailable))
			})
		})
	})

	// ////////////////////////////////////////////////////////////////////////
})
