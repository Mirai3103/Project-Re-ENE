package store

import (
	"embed"

	"database/sql"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	_ "github.com/tursodatabase/turso-go"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func NewDB() (*sqlx.DB, error) {
	// mở SQLite
	db, err := sql.Open("turso", "./data.db")
	if err != nil {
		return nil, err
	}

	// chạy migration runtime
	source := migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationFiles,
		Root:       "migrations",
	}
	_, err = migrate.Exec(db, "sqlite3", source, migrate.Up)
	if err != nil {
		return nil, err
	}
	dbx := sqlx.NewDb(db, "turso")

	return dbx, nil
}
