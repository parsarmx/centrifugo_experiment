package service

import (
	"golang_template/internal/config"
	"golang_template/internal/producers"
	"golang_template/internal/repository"
	rpc_service "golang_template/proto"

	"go.uber.org/zap"
)

type Service interface {
	UserService() UserService
	RoomService() RoomService
}

type service struct {
	userService UserService
	roomService RoomService
}

func NewService(
	repo repository.Repository,
	logger *zap.Logger,
	redis producers.RedisClient,
	config config.Config,
	grpc rpc_service.CentrifugoApiClient,
) Service {
	userService := NewUserService(repo.UserRepository(), logger, redis, config.Auth)
	roomService := NewRoomService(repo.RoomRepository(), logger, grpc)
	return &service{
		userService: userService,
		roomService: roomService,
	}
}

func (s *service) UserService() UserService {
	return s.userService
}

func (s *service) RoomService() RoomService {
	return s.roomService
}
