package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/database"
	"github.com/arttet/reddit-feed-api/internal/server"
	"github.com/arttet/reddit-feed-api/internal/tracer"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

func main() {
	migration := flag.String("migration", "", "Defines the migration start option")
	configYML := flag.String("cfg", "config.yml", "Defines the configuration file option")
	flag.Parse()

	if err := config.ReadConfigYML(*configYML); err != nil {
		log.Fatal().
			Err(err).
			Msg("Reading configuration")
	}

	cfg := config.GetConfigInstance()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

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
		log.Fatal().Err(err).Msg("db initialization")
	}
	defer db.Close()

	if *migration != "" {
		if err := database.Migrate(db.DB, *migration); err != nil {
			log.Error().Err(err).Msg("migrations initialization")
			return
		}
	}

	tracing, err := tracer.NewTracer(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to init tracing")
		return
	}
	defer tracing.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if err := server.NewServer(db).Start(&cfg); err != nil {
			log.Error().Err(err).Msg("Failed creating gRPC server")
			return
		}
		wg.Done()
	}()

	wg.Wait()
}
