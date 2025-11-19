package application

import (
	"golang_template/internal/repository"
	"golang_template/internal/service"

	"go.uber.org/zap"
)

func (a *application) InitService(repo repository.Repository, logger *zap.Logger) service.Service {
	return service.NewService(repo, logger)
}
