package ai_chat

import (
	"database/sql"
	"log/slog"
)

type aiChatStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *aiChatStorage {
	return &aiChatStorage{db: db, l: l}
}
