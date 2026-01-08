package app

import (
	"compConfigurator/internal/config"
	"compConfigurator/internal/db"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run(ctx context.Context, cfg config.Config) error {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	pool, err := db.New(ctx, cfg.DB.DSN, cfg.DB.MinConns, cfg.DB.MaxConns, cfg.DB.ConnMaxIdle, cfg.DB.ConnMaxLife)
	if err != nil {
		return err
	}
	defer pool.Close()

	handler := NewRouter(cfg, pool)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      handler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	errChan := make(chan error, 1)
	go func() {
		log.Info().Str("addr", srv.Addr).Msg("http server started")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()

		log.Info().Msg("shutting down http server")
		_ = srv.Shutdown(shutdownCtx)
		return nil
	case err := <-errChan:
		return err
	}
}
