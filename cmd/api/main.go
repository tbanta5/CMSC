package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Passed from Makefile, through dockerfile
const (
	VERSION = "0.0.1"
	PORT    = "8585"
	ENV     = "develop"
)

type config struct {
	port    string
	env     string
	version string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	// Instantiate the config struct
	var cfg config
	cfg.env = ENV
	cfg.port = PORT
	cfg.version = VERSION

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}
