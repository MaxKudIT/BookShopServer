package bookviews

import (
	"database/sql"
	"log/slog"
)

type bookViewsStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *bookViewsStorage {
	return &bookViewsStorage{db: db, l: l}
}
