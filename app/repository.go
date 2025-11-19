package application

import (
	"golang_template/internal/database/postgres"
	"golang_template/internal/repository"
	rpc_service "golang_template/proto"

	"go.uber.org/zap"
)

func (a *application) InitRepository(db postgres.Database, logger *zap.Logger, grpc rpc_service.CentrifugoApiClient) repository.Repository {
	return repository.NewRepository(a.ctx, db, logger, grpc)
}
