package cmd

import (
	"context"
	"log"
	"time"

	"github.com/arttet/reddit-feed-api/internal/model"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api/v1"
	"github.com/arttet/reddit-feed-api/pkg/transform"
)

type createConfig struct {
	Posts model.Posts `yaml:"posts"`
}

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Client: create new posts",

	Run: func(cmd *cobra.Command, args []string) {
		cfg := getCreateConfig()

		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewRedditFeedAPIServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		response, err := client.CreatePostsV1(ctx, &pb.CreatePostsV1Request{
			Posts: transform.PostToPbPtrList(cfg.Posts),
		})
		if err != nil {
			log.Fatalf("could not create: %v", err)
		}

		log.Printf("Number of created posts: %d", response.NumberOfCreatedPosts)
	},
}

func init() {
	createCmd.Flags().StringVarP(&cfgFile, "config", "c", "create.yml", "the path to the configuration file")
	createCmd.Flags().StringVarP(&addr, "url", "", "localhost:8082", "the address to connect to the gRPC server")
}

func getCreateConfig() *createConfig {
	viper.SetConfigName("create")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(cfgFile)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("./configs/reddit-feed-api/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read the config file: %v", err)
	}

	cfg := &createConfig{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
