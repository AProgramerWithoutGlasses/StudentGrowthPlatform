package gorm_model

import "gorm.io/gorm"

type NotificationConfig struct {
	gorm.Model
	NoticeTypeId int
	NoticeDesc   string // 消息描述
}
