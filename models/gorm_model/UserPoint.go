package gorm_model

import "gorm.io/gorm"

type UserPoint struct {
	gorm.Model
	UserID  uint
	User    User
	TopicID uint
	Topic   Topic
	Point   int `gorm:"default:0"`
}
