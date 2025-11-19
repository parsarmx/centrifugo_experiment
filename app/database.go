package application

import (
	"golang_template/internal/database/postgres"

	"go.uber.org/zap"
)

func (a *application) InitDatabase(logger *zap.Logger) postgres.Database {
	db, err := postgres.NewDatabase(a.ctx, &a.config.DB, logger)
	if err != nil {
		logger.Fatal("failed to connect to Database", zap.Error(err))
	}

	return db
}
