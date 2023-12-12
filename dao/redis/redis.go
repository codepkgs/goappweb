package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/x-hezhang/gowebapp/settings"
)

var rdb *redis.Client

func Init() (err error) {
	var ctx context.Context
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", settings.Conf.RedisConfig.Host, settings.Conf.RedisConfig.Port),
		Username: settings.Conf.RedisConfig.Username,
		Password: settings.Conf.RedisConfig.Password,
		DB:       settings.Conf.RedisConfig.Database,
		PoolSize: settings.Conf.RedisConfig.PoolSize,
	})
	_, err = rdb.Ping(ctx).Result()
	return
}

func Close() {
	_ = rdb.Close()
}
