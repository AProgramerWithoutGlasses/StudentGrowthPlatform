package gorm_model

import "gorm.io/gorm"

type UserReportArticleRecord struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	User      User
	ArticleID uint `gorm:"not null"`
	Article   Article
	Msg       string
	IsRead    bool `gorm:"default:false"`
}
