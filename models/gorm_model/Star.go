package gorm_model

import "gorm.io/gorm"

type Star struct {
	gorm.Model
	Username string
	Name     string
	Type     int
	Session  int `gorm:"default:0"`
}
