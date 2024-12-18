package gorm_model

import "gorm.io/gorm"

type UserCollectRecord struct {
	gorm.Model
	UserID    uint    `gorm:"not null"` //收藏属于用户
	User      User    //收藏属于用户
	ArticleID uint    `gorm:"not null"` //收藏属于文章
	Article   Article //收藏属于文章
	IsRead    bool    `gorm:"default:false"`
	PostTime  string  `gorm:"-" json:"post_time"`
}
