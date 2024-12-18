package gorm_model

import "gorm.io/gorm"

type UserClass struct {
	gorm.Model
	Class string `json:"user_class"`
}
