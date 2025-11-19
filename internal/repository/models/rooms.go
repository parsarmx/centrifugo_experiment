package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RoomName  string    `gorm:"type:varchar(255)" json:"room_name"`
	Channel   string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"channel"`
	Capacity  int       `gorm:"type:int" json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoomMember struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RoomID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_room_user" json:"room_id"`
	Room     Room      `gorm:"constraint:OnDelete:CASCADE;" json:"room"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_room_user" json:"user_id"`
	User     User      `gorm:"constraint:OnDelete:CASCADE;" json:"user"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

type Message struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RoomID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"room_id"`
	Room      Room       `gorm:"constraint:OnDelete:CASCADE;" json:"room"`
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	User      *User      `gorm:"constraint:OnDelete:SET NULL;" json:"user"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
