package timeConverter

import (
	"fmt"
	"time"
)

// IntervalConversion 换算时间间隔 例如：X小时前 xxxx
func IntervalConversion(t time.Time) string {
	now := time.Now()
	delta := now.Sub(t)

	if delta >= 24*time.Hour {
		days := int(delta.Hours() / 24)
		return fmt.Sprintf("%d天前", days)
	} else if delta >= time.Hour {
		hours := int(delta.Hours())
		return fmt.Sprintf("%d小时前", hours)
	} else if delta >= time.Minute {
		minutes := int(delta.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	} else {
		return "刚刚"
	}
}
