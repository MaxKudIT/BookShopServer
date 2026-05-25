package book_revs

import (
	"database/sql"
	"log/slog"
)

type bookRevsStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *bookRevsStorage {
	return &bookRevsStorage{db: db, l: l}
}
