package schema

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/ardanlabs/darwin"
	_ "github.com/jackc/pgx/v5/stdlib" // Must have because we are mirroring database/sql
)

var (
	//go:embed sql/schema.sql
	schemaDoc string
	//go:embed sql/delete.sql
	deleteDoc string
)

func Migrate(ctx context.Context, db *sql.DB) error {
	driver, err := darwin.NewGenericDriver(db, darwin.PostgresDialect{})
	if err != nil {
		return fmt.Errorf("construct darwin driver: %w", err)
	}

	d := darwin.New(driver, darwin.ParseMigrations(schemaDoc))
	return d.Migrate()
}

func PrintSchema() {
	fmt.Println(schemaDoc)
}

func DeleteAll(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()

}
