package application

import (
	"golang_template/api/controllers"
	"golang_template/api/router"

	"go.uber.org/zap"
)

func (a *application) InitRouter(controllers controllers.Controllers, logger *zap.Logger) router.Router {
	return router.NewRouter(controllers, logger, *a.config)
}
