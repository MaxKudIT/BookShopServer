package orders

import (
	"database/sql"
	"log/slog"
)

type ordersStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *ordersStorage {
	return &ordersStorage{db: db, l: l}
}
