package utils

import (
	"time"
)

// 辅助函数：获取周一日期
func GetMonday(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return t.AddDate(0, 0, -weekday+1)
}
