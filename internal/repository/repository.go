package repository

import (
	"context"
	"golang_template/internal/database/postgres"

	"go.uber.org/zap"
)

type Repository interface {
}

type repository struct {
	ctx context.Context
	db  postgres.Database
}

func NewRepository(ctx context.Context, db postgres.Database, logger *zap.Logger) Repository {
	return &repository{
		ctx: ctx,
		db:  db,
	}
}
