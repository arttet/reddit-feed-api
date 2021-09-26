package main

import (
	"flag"
	"fmt"

	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/database"
	"github.com/arttet/reddit-feed-api/internal/server"
	"github.com/arttet/reddit-feed-api/internal/tracer"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

func main() {
	migration := flag.String("migration", "", "Defines the migration start option")
	configYML := flag.String("cfg", "config.yml", "Defines the configuration file option")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	if err := config.ReadConfigYML(*configYML); err != nil {
		logger.Fatal("reading configuration", zap.Error(err))
	}

	cfg := config.GetConfigInstance()

	logger, err = cfg.Logger.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync() // nolint:errcheck

	logger.Info("starting service",
		zap.String("name", cfg.Project.Name),
		zap.Bool("debug", cfg.Project.Debug),
		zap.String("environment", cfg.Project.Environment),
		zap.String("commit_hash", cfg.Project.CommitHash),
		zap.String("version", cfg.Project.Version),
	)

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := database.NewConnection(dsn, cfg.Database.Driver)
	if err != nil {
		logger.Fatal("database initialization", zap.Error(err))
	}
	defer db.Close()

	if *migration != "" {
		if err := database.Migrate(db.DB, *migration); err != nil {
			logger.Error("migrations initialization", zap.Error(err))
			return
		}
	}

	tracing, err := tracer.NewTracer(&cfg)
	if err != nil {
		logger.Error("tracing initialization", zap.Error(err))
		return
	}
	defer tracing.Close()

	if err := server.NewServer(logger, db).Start(&cfg); err != nil {
		logger.Error("server initialization", zap.Error(err))
	}
}
