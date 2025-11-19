package controllers

import (
	"golang_template/internal/config"
	"golang_template/internal/service"

	"go.uber.org/zap"
)

type Controllers interface {
	UserController() UserController
	RoomController() RoomController
	// rpcServiceController() RpcServiceController
}

type controllers struct {
	userController UserController
	roomController RoomController
}

func NewController(s service.Service, logger *zap.Logger, config config.Config) Controllers {
	userController := NewUserController(s.UserService(), logger, config.Auth)
	roomController := NewRoomController(s.RoomService(), logger)
	return &controllers{
		userController: userController,
		roomController: roomController,
	}
}

func (c *controllers) UserController() UserController {
	return c.userController
}

func (c *controllers) RoomController() RoomController {
	return c.roomController
}

// func (c *controllers) RpcServiceController() RpcServiceController {
// 	return c.rpcServiceController
// }
