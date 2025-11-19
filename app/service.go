package application

import (
	"golang_template/internal/producers"
	"golang_template/internal/repository"
	"golang_template/internal/service"
	rpc_service "golang_template/proto"

	"go.uber.org/zap"
)

func (a *application) InitService(repo repository.Repository, logger *zap.Logger, redis producers.RedisClient, grpc rpc_service.CentrifugoApiClient) service.Service {
	return service.NewService(repo, logger, redis, *a.config, grpc)
}
