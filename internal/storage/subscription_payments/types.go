package subscription_payments

import (
	"database/sql"
	"log/slog"
)

type subscriptionPaymentsStorage struct {
	db *sql.DB
	l  *slog.Logger
}

func New(db *sql.DB, l *slog.Logger) *subscriptionPaymentsStorage {
	return &subscriptionPaymentsStorage{db: db, l: l}
}
