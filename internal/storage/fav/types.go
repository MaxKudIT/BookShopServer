package fav

import (
	"database/sql"
	"log/slog"
)

type fStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *fStorage {
	return &fStorage{db: db, l: l}
}
