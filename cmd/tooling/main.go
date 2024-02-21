package main

import (
	"database/sql"
	"log/slog"
	"os"

	"cmsc.group2.coffee-api/internal/schema"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// This tool is used to seed the database

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// The Darwin lib reqiures that pgx follow database/sql interfaces.
	db, err := sql.Open("pgx", "user=postgres password=p@55word123 host=localhost port=5432 database=postgres sslmode=disable")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1) // Exit if there is a problem with DB connection
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1) // Exit if cannot ping DB
	}

	// Perform migrations
	err = schema.Migrate(db)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Migrations Complete")

	// Seed the database
	err = schema.Seed(db)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Seed Complete")
}
