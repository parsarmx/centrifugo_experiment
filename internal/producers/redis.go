package producers

import (
	"context"
	"errors"
	"fmt"
	err "golang_template/api/errors"
	"golang_template/internal/config"
	"time"

	"github.com/gofiber/storage/redis/v3"
	go_redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	ErrFailedToRetrieve = &err.AppError{
		Msg: "failed to get TTL",
		Err: errors.New("failed to get TTL"),
	}
)

type RedisClient interface {
	Conn() *go_redis.Client
	RedisStorage() *redis.Storage
	Close() error
	TTL(ctx context.Context, key string) (time.Duration, error)
}

type Redis struct {
	client *go_redis.Client
	store  *redis.Storage
	logger *zap.Logger
}

func NewRedis(cfg *config.RedisConfig, logger *zap.Logger) RedisClient {
	// Initialize go-redis client
	client := go_redis.NewClient(&go_redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), // Host:Port
		Password: cfg.Password,                             // Password
		DB:       cfg.DB,                                   // Database
		PoolSize: cfg.PoolSize,                             // Pool size
	})

	// Initialize redis.Storage for middleware
	store := redis.New(redis.Config{
		Host:      cfg.Host,
		Port:      cfg.Port,
		Password:  cfg.Password,
		Database:  cfg.DB,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  cfg.PoolSize,
	})

	return &Redis{
		client: client,
		store:  store,
		logger: logger,
	}
}

func (r *Redis) Conn() *go_redis.Client {
	return r.client
}

func (r *Redis) RedisStorage() *redis.Storage {
	return r.store
}

func (r *Redis) Close() error {
	// Close both the go-redis client and redis.Storage
	if err := r.client.Close(); err != nil {
		return err
	}
	return r.store.Close()
}

func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttlCmd := r.client.TTL(ctx, key)
	ttlDuration, err := ttlCmd.Result()
	if err != nil {
		return time.Duration(0), ErrFailedToRetrieve
	}

	if ttlDuration < 0 {
		return time.Duration(0), fmt.Errorf("key %s does not exist or has no expiration", key)
	}

	return ttlDuration, nil
}
