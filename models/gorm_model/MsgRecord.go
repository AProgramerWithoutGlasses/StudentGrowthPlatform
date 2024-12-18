package gorm_model

import (
	"gorm.io/gorm"
)

// 每条通知对应一个用户
// 添加管理员通知或系统通知需要广播到相应的或所有用户

type MsgRecord struct {
	gorm.Model
	Username string //消息发布者的username
	Content  string `gorm:"not null" json:"msg_content"`
	Type     int    `gorm:"not null" json:"msg_type"`
	Time     string `gorm:"-" json:"msg_time"`
	IsRead   bool   `gorm:"default:false" json:"is_read"`
	UserID   uint
	User     User
}
