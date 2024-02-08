package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
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
	db      struct {
		dsn string
	}
}

type application struct {
	config         config
	logger         *slog.Logger
	sessionManager *scs.SessionManager
}

func main() {
	// Instantiate the config struct
	var cfg config
	cfg.env = ENV
	cfg.port = PORT
	cfg.version = VERSION
	cfg.db.dsn = "postgres://postgres:pa55word123@localhost:5432/postgres?sslmode=disable"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Initialize a new session manager and configure it to use postgresstore as the session store.
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 5 * time.Minute

	app := &application{
		config:         cfg,
		logger:         logger,
		sessionManager: sessionManager,
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

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(cfg config) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
