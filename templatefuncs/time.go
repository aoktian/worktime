package templatefuncs

import (
	"fmt"
	"time"
)

func formatDate(unix int64) string {
	t := time.Unix(unix, 0)
	return t.Format("2006-01-02 15:04:05")
}

func unixToDay(unix int64) string {
	t := time.Unix(unix, 0)
	return t.Format("2006-01-02")
}

func secondsToTime(unix int64) string {
	seconds := time.Now().Unix() - unix
	years := seconds / 31536000
	seconds -= years * 31536000
	days := seconds / 86400

	if years > 0 {
		return fmt.Sprintf("历时 %d年 %d天", years, days)
	}
	if days > 0 {
		return fmt.Sprintf("历时 %d天", days)
	}

	return ""
}
