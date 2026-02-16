package cart_items

import (
	"database/sql"
	"log/slog"
)

type ciStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *ciStorage {
	return &ciStorage{db: db, l: l}
}
