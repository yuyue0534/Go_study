package database

import (
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed migrations/001_init.sql
var initSchema string

func Migrate(db *sql.DB) error {
	if _, err := db.Exec(initSchema); err != nil {
		return fmt.Errorf("execute initial schema: %w", err)
	}
	return nil
}
