package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/arttet/reddit-feed-api/internal/api"
	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/database"
	"github.com/arttet/reddit-feed-api/internal/repo"
	"github.com/arttet/reddit-feed-api/internal/tracer"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	pb "github.com/arttet/reddit-feed-api/pkg/reddit-feed-api"
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
	r := repo.NewRepo(db)

	wg.Add(1)
	go func() {
		if err := runGRPC(r, cfg.GRPC.Host, cfg.GRPC.Port, cfg.Project.Debug); err != nil {
			log.Fatal().Err(err).Msg("Failed creating gRPC server")
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := runHTTP(&cfg); err != nil {
			log.Fatal().Err(err).Msg("Failed creating REST server")
		}
		wg.Done()
	}()

	wg.Wait()
}

func runGRPC(r repo.Repo, host string, port int, debug bool) error {
	listenOn := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", listenOn, err)
	}

	server := grpc.NewServer()
	if debug {
		reflection.Register(server)
	}
	pb.RegisterRedditFeedAPIServiceServer(server, api.NewRedditFeedAPI(r))

	log.Info().Msgf("Listening GRPC server on %s", listenOn)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}
	return nil
}

func runHTTP(cfg *config.Config) error {
	grpcEndpoint := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	httpEndpoint := fmt.Sprintf("%s:%d", cfg.REST.Host, cfg.REST.Port)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterRedditFeedAPIServiceHandlerFromEndpoint(ctx, gwmux, grpcEndpoint, opts)
	if err != nil {
		return fmt.Errorf("failed to register gRPC server: %w", err)
	}

	mux := http.NewServeMux()
	if cfg.Project.Debug {
		mux.HandleFunc("/swagger/", serveSwagger)
	}
	mux.Handle("/", gwmux)

	log.Info().Msgf("Listening HTTP server on %s", httpEndpoint)
	return http.ListenAndServe(httpEndpoint, mux)
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		log.Error().Msgf("Swagger JSON not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	log.Info().Msgf("Serving %s", r.URL.Path)

	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	cfg := config.GetConfigInstance()
	p = path.Join(cfg.Project.SwaggerDir, p)

	http.ServeFile(w, r, p)
}
