package history

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type historyStorage struct {
	rs *redis.Client
	l  *slog.Logger
}

func New(rs *redis.Client, l *slog.Logger) *historyStorage {
	return &historyStorage{rs: rs, l: l}
}
