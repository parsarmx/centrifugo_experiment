package service

import (
	"context"
	"time"

	"golang_template/internal/repository"
	"golang_template/internal/repository/models"

	"go.uber.org/zap"
)

const (
	roomServiceDefaultTimeout = 10 * time.Second
)

type RoomService interface {
	CreateRoom(ctx context.Context, roomName string, capacity int) (*models.Room, error)
}

type roomService struct {
	repo   repository.RoomRepository
	logger *zap.Logger
}

func NewRoomService(
	repo repository.RoomRepository,
	logger *zap.Logger,
) RoomService {
	return &roomService{
		repo:   repo,
		logger: logger,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, roomName string, capacity int) (*models.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, roomServiceDefaultTimeout)
	defer cancel()

	if roomName == "" {
		s.logger.Debug("Room name is empty")
		return nil, ErrInvalidRoomName
	}

	if capacity <= 0 {
		s.logger.Debug("Invalid room capacity",
			zap.Int("capacity", capacity))
		return nil, ErrInvalidRoomCapacity
	}

	room, err := s.repo.CreateRoom(ctx, roomName, capacity)
	if err != nil {
		s.logger.Error("Failed to create room",
			zap.String("room_name", roomName),
			zap.Int("capacity", capacity),
			zap.Error(err))
		return nil, err
	}

	s.logger.Debug("Room created successfully in service",
		zap.String("room_id", room.ID.String()),
		zap.String("channel", room.Channel))

	return room, nil
}
