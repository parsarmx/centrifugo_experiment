package repository

import (
	"context"
	"errors"
	"golang_template/internal/database/postgres"
	"golang_template/internal/repository/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	UserDataByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error)
	UserCreate(ctx context.Context, phoneNumber string) (*models.User, error)
}

type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db postgres.Database, logger *zap.Logger) UserRepository {
	return &userRepository{
		db:     db.Gorm(),
		logger: logger,
	}
}

func (r userRepository) UserDataByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	var userData models.User

	err := r.db.WithContext(ctx).
		Select("id").
		Where("phone_number = ?", phoneNumber).
		First(&userData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("User not found while getting user by phone number",
				zap.String("phone_number", phoneNumber))
			return nil, err
		}

		r.logger.Error("Error while getting user by phone number", zap.Error(err))
		return nil, err
	}

	return &userData, nil
}

func (r userRepository) UserCreate(ctx context.Context, phoneNumber string) (*models.User, error) {
	user := &models.User{
		PhoneNumber: phoneNumber,
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Error("Error while creating user", zap.Error(err), zap.String("phone_number", phoneNumber))
		return nil, err
	}

	r.logger.Debug("User created successfully", zap.String("user_id", user.ID.String()), zap.String("phone_number", phoneNumber))
	return user, nil
}
