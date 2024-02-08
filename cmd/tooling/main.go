package main

import (
	"context"
	"time"

	"log/slog"
	"os"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"cmsc.group2.coffee-api/internal/schema"
)

// This tool is used to seed the database

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Note terribly clear hear but due to darwin lib, must use pgx adapter or just database/sql
	// db, err := sql.Open("pgx", "postgres://postgres:p@55word123@localhost/postgres?sslmode=disable")
	db, err := sql.Open("pgx", "user=postgres password=p@55word123 host=localhost port=5432 database=postgres sslmode=disable")
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.Ping()
	if err != nil {
		logger.Error(err.Error())
	}

	schema.PrintSchema()

	err = schema.Migrate(ctx, db)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("Migrations Complete")

}
