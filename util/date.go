package util

import (
	"time"
)

// StrToIntMonth 字符串月份转整数月份
func TimeToIntMonth(t time.Time) int {
	return int(t.Month())
}

// GetTodayYMD 得到以sep为分隔符的年、月、日字符串(今天)
func GetTodayYMD(sep string) string {
	return time.Now().Format("2006" + sep + "01" + sep + "02")
}

// GetTodayYM 得到以sep为分隔符的年、月字符串(今天所属于的月份)
func GetTodayYM(sep string) string {
	return time.Now().Format("2006" + sep + "01")
}

// GetYesterdayYMD 得到以sep为分隔符的年、月、日字符串(昨天)
func GetYesterdayYMD(sep string) string {
	return time.Now().Add(-time.Hour * 24).Format("2006" + sep + "01" + sep + "02")
}

// GetTomorrowYMD 得到以sep为分隔符的年、月、日字符串(明天)
func GetTomorrowYMD(sep string) string {
	return time.Now().Add(time.Hour * 24).Format("2006" + sep + "01" + sep + "02")
}

// GetTodayTime 返回今天零点的time
func GetTodayTime() time.Time {
	now := time.Now()
	// now.Year(), now.Month(), now.Day() 是以本地时区为参照的年、月、日
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

// GetYesterdayTime 返回昨天零点的time
func GetYesterdayTime() time.Time {
	ysd := time.Now().Add(-time.Hour * 24)
	return time.Date(ysd.Year(), ysd.Month(), ysd.Day(), 0, 0, 0, 0, time.Local)
}
