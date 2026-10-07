// Command server runs the Pantry HTTP API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/darko-komericki/pantry/backend/internal/api"
	"github.com/darko-komericki/pantry/backend/internal/config"
)

const shutdownTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "server:", err)
		os.Exit(1)
	}
}

// run holds the real startup logic. Keeping it out of main means every
// failure is an ordinary returned error and deferred cleanup always runs
// (os.Exit skips defers, so only main calls it).
func run() error {
	// ctx is cancelled on Ctrl+C or SIGTERM; that is the shutdown signal.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := newLogger(cfg.LogFormat)

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("create db pool: %w", err)
	}
	defer pool.Close()

	// pgxpool connects lazily; ping so a bad DATABASE_URL fails at startup.
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	// pool satisfies api.Pinger, so it is passed in directly.
	strict := api.NewStrictHandler(api.NewServer(pool, logger), nil)
	handler := api.HandlerFromMuxWithBaseURL(strict, http.NewServeMux(), "/api")

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// ListenAndServe blocks, so it runs in its own goroutine. If it fails
	// (e.g. port in use) the error comes back through serveErr.
	serveErr := make(chan error, 1)
	go func() {
		logger.Info("http server listening", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()

	// Wait for whichever happens first: a server failure or a shutdown signal.
	select {
	case err := <-serveErr:
		return fmt.Errorf("http server: %w", err)
	case <-ctx.Done():
	}

	logger.Info("shutting down")

	// ctx is already cancelled, so shutdown gets a fresh deadline of its own.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}
	logger.Info("stopped")
	return nil
}

func newLogger(format string) *slog.Logger {
	if format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
