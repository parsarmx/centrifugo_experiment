package dto

type CreateRoomRequest struct {
	RoomName string `json:"room_name" validate:"required,min=3"`
	Capacity int    `json:"capacity" validate:"required,min=1"`
}

type SendMessageRequest struct {
	Message string `json:"message" validate:"required,min=1"`
	Channel string `json:"channel" validate:"required,min=3"`
}
