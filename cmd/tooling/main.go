package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"cmsc.group2.coffee-api/cmd/tooling/cmd"
	"cmsc.group2.coffee-api/internal/schema"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// This tool is used to seed the database

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// The Darwin lib reqiures that pgx follow database/sql interfaces.
	// db, err := sql.Open("pgx", "user=postgres password=p@55word123 host=localhost port=5432 database=postgres sslmode=disable")
	db_dsn := os.Getenv("DB_DSN")
	db, err := sql.Open("pgx", db_dsn)
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

	// Generate an admin password hash
	// here we use 'password123$' as default
	hash, err := cmd.CreateAdminPassword()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// here we create an annonymous struct
	admin := struct {
		name  string
		email string
		hash  []byte
	}{
		name:  "admin",
		email: "admin@coffeeshop.com",
		hash:  hash,
	}

	// Hard coded and in need of some love at a later date.
	adminSeed := `INSERT into users (name, email, password_hash) VALUES
    ($1,$2,$3);`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, adminSeed, admin.name, admin.email, admin.hash)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
