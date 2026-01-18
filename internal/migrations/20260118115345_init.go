package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInit, downInit)
}

func upInit(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS secrets (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		type TEXT NOT NULL,
		title TEXT NOT NULL,
		data BYTEA NOT NULL,
		metadata TEXT,
		version INTEGER DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS sync_log (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		secret_id TEXT NOT NULL,
		operation TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_secrets_user_id ON secrets(user_id);
	CREATE INDEX IF NOT EXISTS idx_sync_log_user_id ON sync_log(user_id);
	CREATE INDEX IF NOT EXISTS idx_sync_log_timestamp ON sync_log(timestamp);
	`

	_, err := tx.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema creation: %w", err)
	}
	return nil
}

func downInit(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
	DROP TABLE IF EXISTS sync_log CASCADE;
	DROP TABLE IF EXISTS secrets CASCADE;
	DROP TABLE IF EXISTS users CASCADE;
	`)
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	return nil
}
