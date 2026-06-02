package user_subscriptions

import (
	"database/sql"
	"log/slog"
)

type userSubscriptionsStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *userSubscriptionsStorage {
	return &userSubscriptionsStorage{db: db, l: l}
}
