package reading_sessions

import (
	"database/sql"
	"log/slog"
)

type readingSessionsStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *readingSessionsStorage {
	return &readingSessionsStorage{db: db, l: l}
}
