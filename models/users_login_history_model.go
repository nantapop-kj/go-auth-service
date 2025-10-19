package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/enums"
)

type UsersLoginHistory struct {
	ID        uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID         `gorm:"type:uuid;index" json:"user_id"`
	IPAddress string            `gorm:"type:varchar(50)" json:"ip_address"`
	UserAgent string            `gorm:"type:varchar(255)" json:"user_agent"`
	Action    enums.LoginAction `gorm:"type:varchar(20)" json:"action"`

	CreatedAt time.Time
}

func (UsersLoginHistory) TableName() string {
	return "users_login_history"
}
