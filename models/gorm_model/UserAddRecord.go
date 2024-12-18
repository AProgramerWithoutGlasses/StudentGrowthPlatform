package gorm_model

import (
	"gorm.io/gorm"
)

type UserAddRecord struct {
	gorm.Model
	Username    string `json:"username"`
	AddUsername string `json:"add_username"`
}
