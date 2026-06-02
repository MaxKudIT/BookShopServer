package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connection(redisOptions *redis.Options) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rdb := redis.NewClient(redisOptions)

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	return rdb, nil
}
