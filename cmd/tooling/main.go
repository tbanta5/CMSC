package main

import (
	"log/slog"
	"os"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"cmsc.group2.coffee-api/internal/schema"
)

// This tool is used to seed the database

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// The Darwin lib reqiures that pgx follow database/sql interfaces.
	db, err := sql.Open("pgx", "user=postgres password=p@55word123 host=localhost port=5432 database=postgres sslmode=disable")
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Error(err.Error())
	}

	err = schema.Migrate(db)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Migrations Complete")

	err = schema.Seed(db)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Seed Complete")

}
