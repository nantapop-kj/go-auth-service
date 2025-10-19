package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email          string    `gorm:"type:varchar(100);uniqueIndex"`
	Password       string    `gorm:"type:varchar(100)" json:"-"`
	FirstName      string    `gorm:"type:varchar(100)" json:"first_name"`
	LastName       string    `gorm:"type:varchar(100)" json:"last_name"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	IsLoggedIn     bool      `gorm:"default:false" json:"is_logged_in"`
	CurrentSession string    `gorm:"type:varchar(100)" json:"current_session"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Users) TableName() string {
	return "users"
}
