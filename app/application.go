package application

import (
	"context"
	"fmt"
	"golang_template/internal/broker"
	"golang_template/internal/config"

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

func (a *application) RunEchoServer(lc fx.Lifecycle, e *echo.Echo, logger *zap.Logger) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("starting server...")

				go func() {
					e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", a.config.Server.Host, a.config.Server.Port)))
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})
}

func (a *application) RunRabbitConnection(lc fx.Lifecycle, r broker.Rabbit, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("rabbitmq connected")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing rabbitmq connection...")
			r.Close()
			return nil
		},
	})
}

func (a *application) Setup() {

	app := fx.New(
		fx.Provide(
			a.InitRepository,
			a.InitDatabase,
			a.InitHandler,
			a.InitLogger,
			a.InitFramework,
			a.InitRabbit,
		),
		fx.Invoke(a.RunEchoServer),
		fx.Invoke(a.RunRabbitConnection),
	)

	app.Run()
}
