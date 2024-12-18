package gorm_model

import (
	"gorm.io/gorm"
)

type ArticleTag struct {
	gorm.Model
	ArticleID uint `gorm:"not null"`
	Article   Article
	TagID     uint `gorm:"not null"`
	Tag       Tag
}
