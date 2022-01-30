// source file: reddit/reddit_feed_api/v1/reddit_feed_api.proto
// source package: reddit.reddit_feed_api.v1

package transform

import (
	"github.com/arttet/reddit-feed-api/internal/model"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
)

func PbToPostPtr(src *pb.Post) *model.Post {
	if src == nil {
		return nil
	}

	d := &model.Post{
		Title:          src.Title,
		Author:         src.Author,
		Subreddit:      src.Subreddit,
		Link:           src.GetLink(),
		Content:        src.GetContent(),
		Score:          src.Score,
		Promoted:       src.Promoted,
		NotSafeForWork: src.NotSafeForWork,
	}

	return d
}

func PbToPostPtrList(src []*pb.Post) []*model.Post {
	resp := make([]*model.Post, len(src))

	for i, s := range src {
		resp[i] = PbToPostPtr(s)
	}

	return resp
}

func PostToPbPtr(src *model.Post) *pb.Post {
	if src == nil {
		return nil
	}

	d := &pb.Post{
		Title:          src.Title,
		Author:         src.Author,
		Subreddit:      src.Subreddit,
		Score:          src.Score,
		Promoted:       src.Promoted,
		NotSafeForWork: src.NotSafeForWork,
	}

	if src.Link != "" {
		d.PostType = &pb.Post_Link{Link: src.Link}
	} else {
		d.PostType = &pb.Post_Content{Content: src.Content}
	}

	return d
}

func PostToPbPtrList(src []*model.Post) []*pb.Post {
	resp := make([]*pb.Post, len(src))

	for i, s := range src {
		resp[i] = PostToPbPtr(s)
	}

	return resp
}
