package repository

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"golang_template/internal/database/postgres"
	"golang_template/internal/repository/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, roomName string, capacity int) (*models.Room, error)
	LeaveRoom(ctx context.Context, roomID, userID string) error
	FetchAllRooms(ctx context.Context) ([]*models.Room, error)
}

type roomRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRoomRepository(db postgres.Database, logger *zap.Logger) RoomRepository {
	return &roomRepository{
		db:     db.Gorm(),
		logger: logger,
	}
}

func (r roomRepository) CreateRoom(ctx context.Context, roomName string, capacity int) (*models.Room, error) {
	channel := slugifyChannel(roomName)

	// Check if slug already exists
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Room{}).
		Where("channel = ?", channel).
		Count(&count).Error; err != nil {

		r.logger.Error("DB error while checking room slug", zap.Error(err))
		return nil, err
	}

	if count > 0 {
		r.logger.Debug("Room name already exists",
			zap.String("room_name", roomName),
			zap.String("channel", channel),
		)
		return nil, errors.New("Channel already exists")
	}

	// Create the room
	room := &models.Room{
		RoomName: roomName,
		Channel:  channel,
		Capacity: capacity,
	}

	if err := r.db.WithContext(ctx).Create(room).Error; err != nil {
		r.logger.Error("Error while creating room",
			zap.Error(err),
			zap.String("channel", channel),
		)
		return nil, err
	}

	r.logger.Debug("Room created successfully",
		zap.String("room_id", room.ID.String()),
		zap.String("channel", room.Channel),
	)

	return room, nil
}

func (r roomRepository) FetchAllRooms(ctx context.Context) ([]*models.Room, error) {
	var rooms []*models.Room

	if err := r.db.WithContext(ctx).Find(&rooms).Error; err != nil {
		r.logger.Error("Error while fetching rooms", zap.Error(err))
		return nil, err
	}

	r.logger.Debug("Fetched all rooms", zap.Int("count", len(rooms)))

	return rooms, nil
}

func (r roomRepository) LeaveRoom(ctx context.Context, roomID, userID string) error {
	result := r.db.WithContext(ctx).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(&models.RoomMember{})

	if result.Error != nil {
		r.logger.Error("Error while removing user from room",
			zap.Error(result.Error),
			zap.String("room_id", roomID),
			zap.String("user_id", userID),
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Debug("No membership found for user in room",
			zap.String("room_id", roomID),
			zap.String("user_id", userID),
		)
		return gorm.ErrRecordNotFound
	}

	r.logger.Debug("User left the room successfully",
		zap.String("room_id", roomID),
		zap.String("user_id", userID),
	)

	return nil
}

// generate channel names without any whitespace
var slugRegex = regexp.MustCompile(`[^a-z0-9-]+`)

func slugifyChannel(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = slugRegex.ReplaceAllString(s, "")
	return s
}
