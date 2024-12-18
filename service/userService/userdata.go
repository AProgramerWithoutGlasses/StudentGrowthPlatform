package userService

import (
	"studentGrow/dao/mysql"
	"time"
)

// IntervalInDays 获取时间间隔的天数
func IntervalInDays(t time.Time) int {
	now := time.Now()
	delta := now.Sub(t)

	// 如果时间差大于或等于24小时，计算天数
	if delta >= 24*time.Hour {
		return int(delta.Hours() / 24)
	}
	// 如果时间差小于一天，返回0
	return 0
}

// UpdateStatus 是否可以解禁
func UpdateStatus(dataStr, username string) (bool, error) {
	date, err := time.Parse("2006-01-02", dataStr)
	if err != nil {
		return false, err
	}
	BanEndTime, err := mysql.SelEndTime(username)
	// 比较两个 time.Time 对象
	if date.Before(BanEndTime) {
		return false, nil
	} else if date.After(BanEndTime) {
		return true, nil
	} else {
		return true, nil
	}
}
