package fav_items

import (
	"database/sql"
	"log/slog"
)

type fiStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *fiStorage {
	return &fiStorage{db: db, l: l}
}
