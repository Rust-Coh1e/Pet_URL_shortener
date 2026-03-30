package repository //наверно он должен называться repository

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
	"url-shortener/storage"
)

type DB struct {
	conn *sql.DB
}

func NewDB(connString string) (*DB, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	query := `
    CREATE TABLE IF NOT EXISTS urls (
        id          TEXT PRIMARY KEY,
        original    TEXT NOT NULL,
        clicks      BIGINT DEFAULT 0,
        created_at  TIMESTAMP DEFAULT NOW()
    );
    `

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return &DB{conn: db}, nil
}

func (db *DB) Save(url storage.URL) error {
	query := `INSERT INTO urls (id, original, created_at) VALUES ($1, $2, $3)`
	if _, err := db.conn.Exec(query, url.ID, url.Original, url.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (db *DB) Get(id string) (*storage.URL, error) {
	var url storage.URL

	query := `SELECT id, original, clicks, created_at FROM urls WHERE id = $1`

	err := db.conn.QueryRow(query, id).Scan(
		&url.ID,
		&url.Original,
		&url.Clicks,
		&url.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // не найдено — не ошибка
	}

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (db *DB) IncrementClick(id string) error {
	query := `UPDATE urls SET clicks = clicks + 1 WHERE id = $1`

	if _, err := db.conn.Exec(query, id); err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() {
	db.conn.Close()
}
