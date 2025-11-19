package router

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Router interface {
	AddRoutes(e *echo.Echo)
}

type router struct {
	logger *zap.Logger
}

func NewRouter(logger *zap.Logger) Router {
	return &router{logger: logger}
}

func (r *router) AddRoutes(e *echo.Echo) {
	apiGroup := e.Group("/api/v0")
	apiGroup.GET("/heath-check", nil)
}
