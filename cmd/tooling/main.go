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
	// This allows us to add an initial admin password to the database.
	pwd := os.Getenv("ADMIN_PASSWD")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// The Darwin lib reqiures that pgx follow database/sql interfaces.
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
	// Test to ensure db is not already up.
	// This is useful in event of a pod deletion/death
	// since database seeding and migrations are done via
	// the intiContainer which runs for every newly started pod.
	// This action cause an ERROR originally in the db
	// assuming nothing has been migrated yet.
	stmt := `select coffee_id from coffee where coffee_id=1;`
	_, err = db.Exec(stmt)
	if err == nil {
		// Assume we got the statement back
		os.Exit(0)
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
	hash, err := cmd.CreateAdminPassword(pwd)
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
