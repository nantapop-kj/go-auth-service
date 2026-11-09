package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/enums"
)

type UsersLoginHistory struct {
	ID        uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    *uuid.UUID        `gorm:"type:uuid;index"`
	Email     string            `gorm:"type:varchar(255);index"`
	IPAddress string            `gorm:"type:varchar(50)"`
	UserAgent string            `gorm:"type:varchar(255)"`
	Action    enums.LoginAction `gorm:"type:varchar(20)"`
	CreatedAt time.Time
}

func (UsersLoginHistory) TableName() string {
	return "users_login_history"
}
