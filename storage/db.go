package storage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB() (*pgxpool.Pool, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://trellouser:wzQ%7CjouYF1HeSp%2AIeb2nKKS%24JZeAn%40k%7C1nszW9A2qRwHHmtgv%2As9bT%3Fj%406pCuMNuFugu@localhost:5432/trello_db?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return pool, nil
}

func Migrate(pool *pgxpool.Pool) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS boards (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(100) NOT NULL,
		description VARCHAR(500) DEFAULT '',
		is_public BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS lists (
		id UUID PRIMARY KEY,
		board_id UUID NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
		name VARCHAR(100) NOT NULL,
		position INT DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS cards (
		id UUID PRIMARY KEY,
		list_id UUID NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
		title VARCHAR(200) NOT NULL,
		description VARCHAR(1000) DEFAULT '',
		position INT DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);
	`

	_, err := pool.Exec(context.Background(), schema)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	_, _ = pool.Exec(context.Background(), `ALTER TABLE boards ADD COLUMN IF NOT EXISTS is_public BOOLEAN DEFAULT FALSE`)

	log.Println("Database migration completed")
	return nil
}
