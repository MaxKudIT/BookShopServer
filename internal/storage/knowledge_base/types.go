package knowledge_base

import (
	"database/sql"
	"log/slog"
)

type netStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *netStorage {
	return &netStorage{db: db, l: l}
}
