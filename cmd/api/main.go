package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guidiguidi/portfolio-tracker/internal/assets"
	"github.com/guidiguidi/portfolio-tracker/internal/config"
	httpapi "github.com/guidiguidi/portfolio-tracker/internal/http"
	"github.com/guidiguidi/portfolio-tracker/internal/logger"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		// Can't use logger here as it's not initialized yet
		log.Fatalf("failed to load config: %v", err)
	}

	log := logger.New(cfg.Logger.Level, cfg.Logger.JSON)
	log.Info("starting application", slog.String("version", cfg.App.Version))

	runMigrations(cfg.Database.DSN, log)

	db, err := config.NewDB(cfg)
	if err != nil {
		log.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Info("database connection established")

	var assetsRepo assets.Repository
	if cfg.Database.Driver == "postgres" {
		assetsRepo = assets.NewPostgresRepo(db, log)
		log.Info("using postgres assets repository")
	} else {
		assetsRepo = assets.NewMemoryRepo(log)
		log.Info("using in-memory assets repository")
	}

	assetsHandler := assets.NewHandler(assetsRepo, log)

	r := httpapi.NewRouter(log, assetsHandler)

	addr := ":" + cfg.App.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Info("starting server", slog.String("address", addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	shutdownTimeout, err := time.ParseDuration(cfg.App.ShutdownTimeout)
	if err != nil {
		log.Error("invalid shutdown timeout duration in config, using default 5s", "error", err)
		shutdownTimeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("server exiting")
}

func runMigrations(dsn string, log *slog.Logger) {
	if dsn == "" {
		log.Warn("database DSN is not set, skipping migrations")
		return
	}

	log.Info("running database migrations")
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Error("could not create migrate instance", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	log.Info("migrations completed successfully")
}