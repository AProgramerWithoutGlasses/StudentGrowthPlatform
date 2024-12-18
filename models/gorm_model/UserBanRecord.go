package gorm_model

import (
	"gorm.io/gorm"
)

type UserBanRecord struct {
	gorm.Model
	UserID  int `json:"user_id"`
	BanId   int `json:"ban_id"`
	BanTime int `json:"ban_time"`
}
