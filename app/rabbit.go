package application

import (
	"golang_template/internal/broker"

	"go.uber.org/zap"
)

func (a *application) InitRabbit(logger *zap.Logger) broker.Rabbit {
	mq, err := broker.NewRabbit(&a.config.Rabbit, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return mq
}
