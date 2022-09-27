package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/broker"
	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/consume"
	"github.com/arttet/reddit-feed-api/internal/telemetry"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

type consumeConfig struct {
	Consumer config.Consume `yaml:"consumer"`
	Logger   zap.Config     `yaml:"logger"`
}

// consumerCmd represents the Kafka consumer command.
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Kafka: read the events",

	Run: func(cmd *cobra.Command, args []string) {
		cfg := getConsumeConfig()

		logger, err := telemetry.NewLogger(&cfg.Logger)
		if err != nil {
			log.Fatal(err)
		}
		defer logger.Sync() //nolint:errcheck

		ctx, cancel := context.WithCancel(context.Background())

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err = consume.StartConsuming(ctx, cfg.Consumer, broker.ReceiveMessage(logger), logger); err != nil {
				log.Fatal(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		select {
		case v := <-quit:
			logger.Info("", zap.String("signal.Notify", v.String()))
		case done := <-ctx.Done():
			logger.Info("", zap.String("ctx.Done", fmt.Sprintf("%v", done)))
		}

		cancel()
		wg.Wait()
	},
}

func init() {
	consumerCmd.Flags().StringVarP(&cfgFile, "config", "c", "consumer.yml", "the path to the configuration file")
}

func getConsumeConfig() *consumeConfig {
	viper.SetConfigName("consumer")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(cfgFile)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("./configs/reddit-feed-api/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read the config file: %v", err)
	}

	cfg := &consumeConfig{}
	if err := config.ReadConfigYML(viper.ConfigFileUsed(), cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
