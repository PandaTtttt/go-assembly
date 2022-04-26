package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	fmt.Println(GetTimeNumMinute(1650945718, 1, TimezoneShanghai))
	fmt.Println(GetTimeNumHour(1650945718, 5, TimezoneShanghai))
	fmt.Println(ChangeDayFormat("20220106", FormatYyyyMmDdNoSymbol, FormatYyyyMmDdHhIiNormal))
	fmt.Println(DayNoSymbol2Normal("20220106"))
	fmt.Println(DayNormal2NoSymbol("2022-01-06"))
	fmt.Println(Str2Time("2022-01-06", FormatYyyyMmDdNormal, TimezoneJp))
	fmt.Println(Time2Str(time.Now(), FormatYyyyMmDdNormal, TimezoneJp))
	fmt.Println(DayNoSymbol2Second("20220106", TimezoneShanghai))
	fmt.Println(Second2DayNoSymbol(time.Now().Unix(), TimezoneShanghai))
	fmt.Println(GetRangeDayNoSymbol("20220106", "20220109"))
	fmt.Println(GetRangeDayNoSymbolByTZone("20220106", "20220109", TimezoneUtc))
	fmt.Println(DayNoSymbolAdd("20220106", 11))
	fmt.Println(GetMultiDuration(time.Minute, 11))
	fmt.Println(GetYesterdayNoSymbol(TimezoneShanghai))
	fmt.Println(GetTodayNoSymbol(TimezoneShanghai))
	fmt.Println(GetDayNoSymbolBeforeYesterday(TimezoneShanghai))
	fmt.Println(GetHour(TimezoneShanghai), GetMinute(TimezoneShanghai), GetCurrentTime(TimezoneShanghai, FormatYyyyMmDdHhIiNormal))
	fmt.Println(GetCurrentUnixTimestamp())
	fmt.Println(IsTodayNoSymbol("20220426", TimezoneJp))
	fmt.Println(DayNoSymbolDiff("20220424", "20220425"))
	fmt.Println(GetMondayNoSymbolOfWeek("20220426", TimezoneShanghai))
	fmt.Println(GetTuesdayNoSymbolOfWeek("20220426", TimezoneShanghai))
	fmt.Println(CheckDateFormat("20220426", FormatYyyyMmDdNoSymbol))
}
