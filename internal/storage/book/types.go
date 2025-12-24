package book

import (
	"database/sql"
	"log/slog"
)

type bookStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *bookStorage {
	return &bookStorage{db: db, l: l}
}
