package controllers

import (
	"errors"
	"fmt"
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
	SendMessage(ctx echo.Context) error
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

	if err := ctx.Bind(&req); err != nil {
		c.logger.Error("Failed to bind request", zap.Error(err))
		return response.SendError(ctx, "Invalid request body", http.StatusBadRequest, err.Error())
	}

	if err := c.requestValidator.ValidateStruct(req); err != nil {
		return response.SendError(ctx, "Invalid input", http.StatusBadRequest, err.Error())
	}

	room, err := c.service.CreateRoom(ctx.Request().Context(), req.RoomName, req.Capacity)
	if err != nil {
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

func (c *roomController) SendMessage(ctx echo.Context) error {
	var req dto.SendMessageRequest

	userId := fmt.Sprintf("%v", ctx.Get("X-User-ID"))
	if err := ctx.Bind(&req); err != nil {
		c.logger.Error("Failed to bind request", zap.Error(err))
		return response.SendError(ctx, "Invalid request body", http.StatusBadRequest, err.Error())
	}

	if err := c.requestValidator.ValidateStruct(req); err != nil {
		return response.SendError(ctx, "Invalid input", http.StatusBadRequest, err.Error())
	}

	// it needs at least an error
	c.service.SendMessage(ctx.Request().Context(), userId, req.Message, req.Channel)

	return ctx.JSON(http.StatusAccepted, map[string]interface{}{
		"ok": true,
	})
}
