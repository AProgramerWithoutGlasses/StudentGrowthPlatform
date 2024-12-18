package timeConverter

import (
	"fmt"
	"time"
)

// GetEnrollmentYear 根据年级计算入学年份xxxx年9月1日
func GetEnrollmentYear(grade int) (time.Time, error) {
	// 获取当前年份
	currentYear, currentMonth, _ := time.Now().Date()

	if currentMonth < 9 {
		currentYear--
	}

	// 根据年级计算入学年份
	switch grade {
	case 1:
		return time.Date(currentYear, 9, 1, 0, 0, 0, 0, time.UTC), nil
	case 2:
		return time.Date(currentYear-1, 9, 1, 0, 0, 0, 0, time.UTC), nil
	case 3:
		return time.Date(currentYear-2, 9, 1, 0, 0, 0, 0, time.UTC), nil
	case 4:
		return time.Date(currentYear-3, 9, 1, 0, 0, 0, 0, time.UTC), nil
	default:
		return time.Time{}, fmt.Errorf("无效的年级: %s", grade)
	}
}

// GetUserGrade 根据入学年份计算年级
func GetUserGrade(plusTime time.Time) int {
	// 获取当前年份
	currentYear, currentMonth, _ := time.Now().Date()
	year := currentYear - plusTime.Year()

	if currentMonth < 8 {
		year--
	}

	switch year {
	case 0:
		return 1
	case 1:
		return 2
	case 2:
		return 3
	case 3:
		return 4
	default:
		return -1
	}

}
