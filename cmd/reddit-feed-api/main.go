package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/server"
	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/service/repo"
	"github.com/arttet/reddit-feed-api/internal/broker"
	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/database"
	"github.com/arttet/reddit-feed-api/internal/telemetry"
	"github.com/pressly/goose"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

func main() {
	migration := flag.Bool("migration", true, "Defines the migration start option")
	configYML := flag.String("cfg", "config.yml", "Defines the configuration file option")
	flag.Parse()

	if err := config.ReadConfigYML(*configYML); err != nil {
		log.Fatal(err)
	}

	cfg := config.GetConfigInstance()

	logger, err := telemetry.NewLogger(&cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("starting service",
		zap.String("name", cfg.Project.Name),
		zap.String("yml", *configYML),
		zap.Bool("debug", cfg.Project.Debug),
		zap.String("environment", cfg.Project.Environment),
		zap.String("commit_hash", cfg.Project.CommitHash),
		zap.String("version", cfg.Project.Version),
	)

	db, err := database.NewConnection(cfg.Database.String(), cfg.Database.Driver)
	if err != nil {
		logger.Fatal("database initialization", zap.Error(err))
	}
	defer db.Close()

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.MigrationsDir); err != nil {
			logger.Error("migrations initialization", zap.Error(err))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	producer, err := broker.NewProducer(ctx, &cfg.Kafka, logger)
	if err != nil {
		logger.Error("failed to create a producer", zap.Error(err))
		return
	}
	logger.Info("the Kafka producer is running", zap.Strings("brokers", cfg.Kafka.Brokers))

	repository := repo.NewRepo(db)

	if err := server.NewServer(producer, repository, logger).Start(&cfg); err != nil {
		logger.Error("server initialization", zap.Error(err))
	}
}
