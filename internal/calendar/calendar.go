package calendar

import (
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// SiZhu 返回给定北京时间的四柱（年月日时干支）
func SiZhu(t time.Time) (year, month, day, hour string) {
	solar := calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
	lunar := solar.GetLunar()
	return lunar.GetYearInGanZhi(),
		lunar.GetMonthInGanZhi(),
		lunar.GetDayInGanZhi(),
		lunar.GetTimeInGanZhi()
}

// CurrentJieQi 返回当前所处节气（最近过去的节气）名称
func CurrentJieQi(t time.Time) string {
	solar := calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
	lunar := solar.GetLunar()
	jq := lunar.GetPrevJieQi()
	if jq == nil {
		return ""
	}
	return jq.GetName()
}

// CurrentJie 返回当前所处节令（非气令）名称，用于定月干支
func CurrentJie(t time.Time) string {
	solar := calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
	lunar := solar.GetLunar()
	jq := lunar.GetPrevJie()
	if jq == nil {
		return ""
	}
	return jq.GetName()
}

// IsJie 判断给定节气名称是否为“节”（用于定月）
func IsJie(name string) bool {
	for _, j := range JieList {
		if j == name {
			return true
		}
	}
	return false
}

// LunarString 返回农历字符串描述
func LunarString(t time.Time) string {
	solar := calendar.NewSolar(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
	return solar.GetLunar().String()
}
