package timeutil

import (
	"time"
)

const (
	FormatYyyyMmDdHhIiNormal            = "2006-01-02 15:04"
	FormatYyyyMmDdHhIiSsNormal          = "2006-01-02 15:04:05"
	FormatYyyyMmDdHhIiSsNormalWithMilli = "2006-01-02 15:04:05.000"
	FormatYyyyMmDdNormal                = "2006-01-02"
	FormatYyyyMmDdNoSymbol              = "20060102"
	FormatYyyyMmDdHhIiSsNoSymbol        = "20060102150405"
	FormatYyyyMmDdHhIiNoSymbol          = "200601021504"
	FormatYyyy                          = "2006"
	FormatMm                            = "01"
	FormatDd                            = "02"
	FormatRFC1123                       = time.RFC1123
	FormatRFC3339Nano                   = time.RFC3339Nano
	Hour                                = time.Hour
	Minute                              = time.Minute
	Second                              = time.Second
)

type Timezone struct {
	Name    string         //时区名字
	Offset  int            //时区偏移
	TimeLoc *time.Location //golang 时区
}

var TimezoneUtc = time.UTC
var TimezoneShanghai, _ = time.LoadLocation("Asia/Shanghai")
var TimezoneJp, _ = time.LoadLocation("Asia/Tokyo")

var TimezoneEntityUTC = Timezone{
	Name:    "UTC",
	Offset:  0,
	TimeLoc: TimezoneUtc,
}

var TimezoneEntityShanghai = Timezone{
	Name:    "Asia/Shanghai",
	Offset:  8,
	TimeLoc: TimezoneShanghai,
}

var TimezoneEntityJP = Timezone{
	Name:    "Asia/Tokyo",
	Offset:  9,
	TimeLoc: TimezoneJp,
}

// ChangeDayFormat 将某种格式的时间字符串，转为另一种格式
func ChangeDayFormat(day string, from string, to string) string {
	return Str2Time(day, from, TimezoneUtc).Format(to)
}

// DayNoSymbol2Normal FormatYyyyMmDdNoSymbol 转换 FormatYyyyMmDdNormal
func DayNoSymbol2Normal(day string) string {
	return ChangeDayFormat(day, FormatYyyyMmDdNoSymbol, FormatYyyyMmDdNormal)
}

// DayNormal2NoSymbol FormatYyyyMmDdNormal 转换 FormatYyyyMmDdNoSymbol
func DayNormal2NoSymbol(day string) string {
	return ChangeDayFormat(day, FormatYyyyMmDdNormal, FormatYyyyMmDdNoSymbol)
}

// Str2Time 时间字符串按照格式解析出time.Time
func Str2Time(timeStr string, format string, tZone *time.Location) time.Time {
	theTime, _ := time.ParseInLocation(format, timeStr, tZone)
	return theTime
}

// Time2Str time.Time 按照格式返回时间字符串
func Time2Str(t time.Time, format string, tZone *time.Location) string {
	return t.In(tZone).Format(format)
}

// DayNoSymbol2Second FormatYyyyMmDdNoSymbol 转换 Second
func DayNoSymbol2Second(day string, tZone *time.Location) int64 {
	theTime, _ := time.ParseInLocation(FormatYyyyMmDdNoSymbol, day, tZone)
	return theTime.Unix()
}

// Second2DayNoSymbol Second 转换 FormatYyyyMmDdNoSymbol
func Second2DayNoSymbol(unix int64, tZone *time.Location) string {
	tm := time.Unix(unix, 0)
	return tm.In(tZone).Format(FormatYyyyMmDdNoSymbol)
}

// GetRangeDayNoSymbol  获取从某日到某日的所有天（包括起止点，格式为 FormatYyyyMmDdNoSymbol）
func GetRangeDayNoSymbol(from string, to string) []string {
	begin := DayNoSymbol2Second(from, TimezoneUtc)
	end := DayNoSymbol2Second(to, TimezoneUtc)
	var ret []string
	for i := begin; i <= end; i += 86400 {
		ret = append(ret, Second2DayNoSymbol(i, TimezoneUtc))
	}
	return ret
}

// GetRangeDayNoSymbolByTZone 设置时区-获取从某日到某日的所有天（包括起止点，格式为 FormatYyyyMmDdNoSymbol）
func GetRangeDayNoSymbolByTZone(from string, to string, tZone *time.Location) []string {
	begin := DayNoSymbol2Second(from, tZone)
	end := DayNoSymbol2Second(to, tZone)
	var ret []string
	for i := begin; i <= end; i += 86400 {
		ret = append(ret, Second2DayNoSymbol(i, tZone))
	}
	return ret
}

// DayNoSymbolAdd 字符串日期加法
func DayNoSymbolAdd(day string, add int64) string {
	return Second2DayNoSymbol(DayNoSymbol2Second(day, TimezoneUtc)+add*86400, TimezoneUtc)
}

// GetMultiDuration 获取复数时间(时长):例如获取5小时：GetMultiDuration(time.Hour,5)
func GetMultiDuration(duration time.Duration, multi int) time.Duration {
	return duration * time.Duration(multi)
}

