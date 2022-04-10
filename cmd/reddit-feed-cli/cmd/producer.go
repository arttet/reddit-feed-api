package cmd

import (
	"context"
	"log"
	"time"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/broker"
	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/model"
	"github.com/arttet/reddit-feed-api/internal/telemetry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

type produceConfig struct {
	Producer config.Producer `yaml:"producer"`
	Logger   zap.Config      `yaml:"logger"`
	Posts    model.Posts     `yaml:"posts"`
}

// producerCmd represents the Kafka producer command.
var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "Kafka: write some events into the topic",

	Run: func(cmd *cobra.Command, args []string) {
		cfg := getProduceConfig()

		logger, err := telemetry.NewLogger(&cfg.Logger)
		if err != nil {
			log.Fatal(err)
		}
		defer logger.Sync() // nolint:errcheck

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		producer, err := broker.NewProducer(ctx, &cfg.Producer, logger)
		if err != nil {
			log.Fatal(err)
		}
		logger.Info("the Kafka producer is running", zap.Strings("brokers", cfg.Producer.Brokers))

		producer.CreatePosts(cfg.Posts)
		<-ctx.Done()
	},
}

func init() {
	producerCmd.Flags().StringVarP(&cfgFile, "config", "c", "producer.yml", "the path to the configuration file")
}

func getProduceConfig() *produceConfig {
	viper.SetConfigName("producer")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(cfgFile)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("./configs/reddit-feed-api/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read the config file: %v", err)
	}

	cfg := &produceConfig{}
	if err := config.ReadConfigYML(viper.ConfigFileUsed(), cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
