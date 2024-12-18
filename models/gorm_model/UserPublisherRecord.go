package gorm_model

import "gorm.io/gorm"

// UserPublisherRecord 用于记录用户添加时的时间和添加人
type UserPublisherRecord struct {
	gorm.Model
	Username string `json:"username"` // 添加用户的人 的账号
	Users    []User // 添加用户的人 和 用户 为一对多关系
}
