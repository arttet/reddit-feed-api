package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/api"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/broker"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/server"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/feed"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repository"
	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/database"
	"github.com/arttet/reddit-feed-api/internal/telemetry"

	"github.com/pressly/goose"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

var (
	cfgFile   string
	migration bool
)

var rootCmd = &cobra.Command{
	Use:   "reddit-feed-api",
	Short: "The Reddit Feed API is the interface used to create new posts and generate feeds of posts.",

	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConfig()

		db, err := database.NewConnection(cfg.Database.String(), cfg.Database.Driver)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if migration {
			if err = goose.Up(db.DB, cfg.Database.MigrationsDir); err != nil {
				log.Fatal(err)
			}
		}

		logger, err := telemetry.NewLogger(&cfg.Logger)
		if err != nil {
			log.Fatal(err)
		}
		defer logger.Sync() // nolint:errcheck

		logger.Info("starting service",
			zap.String("name", cfg.Project.Name),
			zap.String("yml", cfgFile),
			zap.Bool("debug", cfg.Project.Debug),
			zap.String("environment", cfg.Project.Environment),
			zap.String("commit_hash", cfg.Project.CommitHash),
			zap.String("version", cfg.Project.Version),
		)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		producer, err := broker.NewProducer(ctx, &cfg.Producer, logger)
		if err != nil {
			log.Fatal(err)
		}

		logger.Info("the Kafka producer is running", zap.Strings("brokers", cfg.Producer.Brokers))

		repo := repository.NewRepository(db)

		srv := api.NewRedditFeedAPIServiceServer(
			feed.NewFeed(repo, logger),
			producer,
			logger,
		)

		if err := server.NewServer(srv, logger).Serve(cfg); err != nil {
			logger.Error("server initialization", zap.Error(err))
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "config.yml", "the path to the configuration file")
	rootCmd.Flags().BoolVarP(&migration, "migration", "m", true, "the migration start flag")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getConfig() *config.Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if cfgFile != "config.yml" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("./configs/reddit-feed-api/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read the config file: %v", err)
	}

	cfg := &config.Config{}
	if err := config.ReadConfigYML(viper.ConfigFileUsed(), cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
