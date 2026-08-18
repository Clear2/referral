// Package redis provides the application's instrumented Redis client.
package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9"
	goredis "github.com/redis/go-redis/v9"
)

const defaultMaxPoolSize = 10

// Redis owns the underlying go-redis client.
type Redis struct {
	maxPoolSize int
	password    string

	Client *goredis.Client
}

// New creates, verifies, and instruments a Redis client.
func New(url string, opts ...Option) (*Redis, error) {
	client := &Redis{maxPoolSize: defaultMaxPoolSize}
	for _, opt := range opts {
		opt(client)
	}

	config, err := goredis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}
	config.PoolSize = client.maxPoolSize
	if client.password != "" {
		config.Password = client.password
	}
	client.Client = goredis.NewClient(config)

	if err = client.Client.Ping(context.Background()).Err(); err != nil {
		_ = client.Client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	if err = errors.Join(
		redisotel.InstrumentTracing(client.Client),
		redisotel.InstrumentMetrics(client.Client),
	); err != nil {
		_ = client.Client.Close()
		return nil, fmt.Errorf("instrument redis: %w", err)
	}
	return client, nil
}

// Close releases the Redis client's resources.
func (client *Redis) Close() error {
	if client == nil || client.Client == nil {
		return nil
	}
	return client.Client.Close()
}
