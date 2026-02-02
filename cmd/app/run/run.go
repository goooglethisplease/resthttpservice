package run

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

	_ "github.com/jackc/pgx/v5/stdlib"

	"restservice/internal/config"
	httpa "restservice/internal/http"
	"restservice/internal/repo"
	"restservice/internal/usecase/subscription"
)

const (
	defaultHTTPAddr    = ":8080"
	serverReadTimeout  = 5 * time.Second
	serverWriteTimeout = 10 * time.Second
	serverIdleTimeout  = 60 * time.Second
	shutdownTimeout    = 10 * time.Second
)

func Run() error {
	cfg := config.MustLoad()
	if cfg.DSN == "" {
		return errors.New("DSN env is empty")
	}

	db, err := repo.Open(cfg.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	subRepo := repo.NewSubscriptionRepo(db)
	subService := subscription.NewService(subRepo, logger)
	handler := httpa.NewHandler(subService, logger)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:         defaultHTTPAddr,
		Handler:      mux,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
		return nil
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	}
}
