package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PhoneNumber string    `gorm:"type:varchar(255);check:phone_number ~ '^\\+[1-9][0-9]{7,14}$'" json:"phone_number"`
	FirstName   string    `gorm:"type:varchar(255)" json:"first_name"`
	LastName    string    `gorm:"type:varchar(255)" json:"last_name"`
	Age         int       `gorm:"type:int;check:age >= 0" json:"age"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
