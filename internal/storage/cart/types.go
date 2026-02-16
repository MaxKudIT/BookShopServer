package cart

import (
	"database/sql"
	"log/slog"
)

type cStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *cStorage {
	return &cStorage{db: db, l: l}
}
