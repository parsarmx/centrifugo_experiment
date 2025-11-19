package repository

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang_template/internal/database/postgres"
	"golang_template/internal/repository/models"

	rpc_service "golang_template/proto"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, roomName string, capacity int) (*models.Room, error)
	LeaveRoom(ctx context.Context, roomID, userID string) error
}

type roomRepository struct {
	db     *gorm.DB
	logger *zap.Logger
	grpc   rpc_service.CentrifugoApiClient
}

func NewRoomRepository(db postgres.Database, logger *zap.Logger, grpc rpc_service.CentrifugoApiClient) RoomRepository {
	return &roomRepository{
		db:     db.Gorm(),
		logger: logger,
		grpc:   grpc,
	}
}

func (r roomRepository) CreateRoom(ctx context.Context, roomName string, capacity int) (*models.Room, error) {
	channel := slugifyChannel(roomName)

	originalSlug := channel
	counter := 1

	for {
		var count int64
		r.db.WithContext(ctx).Model(&models.Room{}).Where("channel = ?", channel).Count(&count)
		if count == 0 {
			break
		}
		counter++
		channel = originalSlug + "-" + strconv.Itoa(counter)
	}

	// Step 3: create the room
	room := &models.Room{
		RoomName: roomName,
		Channel:  channel,
		Capacity: capacity,
	}

	if err := r.db.WithContext(ctx).Create(room).Error; err != nil {
		r.logger.Error("Error while creating room", zap.Error(err))
		return nil, err
	}

	// Run a gRPC‌ request to publish joining in channel

	go func() {
		fmt.Println(room.Channel)
		req := &rpc_service.PublishRequest{
			Channel:     room.Channel,
			Data:        []byte(fmt.Sprintf(`{"room_id":"%s","room_name":"%s","capacity":%d}`, room.ID, room.RoomName, room.Capacity)),
			SkipHistory: false,
			Tags: map[string]string{
				"source": "room_service",
			},
		}

		fmt.Println(r.grpc)
		resp, err := r.grpc.Publish(context.Background(), req)
		if err != nil {
			fmt.Println(err)
			// fmt.Println(resp.Error.Code, resp.Error.Message)
			r.logger.Error("Failed to publish room creation", zap.Error(err))
			return
		}

		// resp.Error.Code is uint
		if resp.Error != nil && resp.Error.Code != 0 {
			fmt.Println(resp.Error.Code, resp.Error.Message)
			r.logger.Error("Centrifugo returned an error",
				zap.Uint32("code", resp.Error.Code),
				zap.String("msg", resp.Error.Message),
			)
			return
		}
		r.logger.Debug("Room published to Centrifugo successfully",
			zap.String("room_id", room.ID.String()),
			zap.String("channel", room.Channel),
		)
	}()

	r.logger.Debug("Room created successfully",
		zap.String("room_id", room.ID.String()),
		zap.String("channel", room.Channel),
	)

	return room, nil
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
