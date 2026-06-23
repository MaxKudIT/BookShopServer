package admin

import (
	"database/sql"
	"log/slog"
)

type adminStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *adminStorage {
	return &adminStorage{db: db, l: l}
}
