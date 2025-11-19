package application

import (
	"golang_template/api/controllers"
	"golang_template/internal/service"

	"go.uber.org/zap"
)

func (a *application) InitContoller(service service.Service, logger *zap.Logger) controllers.Controllers {
	return controllers.NewController(service, logger, *a.config)
}
