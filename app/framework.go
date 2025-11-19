package application

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (a *application) InitFramework(logger *zap.Logger) *echo.Echo {
	return echo.New()
}
