package application

import (
	"context"
	"fmt"
	"golang_template/api/router"
	"golang_template/internal/config"
	"golang_template/internal/database/postgres"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Application interface {
	Setup()
}

type application struct {
	ctx    context.Context
	config *config.Config
}

func NewApplication(ctx context.Context, config *config.Config) Application {
	return &application{ctx: ctx, config: config}
}

func (a *application) Setup() {

	app := fx.New(
		fx.Provide(
			a.InitRepository,
			a.InitDatabase,
			a.InitContoller,
			a.InitService,
			a.InitLogger,
			a.InitFramework,
			a.InitRedis,
			a.InitRouter,
			a.InitGRPCServer,
		),
		postgres.Module,
		fx.Invoke(func(lifecycle fx.Lifecycle, e *echo.Echo, logger *zap.Logger, router router.Router) {
			lifecycle.Append(
				fx.Hook{
					OnStart: func(ctx context.Context) error {
						logger.Info("starting server...")
						router.AddRoutes(e)
						go func() {
							if err := e.Start(fmt.Sprintf("%v:%v", a.config.Server.Host, a.config.Server.Port)); err != nil && err != http.ErrServerClosed {
								logger.Fatal("shutting down server", zap.Error(err))
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return nil
					},
				})
		}),
	)

	app.Run()
}
