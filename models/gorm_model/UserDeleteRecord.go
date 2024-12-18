package gorm_model

import (
	"gorm.io/gorm"
)

type UserDeleteRecord struct {
	gorm.Model
	Username       string `json:"username"`
	DeleteUsername string `json:"delete_username"`
}
