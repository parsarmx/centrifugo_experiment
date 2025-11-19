package repository

import (
	"context"
	"golang_template/internal/database/postgres"

	"go.uber.org/zap"
)

type Repository interface {
	UserRepository() UserRepository
	RoomRepository() RoomRepository
}

type repository struct {
	userRepository UserRepository
	roomRepository RoomRepository
}

func NewRepository(ctx context.Context, db postgres.Database, logger *zap.Logger) Repository {
	userRepository := NewUserRepository(db, logger)
	roomRepository := NewRoomRepository(db, logger)
	return &repository{
		userRepository: userRepository,
		roomRepository: roomRepository,
	}
}

func (r *repository) UserRepository() UserRepository {
	return r.userRepository
}

func (r *repository) RoomRepository() RoomRepository {
	return r.roomRepository
}
