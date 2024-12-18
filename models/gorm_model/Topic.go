package gorm_model

import "gorm.io/gorm"

type Topic struct {
	gorm.Model
	TopicName    string `json:"topic_name"`
	TopicContent string `json:"topic_content"`
	Tags         []Tag  `json:"tags"`
}
