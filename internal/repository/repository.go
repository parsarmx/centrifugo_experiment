package repository

import (
	"context"
	"golang_template/internal/database/postgres"

	"go.uber.org/zap"
)

type Repository interface {
	UserRepository() UserRepository
}

type repository struct {
	userRepository UserRepository
}

func NewRepository(ctx context.Context, db postgres.Database, logger *zap.Logger) Repository {
	userRepository := NewUserRepository(db, logger)
	return &repository{
		userRepository: userRepository,
	}
}

func (r *repository) UserRepository() UserRepository {
	return r.userRepository
}
