package app

import (
	"fmt"

	"github.com/keep/sunny/config"
	"github.com/keep/sunny/pkg/logger"
	"github.com/keep/sunny/pkg/postgres"
	"github.com/keep/sunny/pkg/redis"
	"github.com/keep/sunny/pkg/rediskey"
	"go.uber.org/fx"
)

var infrastructureModule = fx.Module("infrastructure",
	fx.Provide(
		func(cfg config.Config) logger.Logger { return logger.New(cfg.Log().Dir, cfg.Log().Level) },
		fx.Annotate(newPostgres, fx.OnStop(closePostgres)),
		fx.Annotate(newRedis, fx.OnStop(closeRedis)),
		func(cfg config.Config) rediskey.RedisKey {
			return rediskey.New(cfg.Redis().Prefix, cfg.App().Env)
		},
	),
)

func newPostgres(cfg config.Config, log logger.Logger) (*postgres.Postgres, error) {
	pg, err := postgres.New(cfg.PG().URL, postgres.MaxPoolSize(cfg.PG().PoolMax))
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	log.Info("app - postgres connected")
	return pg, nil
}

func newRedis(cfg config.Config, log logger.Logger) (*redis.Redis, error) {
	redisConfig := cfg.Redis()
	client, err := redis.New(
		redisConfig.URL,
		redis.MaxPoolSize(redisConfig.PoolMax),
		redis.Password(redisConfig.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}
	log.Info("app - redis connected")
	return client, nil
}

func closeRedis(log logger.Logger, client *redis.Redis) error {
	log.Info("app - redis closed")
	return client.Close()
}

func closePostgres(log logger.Logger, pg *postgres.Postgres) error {
	log.Info("app - postgres closed")
	pg.Close()
	return nil
}
