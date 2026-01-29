package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nationpulse-bff/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

type Redis struct{ Client *redis.Client }

func NewRedis(cfg config.Config) *Redis {
	addr := cfg.RedisAddr
	// if addr == "" {
	// 	addr = "localhost:6379"
	// }
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})
	fmt.Println("Connected to Redis successfully")
	return &Redis{Client: rdb}
}

var ttl time.Duration = time.Duration(15) * time.Minute

func (r *Redis) SetJTI(ctx context.Context, key, userID string, exp time.Time) error {
	return r.Client.Set(ctx, key, userID, time.Until(exp)).Err()
}

func (r *Redis) DelJTI(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) GetUserByJTI(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) SetData(ctx context.Context, key string, value interface{}) error {
	log.Println("Set cached data")
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) GetData(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) DelData(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
