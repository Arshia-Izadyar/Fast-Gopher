package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

var rClient *redis.Client

func InitRedis(cfg *config.Config) error {
	// TODO: config redis for performance
	rClient = redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:   cfg.Redis.Password,
		DB:         0,
		MaxRetries: 10,
	})
	_, err := rClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

func Set[T any](k string, v T, d time.Duration) error {
	value, err := sonic.Marshal(v)
	if err != nil {
		return err
	}
	_, err = rClient.Set(context.Background(), k, value, d).Result()
	if err != nil {
		return err
	}
	return nil
}

func Get[T any](k string) (*T, error) {
	dest := *new(T)
	v, err := rClient.Get(context.Background(), k).Result()
	if err != nil {
		return nil, err
	}
	err = sonic.Unmarshal([]byte(v), &dest)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func GetRedis() *redis.Client {
	return rClient
}

func Close() error {
	err := rClient.Close()
	if err != nil {
		return err
	}
	return nil
}
