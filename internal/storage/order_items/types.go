package order_items

import (
	"database/sql"
	"log/slog"
)

type orderItemsStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *orderItemsStorage {
	return &orderItemsStorage{db: db, l: l}
}
