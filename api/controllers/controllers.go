package controllers

import (
	"golang_template/internal/config"
	"golang_template/internal/service"

	"go.uber.org/zap"
)

type Controllers interface {
	UserController() UserController
}

type controllers struct {
	userController UserController
}

func NewController(s service.Service, logger *zap.Logger, config config.Config) Controllers {
	userController := NewUserController(s.UserService(), logger, config.Auth)
	return &controllers{
		userController: userController,
	}
}

func (c *controllers) UserController() UserController {
	return c.userController
}
