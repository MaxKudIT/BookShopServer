package page

import (
	"database/sql"
	"log/slog"
)

type pageStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *pageStorage {
	return &pageStorage{db: db, l: l}
}
