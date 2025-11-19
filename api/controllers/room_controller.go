package controllers

import (
	"errors"
	"golang_template/api/dto"
	"golang_template/api/response"
	"golang_template/internal/pkg"
	"golang_template/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RoomController interface {
	CreateRoom(ctx echo.Context) error
}

type roomController struct {
	service          service.RoomService
	logger           *zap.Logger
	requestValidator *pkg.Validator
}

func NewRoomController(service service.RoomService, logger *zap.Logger) RoomController {
	return &roomController{
		service:          service,
		logger:           logger,
		requestValidator: pkg.NewValidator(),
	}
}

func (c *roomController) CreateRoom(ctx echo.Context) error {
	var req dto.CreateRoomRequest

	// Parse JSON
	if err := ctx.Bind(&req); err != nil {
		c.logger.Error("Failed to bind request", zap.Error(err))
		return response.SendError(ctx, "Invalid request body", http.StatusBadRequest, err.Error())
	}

	// Validate
	if err := c.requestValidator.ValidateStruct(req); err != nil {
		return response.SendError(ctx, "Invalid input", http.StatusBadRequest, err.Error())
	}

	// Service call
	room, err := c.service.CreateRoom(ctx.Request().Context(), req.RoomName, req.Capacity)
	if err != nil {
		// Custom errors example: room already exists
		if errors.Is(err, service.ErrRoomNameExists) {
			return response.SendError(ctx, "Room already exists", http.StatusConflict, "duplicate_room")
		}

		c.logger.Error("Failed to create room", zap.Error(err))
		return response.SendError(ctx, "Internal error", http.StatusInternalServerError, "internal_error")
	}

	// Success
	return ctx.JSON(http.StatusCreated, map[string]interface{}{
		"ok":       true,
		"room_id":  room.ID,
		"channel":  room.Channel,
		"capacity": room.Capacity,
	})
}
