package main

import (
	"Inquiro/config"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func Run(cfg config.Application, mux *chi.Mux) {
	srv := http.Server{
		Addr:    cfg.Config.Addr,
		Handler: mux,
	}
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		cfg.Logger.Infof("server has received %s signal", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()
	cfg.Logger.Infof("server has started on %s", cfg.Config.Addr)
	err := srv.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		cfg.Logger.Infof("server has stopped on %s", cfg.Config.Addr)
		return
	}
	if err != nil {
		cfg.Logger.Fatalf("failed to run server: %v", err.Error())
	}
}
