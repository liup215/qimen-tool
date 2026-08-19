package tests

import (
	"fmt"
	"testing"

	"qimen-tool/internal/plate"
)

func TestBuildPlate_2024_08_19(t *testing.T) {
	tm := beijingTime(2024, 8, 19, 14, 30, 0)
	p, err := plate.BuildPlate(tm)
	if err != nil {
		t.Fatalf("BuildPlate failed: %v", err)
	}

	if p.YearGanZhi != "甲辰" || p.MonthGanZhi != "壬申" || p.DayGanZhi != "乙卯" || p.HourGanZhi != "癸未" {
		t.Errorf("四柱 = %s %s %s %s", p.YearGanZhi, p.MonthGanZhi, p.DayGanZhi, p.HourGanZhi)
	}
	if p.JieQi != "立秋" {
		t.Errorf("节气 = %s, want 立秋", p.JieQi)
	}
	if p.Dun != "阴遁" || p.Bureau != 5 {
		t.Errorf("定局 = %s %d 局, want 阴遁 5 局", p.Dun, p.Bureau)
	}
	if p.XunShou != "甲戌" {
		t.Errorf("旬首 = %s, want 甲戌（时柱癸未之旬首）", p.XunShou)
	}
	if p.ZhiFuStar != "天辅" || p.ZhiFuPalace != 9 {
		t.Errorf("值符 = %s 在 %d 宫, want 天辅在 9 宫", p.ZhiFuStar, p.ZhiFuPalace)
	}
	if p.ZhiShiDoor != "杜" || p.ZhiShiPalace != 4 {
		t.Errorf("值使 = %s 在 %d 宫, want 杜门在 4 宫", p.ZhiShiDoor, p.ZhiShiPalace)
	}
	if p.KongWang != "申酉" {
		t.Errorf("空亡 = %s, want 申酉", p.KongWang)
	}
	if p.MaXing != "巳" {
		t.Errorf("马星 = %s, want 巳", p.MaXing)
	}

	// 验证各宫关键元素
	expected := map[int]struct {
		heaven string
		earth  string
		door   string
		star   string
		spirit string
		kong   bool
		ma     bool
		dayMa  bool
	}{
		1: {heaven: "乙", earth: "壬", door: "休", star: "天心", spirit: "白虎", kong: false, ma: false, dayMa: false},
		2: {heaven: "癸", earth: "辛", door: "死", star: "天英", spirit: "九天", kong: true, ma: false, dayMa: false},
		3: {heaven: "丁", earth: "庚", door: "伤", star: "天任", spirit: "太阴", kong: false, ma: false, dayMa: false},
		4: {heaven: "庚", earth: "己", door: "杜", star: "天冲", spirit: "螣蛇", kong: false, ma: true, dayMa: true},
		6: {heaven: "丙", earth: "乙", door: "开", star: "天柱", spirit: "玄武", kong: false, ma: false, dayMa: false},
		7: {heaven: "辛", earth: "丙", door: "惊", star: "天芮", spirit: "九地", kong: true, ma: false, dayMa: false},
		8: {heaven: "壬", earth: "丁", door: "生", star: "天蓬", spirit: "六合", kong: false, ma: false, dayMa: false},
		9: {heaven: "己", earth: "癸", door: "景", star: "天辅", spirit: "值符", kong: false, ma: false, dayMa: false},
	}

	for i := 0; i < 9; i++ {
		palace := p.Palaces[i]
		if palace.Number == 5 {
			continue
		}
		exp := expected[palace.Number]
		if palace.HeavenStem != exp.heaven || palace.EarthStem != exp.earth ||
			palace.Door != exp.door || palace.Star != exp.star || palace.Spirit != exp.spirit ||
			palace.IsKongWang != exp.kong || palace.HasMaXing != exp.ma || palace.HasDayMaXing != exp.dayMa {
			t.Errorf("%d宫 期望 heaven=%s earth=%s door=%s star=%s spirit=%s kong=%v ma=%v dayMa=%v, 实际 heaven=%s earth=%s door=%s star=%s spirit=%s kong=%v ma=%v dayMa=%v",
				palace.Number, exp.heaven, exp.earth, exp.door, exp.star, exp.spirit, exp.kong, exp.ma, exp.dayMa,
				palace.HeavenStem, palace.EarthStem, palace.Door, palace.Star, palace.Spirit, palace.IsKongWang, palace.HasMaXing, palace.HasDayMaXing)
		}
	}

	// 打印盘面
	fmt.Printf("Plate: %+v\n", p)
	for _, palace := range p.Palaces {
		fmt.Printf("%d宫 %s: 天%s 地%s %s %s %s 空亡=%v 马星=%v\n",
			palace.Number, palace.Gua,
			palace.HeavenStem, palace.EarthStem,
			palace.Door, palace.Star, palace.Spirit,
			palace.IsKongWang, palace.HasMaXing)
	}
}
