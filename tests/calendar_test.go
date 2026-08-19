package tests

import (
	"testing"
	"time"

	"qimen-tool/internal/calendar"
)

func beijingTime(y, m, d, h, min, s int) time.Time {
	loc := time.FixedZone("Beijing", 8*60*60)
	return time.Date(y, time.Month(m), d, h, min, s, 0, loc)
}

func TestSiZhu_2024_08_19(t *testing.T) {
	tm := beijingTime(2024, 8, 19, 14, 30, 0)
	year, month, day, hour := calendar.SiZhu(tm)

	if year != "甲辰" || month != "壬申" || day != "乙卯" || hour != "癸未" {
		t.Errorf("SiZhu = %s %s %s %s, want 甲辰 壬申 乙卯 癸未", year, month, day, hour)
	}
}

func TestCurrentJieQi_2024_08_19(t *testing.T) {
	tm := beijingTime(2024, 8, 19, 14, 30, 0)
	jq := calendar.CurrentJieQi(tm)
	if jq != "立秋" {
		t.Errorf("CurrentJieQi = %s, want 立秋", jq)
	}
}

func TestCurrentJieQi_2024_02_04(t *testing.T) {
	// 2024年立春为 2月4日 16:26:53，12:00 仍在大寒
	tm := beijingTime(2024, 2, 4, 17, 0, 0)
	jq := calendar.CurrentJieQi(tm)
	if jq != "立春" {
		t.Errorf("CurrentJieQi = %s, want 立春", jq)
	}
}

func TestCurrentJie_2024_08_19(t *testing.T) {
	tm := beijingTime(2024, 8, 19, 14, 30, 0)
	j := calendar.CurrentJie(tm)
	if j != "立秋" {
		t.Errorf("CurrentJie = %s, want 立秋", j)
	}
}
