package dto

type CreateRoomRequest struct {
	RoomName string `json:"room_name" validate:"required,min=3"`
	Capacity int    `json:"capacity" validate:"required,min=1"`
}
