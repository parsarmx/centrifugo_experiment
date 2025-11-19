package application

import (
	"golang_template/internal/producers"

	"go.uber.org/zap"
)

func (a *application) InitRedis(logger *zap.Logger) producers.RedisClient {
	return producers.NewRedis(&a.config.Redis, logger)
}
