package user

import (
	"database/sql"
	"log/slog"
)

type userStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *userStorage {
	return &userStorage{db: db, l: l}
}
