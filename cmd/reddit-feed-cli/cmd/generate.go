package cmd

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
	"github.com/arttet/reddit-feed-api/pkg/transform"
)

// generateCmd represents the generate command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Client: generate a feed of posts",

	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewRedditFeedAPIServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		response, err := client.GenerateFeedV1(ctx, &pb.GenerateFeedV1Request{
			PageId: pageID,
		})
		if err != nil {
			log.Fatalf("could not generate: %v", err)
		}

		posts := transform.PbToPostPtrList(response.Posts)
		for _, post := range posts {
			log.Printf("Post: %v", post)
		}
	},
}

func init() {
	generateCmd.Flags().StringVarP(&addr, "url", "u", "localhost:8082", "the address to connect to the gRPC server")
	generateCmd.Flags().Uint64VarP(&pageID, "page", "p", 1, "the page ID of the feed")
}
