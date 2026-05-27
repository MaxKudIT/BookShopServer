package recommendation

import (
	"database/sql"
	"log/slog"
)

type recommendationStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *recommendationStorage {
	return &recommendationStorage{db: db, l: l}
}
