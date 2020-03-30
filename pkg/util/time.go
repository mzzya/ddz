package util

import "time"

// TimeFormat 时间格式化
func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
