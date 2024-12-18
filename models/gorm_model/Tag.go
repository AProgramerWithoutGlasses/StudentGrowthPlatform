package gorm_model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	TopicID uint
	Topic   Topic
	TagName string `json:"tag_name"`
}