// GetYesterdayNoSymbol 获取昨天日期
func GetYesterdayNoSymbol(tZone *time.Location) string {
	now := time.Now().AddDate(0, 0, -1)
	return now.In(tZone).Format(FormatYyyyMmDdNoSymbol)
}

// GetTodayNoSymbol 获取今天日期
func GetTodayNoSymbol(tZone *time.Location) string {
	now := time.Now()
	return now.In(tZone).Format(FormatYyyyMmDdNoSymbol)
}

// GetDayNoSymbolBeforeYesterday  获取前天日期
func GetDayNoSymbolBeforeYesterday(tZone *time.Location) string {
	now := time.Now().AddDate(0, 0, -2)
	return now.In(tZone).Format(FormatYyyyMmDdNoSymbol)
}

// GetHour 获取当前小时
func GetHour(tZone *time.Location) int {
	now := time.Now()
	return now.In(tZone).Hour()
}

// GetMinute 获取当前分钟
func GetMinute(tZone *time.Location) int {
	now := time.Now()
	return now.In(tZone).Minute()
}

// GetCurrentTime 获取当前时间
func GetCurrentTime(tZone *time.Location, format string) string {
	return time.Now().In(tZone).Format(format)
}

// GetCurrentUnixTimestamp 获取当前时间戳
func GetCurrentUnixTimestamp() int64 {
	return time.Now().Unix()
}

// GetTimePart 获取时间戳的详情
func GetTimePart(ts int64, tZone *time.Location) (int, time.Month, int, int, int, int) {
	tm := time.Unix(ts, 0).In(tZone)
	return tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second()
}

// GetMonthFirstDayNoSymbol 获取月份第一天
func GetMonthFirstDayNoSymbol(day string, tZone *time.Location) string {
	t := Str2Time(day, FormatYyyyMmDdNoSymbol, tZone)
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, tZone).Format(FormatYyyyMmDdNoSymbol)
}

// GetMonthLastDayNoSymbol 获取月份最后一天
func GetMonthLastDayNoSymbol(day string, tZone *time.Location) string {
	t := Str2Time(day, FormatYyyyMmDdNoSymbol, tZone)
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, tZone).
		AddDate(0, 1, 0).
		Add(-time.Nanosecond).
		Format(FormatYyyyMmDdNoSymbol)
}

// GetTimeNumMinute 获取当前时间所在num分钟的整数段时间 （num最小为1，最大为59）
// 例如10:08:09，num为10 得到10:00:00
// 例如10:59:01，num为10 得到10:55:00
func GetTimeNumMinute(ts int64, num int, tZone *time.Location) time.Time {
	if num <= 0 {
		num = 1
	}
	if num >= 60 {
		num = 59
	}
	year, month, day, hour, minute, _ := GetTimePart(ts, tZone)
	theMin := minute / num * num
	return time.Date(year, month, day, hour, theMin, 0, 0, tZone)
}

// GetTimeNumHour 获取当前时间所在小时的整数段时间（num最小为1，最大为23）
// 例如10:08:09，num=1，得到10:00:00
// 例如10:59:01，num=3，得到09:00:00
func GetTimeNumHour(ts int64, num int, tZone *time.Location) time.Time {
	if num <= 0 {
		num = 1
	}
	if num >= 24 {
		num = 23
	}
	year, month, day, hour, _, _ := GetTimePart(ts, tZone)
	thHour := hour / num * num
	return time.Date(year, month, day, thHour, 0, 0, 0, tZone)
}

// IsTodayNoSymbol 判断day是否是今天
func IsTodayNoSymbol(day string, tZone *time.Location) bool {
	return GetTodayNoSymbol(tZone) == day
}

// DayNoSymbolDiff 两天相差的天数
func DayNoSymbolDiff(dayBig string, daySmall string) int64 {
	return (DayNoSymbol2Second(dayBig, TimezoneShanghai) - DayNoSymbol2Second(daySmall, TimezoneShanghai)) / 86400
}

// GetMondayNoSymbolOfWeek 获取本周周一的日期
func GetMondayNoSymbolOfWeek(day string, tZone *time.Location) (dayStr string) {
	t := Str2Time(day, FormatYyyyMmDdNoSymbol, tZone)
	if t.Weekday() == time.Monday {
		//修改hour、min、sec = 0后格式化
		dayStr = t.Format(FormatYyyyMmDdNoSymbol)
	} else {
		offset := int(time.Monday - t.Weekday())
		if offset > 0 {
			offset = -6
		}
		dayStr = t.AddDate(0, 0, offset).Format(FormatYyyyMmDdNoSymbol)
	}
	return dayStr
}

// GetTuesdayNoSymbolOfWeek 获取本周周二的日期
func GetTuesdayNoSymbolOfWeek(day string, tZone *time.Location) (dayStr string) {
	monday := GetMondayNoSymbolOfWeek(day, tZone)
	t := Str2Time(monday, FormatYyyyMmDdNoSymbol, tZone)
	return t.AddDate(0, 0, 1).Format(FormatYyyyMmDdNoSymbol)
}

// CheckDateFormat 验证日期格式
func CheckDateFormat(date, layout string) bool {
	_, err := time.Parse(layout, date)
	if err != nil {
		return false
	}
	return true
}
