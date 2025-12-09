package store

import (
	"database/sql"
	"embed"
	"log"

	_ "modernc.org/sqlite"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

const dsn = "file:data.db?_foreign_keys=on&_journal_mode=WAL"

// dùng sẵn *sql.DB, không mở mới
func migrateDB(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		src,
		"sqlite",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("✅ Migrate done (hoặc không có gì mới)")
	return nil
}

func NewSQLiteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// optional: set busy timeout cho SQLite đỡ báo locked
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000;`); err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// migrate ngay trên connection này
	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}
