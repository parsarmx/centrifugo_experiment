package service

import (
	"golang_template/internal/repository"

	"go.uber.org/zap"
)

type Service interface {
}

type service struct {
}

func NewService(repo repository.Repository, logger *zap.Logger) Service {
	return &service{}
}
