package users_books

import (
	"database/sql"
	"log/slog"
)

type ubStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *ubStorage {
	return &ubStorage{db: db, l: l}
}
