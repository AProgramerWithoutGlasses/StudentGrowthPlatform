package gorm_model

import (
	"gorm.io/gorm"
)

type UserReadRecord struct {
	gorm.Model
	UserID    uint `gorm:"not null"` //属于
	User      User
	ArticleID uint `gorm:"not null"` //属于
	Article   Article
}
