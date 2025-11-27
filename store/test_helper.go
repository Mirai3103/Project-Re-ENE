package store

import (
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	_ "github.com/tursodatabase/turso-go"
)

// SetupTestDB creates an in-memory database for testing
func SetupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	// Create in-memory SQLite database
	db, err := sql.Open("turso", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Run migrations
	source := migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationFiles,
		Root:       "migrations",
	}
	_, err = migrate.Exec(db, "sqlite3", source, migrate.Up)
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	dbx := sqlx.NewDb(db, "turso")
	return dbx
}

// TeardownTestDB closes the test database
func TeardownTestDB(t *testing.T, db *sqlx.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Errorf("Failed to close test database: %v", err)
	}
}
