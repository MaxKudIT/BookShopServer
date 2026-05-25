package stats

import (
	"database/sql"
	"log/slog"
)

type statsStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *statsStorage {
	return &statsStorage{db: db, l: l}
}
