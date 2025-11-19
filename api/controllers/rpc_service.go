package controllers

import (
	"golang_template/internal/service"

	"go.uber.org/zap"
)

type RpcServiceController interface {
}

type rpcServiceController struct {
	rpcServiceService service.RpcServiceService
}

func NewRpcServiceController(service service.RpcServiceService, logger *zap.Logger) RpcServiceController {
	return &rpcServiceController{
		rpcServiceService: service,
	}
}
