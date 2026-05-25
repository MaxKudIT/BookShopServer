package history

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type hRedisStorage struct {
	rs *redis.Client
	l  *slog.Logger
}

func New(rs *redis.Client, l *slog.Logger) *hRedisStorage {
	return &hRedisStorage{rs: rs, l: l}
}
