package reading

import (
	"database/sql"
	"log/slog"
)

type readingStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *readingStorage {
	return &readingStorage{db: db, l: l}
}
