package router

import (
	"golang_template/api/controllers"
	"golang_template/internal/config"
	"golang_template/internal/pkg/echojwt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Router interface {
	AddRoutes(e *echo.Echo)
}

type router struct {
	logger      *zap.Logger
	controllers controllers.Controllers
	config      config.Config
}

func NewRouter(controllers controllers.Controllers, logger *zap.Logger, config config.Config) Router {
	return &router{logger: logger, controllers: controllers, config: config}
}

func (r *router) AddRoutes(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			allowed := map[string]bool{
				"http://localhost:3000": true,
				"http://127.0.0.1:3000": true,
				"http://localhost:8080": true,
				"http://127.0.0.1:8080": true,
			}
			return allowed[origin], nil
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Origin",
			"access",
		},
		AllowMethods: []string{
			echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS,
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	apiGroup := e.Group("/api/v0/public")

	private := e.Group("api/v0/private", echojwt.New(
		echojwt.MiddlewareConfig{
			Secret:        r.config.Auth.JWTSecret,
			Logger:        r.logger,
			AuthHeaderKey: "access",
		},
	))

	apiGroup.POST("/otp/send", r.controllers.UserController().SendOTP)
	apiGroup.POST("/otp/login", r.controllers.UserController().OTPLogin)

	// private
	private.POST("/room/create", r.controllers.RoomController().CreateRoom)
	private.POST("/room/send_message", r.controllers.RoomController().SendMessage)

}
