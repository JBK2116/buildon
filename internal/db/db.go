// Package db provides database functionality in the project.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
)

var (
	once   sync.Once
	dbConn *sql.DB
	dbErr  error
)

const (
	maxOpenConns    = 1
	maxIdleConns    = 1
	maxConnLifetime = 0 // process terminates connection on exit.
	timeoutDB       = 5 * time.Second
)

// newDB opens a database connection to the sqlite db,
// creating the db file if it does not exist.
func newDB(logger *slog.Logger) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()
	const path = "./app.db"
	const dsn = "file:" + path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)&_pragma=cache_size(-64000)&_pragma=foreign_keys(ON)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxConnLifetime)

	if errPing := db.PingContext(ctx); errPing != nil {
		logger.Error("failed to ping database", "error", errPing)
		return nil, err
	}
	return db, nil
}

const schema = `
CREATE TABLE IF NOT EXISTS projects (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	title      TEXT    NOT NULL,
	created_at TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
	updated_at TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS problems (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	title      TEXT    NOT NULL,
	content    TEXT    NOT NULL,
	solved     INTEGER NOT NULL DEFAULT 0,
	created_at TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
	updated_at TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
	FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_problems_project_id ON problems (project_id);
`

// schemaVersion is bumped whenever the schema changes, tracking applied migrations
// via SQLite's PRAGMA user_version.
const schemaVersion = 2

// initializeSchema applies the schema to the database only if it has not already
// been applied. It uses PRAGMA user_version to detect whether the current schema
// version has been executed, and returns whether the schema was (re)applied this run
// plus any error encountered.
func initializeSchema(db *sql.DB, logger *slog.Logger) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		logger.Error("failed to query user_version", "error", err)
		return false, err
	}

	if version >= schemaVersion {
		logger.Debug("no migrations to run")
		return false, nil
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		logger.Error("failed to run migrations", "error", err)
		return false, err
	}

	// Alter existing tables created by older schema versions that predate the
	// migrations below. These are safe to run on fresh databases because the
	// column already exists, in which case the "duplicate column" error is
	// tolerated.
	alterations := []string{
		"ALTER TABLE problems ADD COLUMN solved INTEGER NOT NULL DEFAULT 0",
	}
	for _, alter := range alterations {
		if _, err := db.ExecContext(ctx, alter); err != nil && !strings.Contains(err.Error(), "duplicate column") {
			logger.Error("failed to run table alteration", "error", err)
			return false, err
		}
	}

	if _, err := db.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", schemaVersion)); err != nil {
		logger.Error("failed to confirm migration execution", "error", err)
		return false, err
	}
	logger.Debug("database migrations ran", "version", schemaVersion)
	return true, nil
}

// OpenConn initializes the database connection and applies the schema exactly once per
// process, using [sync.Once]. Subsequent calls return the cached [*sql.DB] and any error
// from the first invocation.
func OpenConn(logger *slog.Logger) (*sql.DB, error) {
	once.Do(func() {
		dbConn, dbErr = newDB(logger)
		if dbErr != nil {
			return
		}
		_, dbErr = initializeSchema(dbConn, logger)
	})
	return dbConn, dbErr
}
