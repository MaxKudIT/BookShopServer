package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connection(redisOptions *redis.Options) (*redis.Client, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	rdb := redis.NewClient(redisOptions)

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Redis:", pong)

	return rdb, nil
}
