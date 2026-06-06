package subscription_plans

import (
	"database/sql"
	"log/slog"
)

type subscriptionPlansStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *subscriptionPlansStorage {
	return &subscriptionPlansStorage{db: db, l: l}
}
