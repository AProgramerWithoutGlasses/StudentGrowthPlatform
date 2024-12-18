package NotificationPush

import (
	"studentGrow/dao/mysql"
	model "studentGrow/models/gorm_model"
	"studentGrow/utils/timeConverter"
	"time"
)

// BuildSystemNotification 构建系统消息
func BuildSystemNotification(username, content string) (*model.SysNotification, error) {
	userId, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		return nil, err
	}

	user, err := mysql.QueryUserByUserId(userId)
	if err != nil {
		return nil, err
	}

	sysNotification := model.SysNotification{
		OwnUserId:  uint(userId),
		NoticeType: 3,
		Content:    content,
		IsRead:     false,
		Time:       timeConverter.IntervalConversion(time.Now()),
		OwnUser:    *user,
	}

	return &sysNotification, err
}

// BuildManagerNotification 构建管理员消息
func BuildManagerNotification(username, content string) (*model.SysNotification, error) {
	userId, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		return nil, err
	}

	user, err := mysql.QueryUserByUserId(userId)
	if err != nil {
		return nil, err
	}

	sysNotification := model.SysNotification{
		OwnUserId:  uint(userId),
		NoticeType: 4,
		Content:    content,
		IsRead:     false,
		Time:       timeConverter.IntervalConversion(time.Now()),
		OwnUser:    *user,
	}

	return &sysNotification, nil
}
