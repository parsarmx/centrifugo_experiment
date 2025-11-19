package application

import (
	"golang_template/internal/database/postgres"
	"golang_template/internal/repository"

	"go.uber.org/zap"
)

func (a *application) InitRepository(db postgres.Database, logger *zap.Logger) repository.Repository {
	return repository.NewRepository(a.ctx, db, logger)
}
