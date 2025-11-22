package service

import (
	"context"
	"fmt"
	"time"

	"golang_template/internal/repository"
	"golang_template/internal/repository/models"
	rpc_service "golang_template/proto"

	"go.uber.org/zap"
)

const (
	roomServiceDefaultTimeout = 10 * time.Second
)

type RoomService interface {
	CreateRoom(ctx context.Context, roomName string, capacity int) (*models.Room, error)
	SendMessage(ctx context.Context, userId, message, channel string)
	GetAllRooms(ctx context.Context) ([]*models.Room, error)
}

type roomService struct {
	repo   repository.RoomRepository
	logger *zap.Logger
	grpc   rpc_service.CentrifugoApiClient
}

func NewRoomService(
	repo repository.RoomRepository,
	logger *zap.Logger,
	grpc rpc_service.CentrifugoApiClient,
) RoomService {
	return &roomService{
		repo:   repo,
		logger: logger,
		grpc:   grpc,
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

func (s *roomService) GetAllRooms(ctx context.Context) ([]*models.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, roomServiceDefaultTimeout)
	defer cancel()

	rooms, err := s.repo.FetchAllRooms(ctx)
	if err != nil {
		s.logger.Error("Failed to Fetch rooms",
			zap.Error(err))
		return nil, err
	}

	return rooms, nil
}

func (s *roomService) GetRoomsById(ctx context.Context, roomId string) (models.Room, error) {
	return models.Room{}, nil
}

// it needs at least an error
func (s *roomService) SendMessage(ctx context.Context, userId, message, channel string) {
	// maybe add a repo to store messages inside db ...

	// Run a gRPC‌ request to publish joining in channel
	go func() { // This should place inside service, anyway...
		req := &rpc_service.PublishRequest{
			Channel:     channel,
			Data:        []byte(fmt.Sprintf(`{"user_id": "%s", "message":"%s"}`, userId, message)),
			SkipHistory: false,
			Tags: map[string]string{
				"source": "room_service",
			},
		}

		fmt.Println(s.grpc)
		resp, err := s.grpc.Publish(context.Background(), req)
		if err != nil {
			fmt.Println(err)
			// fmt.Println(resp.Error.Code, resp.Error.Message)
			s.logger.Error("Failed to publish room creation", zap.Error(err))
			return
		}

		if resp.Error != nil && resp.Error.Code != 0 {
			fmt.Println(resp.Error.Code, resp.Error.Message)
			s.logger.Error("Centrifugo returned an error",
				zap.Uint32("code", resp.Error.Code),
				zap.String("msg", resp.Error.Message),
			)
			return
		}

	}()
	// it needs at least an error
}
