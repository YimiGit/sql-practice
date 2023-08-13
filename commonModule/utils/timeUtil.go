package utils

import "time"

// TodayRemainNanosecond 获取今天剩余的纳秒数
func TodayRemainNanosecond() int64 {
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return endOfDay.Sub(now).Nanoseconds()
}

// TodayFormat 获取今天的日期格式
func TodayFormat(formatter string) string {
	switch formatter {
	case "yyyy":
		return time.Now().Format("2006")
	case "yyyyMM":
		return time.Now().Format("200601")
	case "yyyyMMdd":
		return time.Now().Format("20060102")
	case "yyyyMMddHH":
		return time.Now().Format("2006010215")
	case "yyyyMMddHHmm":
		return time.Now().Format("200601021504")
	case "yyyyMMddHHmmss":
		return time.Now().Format("20060102150405")
	}
	panic("formatter参数错误")
}
