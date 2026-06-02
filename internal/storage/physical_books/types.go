package physical_books

import (
	"database/sql"
	"log/slog"
)

type physicalBooksStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *physicalBooksStorage {
	return &physicalBooksStorage{db: db, l: l}
}
