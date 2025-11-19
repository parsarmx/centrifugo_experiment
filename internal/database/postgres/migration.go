package postgres

import (
	"golang_template/internal/repository/models"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Invoke(Run)

func Run(db Database, logger *zap.Logger) {
	logger.Info("Running migrations...")
	if err := db.Gorm().AutoMigrate(
		&models.User{},
	); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}
}
