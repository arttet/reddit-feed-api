package main

import (
	"flag"

	"github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/server"
	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/database"
	"github.com/arttet/reddit-feed-api/internal/telemetry"

	"github.com/pressly/goose/v3"

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
		panic(err)
	}

	cfg := config.GetConfigInstance()

	logger, err := cfg.Logger.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync() // nolint:errcheck

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

	tracing, err := telemetry.NewTracer(&cfg)
	if err != nil {
		logger.Error("tracing initialization", zap.Error(err))
		return
	}
	defer tracing.Close()

	if err := server.NewServer(logger, db).Start(&cfg); err != nil {
		logger.Error("server initialization", zap.Error(err))
	}
}
