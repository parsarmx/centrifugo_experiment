package application

import (
	"golang_template/internal/producers"
	"golang_template/internal/repository"
	"golang_template/internal/service"

	"go.uber.org/zap"
)

func (a *application) InitService(repo repository.Repository, logger *zap.Logger, redis producers.RedisClient) service.Service {
	return service.NewService(repo, logger, redis, *a.config)
}
