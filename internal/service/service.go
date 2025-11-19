package service

import (
	"golang_template/internal/config"
	"golang_template/internal/producers"
	"golang_template/internal/repository"

	"go.uber.org/zap"
)

type Service interface {
	UserService() UserService
}

type service struct {
	userService UserService
}

func NewService(
	repo repository.Repository,
	logger *zap.Logger,
	redis producers.RedisClient,
	config config.Config,
) Service {
	userService := NewUserService(repo.UserRepository(), logger, redis, config.Auth)
	return &service{
		userService: userService,
	}
}

func (s *service) UserService() UserService {
	return s.userService
}
