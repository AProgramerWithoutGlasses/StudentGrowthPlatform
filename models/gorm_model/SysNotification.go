package gorm_model

import "gorm.io/gorm"

type SysNotification struct {
	gorm.Model
	OwnUserId  uint
	TarUserId  uint
	NoticeType int
	Content    string
	OwnUser    User   `gorm:"foreignKey:OwnUserId"` // 预加载发送者用户
	TarUser    User   `gorm:"foreignKey:TarUserId"` // 预加载接收者用户
	Status     bool   `gorm:"default:false"`
	Time       string `gorm:"-"`
	IsRead     bool   `gorm:"-"`
}
