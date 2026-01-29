package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/config"
)

var (
	redisCacheKey = "reporting:"
	ttl           = time.Duration(5) * time.Minute
)

type Redis struct{ Client *redis.Client }

func NewRedis(cfg config.Config) *Redis {
	addr := cfg.RedisAddr
	// if addr == "" {
	// 	addr = "localhost:6379"
	// }
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	fmt.Println("Connected to Redis successfully")
	return &Redis{Client: rdb}
}

func (r *Redis) CheckInCache(ctx context.Context, exportID string) (*ExportResponse, error) {
	// check if exportID status in cache
	// if processing then return processing with code
	// else if status is completed then return the byte
	// if exportID is not present then return with no data in cache
	key := redisCacheKey + exportID
	var cacheValue ExportResponse

	data, _ := r.Client.Get(ctx, key).Result()

	// if err != nil {
	// 	log.Printf("Error while getting data cache with key: %s %s", key, err)
	// 	return nil, err
	// }
	if data == "" {
		log.Printf("No data in cache with key: %s", key)
		return nil, nil
	}

	if err := json.Unmarshal([]byte(data), &cacheValue); err != nil {
		log.Println("Error unmarshalling data from cache.")
		return nil, errors.New("error unmarshalling data from cache")
	}

	return &cacheValue, nil
}

func (r *Redis) SetInCache(ctx context.Context, exportID string, data []byte) error {

	key := redisCacheKey + exportID
	return r.Client.Set(ctx, key, data, ttl).Err()
}
