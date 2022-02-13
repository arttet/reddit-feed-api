package feed_test

import (
	"context"
	"fmt"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/feed"
	"github.com/arttet/reddit-feed-api/internal/mock"
	"github.com/arttet/reddit-feed-api/internal/model"
	"github.com/arttet/reddit-feed-api/internal/test"

	"github.com/golang/mock/gomock"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Feed", func() {
	var (
		err    error
		errRow = fmt.Errorf("row error")

		ctx context.Context

		ctrl     *gomock.Controller
		mockRepo *mock.MockRepository

		service feed.Feed
	)

	BeforeEach(func() {
		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		Expect(ctrl).ShouldNot(BeNil())

		mockRepo = mock.NewMockRepository(ctrl)
		Expect(mockRepo).ShouldNot(BeNil())

		config := zap.NewDevelopmentConfig()
		config.Level.SetLevel(zapcore.PanicLevel)
		config.DisableCaller = true
		config.DisableStacktrace = true
		logger, _ := config.Build()
		defer logger.Sync()

		service = feed.NewFeed(mockRepo, logger)
	})

	// ////////////////////////////////////////////////////////////////////////

	Context("Create new posts", Label("CreatePosts"), func() {
		var (
			numberOfCreatedPosts       int64
			expectNumberOfCreatedPosts int64
		)

		BeforeEach(func() {
			posts := test.LoadTestData("data/posts.yaml").Posts
			expectNumberOfCreatedPosts = int64(len(posts))

			mockRepo.EXPECT().CreatePosts(gomock.Any(), gomock.Any()).Return(expectNumberOfCreatedPosts, nil).Times(1)
			numberOfCreatedPosts, err = service.CreatePosts(ctx, posts)
		})

		It("should not be an error", func() {
			Expect(err).Should(BeNil())
			Expect(numberOfCreatedPosts).Should(Equal(expectNumberOfCreatedPosts))
		})
	})

	// ////////////////////////////////////////////////////////////////////////

	Context("Generate a feed", Label("GenerateFeed"), func() {
		filenames := []string{
			"data/0_empty_feed.yaml",
			"data/1_tiny_feed.yaml",
			"data/2_not_sorted.yaml",
			"data/3_small_feed_without_promoted_post.yaml",
			"data/4_small_feed.yaml",
			"data/16_normal_feed.yaml",
			"data/17_large_feed.yaml",
			"data/27_huge_feed.yaml",
		}

		for _, filename := range filenames {
			testData := test.LoadTestData(filename)
			var feed model.Posts

			BeforeEach(func() {
				mockRepo.EXPECT().ListPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(testData.Posts, nil).Times(1)
				if len(testData.Posts) >= 3 {
					if testData.PromotedPosts != nil {
						mockRepo.EXPECT().GetPromotedPost(gomock.Any()).Return(testData.PromotedPosts[0], nil).Times(1)
					} else {
						mockRepo.EXPECT().GetPromotedPost(gomock.Any()).Return(nil, nil).Times(1)
					}
				}

				feed, err = service.GenerateFeed(ctx, 1)
			})

			It("should be an nil error", func() {
				Expect(err).Should(BeNil())
				Expect(testData.Feed).Should(Equal(feed))
			})
		}
	})

	Context("The Repository service returns an error", Label("GenerateFeed"), func() {
		var feed model.Posts

		BeforeEach(func() {
			mockRepo.EXPECT().ListPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errRow).Times(1)
			feed, err = service.GenerateFeed(ctx, 1)
		})

		It("should be an error", func() {
			Expect(errRow).Should(MatchError(err))
			Expect(feed).Should(BeNil())
		})
	})

	// ////////////////////////////////////////////////////////////////////////
})
