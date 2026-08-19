package tests

import (
	"fmt"
	"testing"

	"qimen-tool/internal/plate"
	"qimen-tool/internal/verify"
)

// goldenCase 记录一个被验证过的盘例关键字段，用于回归测试。
// plate.JieQi 为当前完整节气（CurrentJieQi），而定局所用之“节”由 DetermineBureau 内部另行计算。
type goldenCase struct {
	name          string
	year          int
	month         int
	day           int
	hour          int
	minute        int
	wantSiZhu     [4]string
	wantJieQi     string
	wantDun       string
	wantBureau    int
	wantXunShou   string
	wantZhiFu     string
	wantZhiFuPal  int
	wantZhiShi    string
	wantZhiShiPal int
	wantKongWang  string
	wantMaXing    string
}

var goldenCases = []goldenCase{
	{
		name:          "2024-08-19 立秋阴遁5局",
		year: 2024, month: 8, day: 19, hour: 14, minute: 30,
		wantSiZhu:     [4]string{"甲辰", "壬申", "乙卯", "癸未"},
		wantJieQi:     "立秋",
		wantDun:       "阴遁",
		wantBureau:    5,
		wantXunShou:   "甲戌",
		wantZhiFu:     "天辅",
		wantZhiFuPal:  9,
		wantZhiShi:    "杜",
		wantZhiShiPal: 4,
		wantKongWang:  "申酉",
		wantMaXing:    "巳",
	},
	{
		name:          "2017-05-04 谷雨阳遁8局",
		year: 2017, month: 5, day: 4, hour: 15, minute: 30,
		wantSiZhu:     [4]string{"丁酉", "甲辰", "辛卯", "丙申"},
		wantJieQi:     "谷雨",
		wantDun:       "阳遁",
		wantBureau:    8,
		wantXunShou:   "甲午",
		wantZhiFu:     "天芮",
		wantZhiFuPal:  6,
		wantZhiShi:    "死",
		wantZhiShiPal: 4,
		wantKongWang:  "辰巳",
		wantMaXing:    "巳",
	},
	{
		name:          "2016-07-31 大暑阴遁1局",
		year: 2016, month: 7, day: 31, hour: 18, minute: 0,
		wantSiZhu:     [4]string{"丙申", "乙未", "甲寅", "癸酉"},
		wantJieQi:     "大暑",
		wantDun:       "阴遁",
		wantBureau:    1,
		wantXunShou:   "甲子",
		wantZhiFu:     "天蓬",
		wantZhiFuPal:  2,
		wantZhiShi:    "休",
		wantZhiShiPal: 1,
		wantKongWang:  "戌亥",
		wantMaXing:    "申",
	},
	{
		name:          "2024-02-04 大寒阳遁3局",
		year: 2024, month: 2, day: 4, hour: 10, minute: 0,
		wantSiZhu:     [4]string{"癸卯", "丙寅", "戊戌", "丁巳"},
		wantJieQi:     "大寒",
		wantDun:       "阳遁",
		wantBureau:    3,
		wantXunShou:   "甲寅",
		wantZhiFu:     "天任",
		wantZhiFuPal:  9,
		wantZhiShi:    "生",
		wantZhiShiPal: 2,
		wantKongWang:  "子丑",
		wantMaXing:    "申",
	},
	{
		name:          "2023-12-22 冬至阳遁7局",
		year: 2023, month: 12, day: 22, hour: 12, minute: 0,
		wantSiZhu:     [4]string{"癸卯", "甲子", "甲寅", "庚午"},
		wantJieQi:     "冬至",
		wantDun:       "阳遁",
		wantBureau:    7,
		wantXunShou:   "甲子",
		wantZhiFu:     "天柱",
		wantZhiFuPal:  9,
		wantZhiShi:    "惊",
		wantZhiShiPal: 4,
		wantKongWang:  "戌亥",
		wantMaXing:    "申",
	},
}

func TestGoldenCases(t *testing.T) {
	for _, tc := range goldenCases {
		t.Run(tc.name, func(t *testing.T) {
			tm := beijingTime(tc.year, tc.month, tc.day, tc.hour, tc.minute, 0)
			p, err := plate.BuildPlate(tm)
			if err != nil {
				t.Fatalf("BuildPlate failed: %v", err)
			}

			if p.YearGanZhi != tc.wantSiZhu[0] || p.MonthGanZhi != tc.wantSiZhu[1] ||
				p.DayGanZhi != tc.wantSiZhu[2] || p.HourGanZhi != tc.wantSiZhu[3] {
				t.Errorf("四柱 = %s %s %s %s, want %v", p.YearGanZhi, p.MonthGanZhi, p.DayGanZhi, p.HourGanZhi, tc.wantSiZhu)
			}
			if p.JieQi != tc.wantJieQi {
				t.Errorf("节气 = %s, want %s", p.JieQi, tc.wantJieQi)
			}
			if p.Dun != tc.wantDun || p.Bureau != tc.wantBureau {
				t.Errorf("定局 = %s %d局, want %s %d局", p.Dun, p.Bureau, tc.wantDun, tc.wantBureau)
			}
			if p.XunShou != tc.wantXunShou {
				t.Errorf("旬首 = %s, want %s", p.XunShou, tc.wantXunShou)
			}
			if p.ZhiFuStar != tc.wantZhiFu || p.ZhiFuPalace != tc.wantZhiFuPal {
				t.Errorf("值符 = %s在%d宫, want %s在%d宫", p.ZhiFuStar, p.ZhiFuPalace, tc.wantZhiFu, tc.wantZhiFuPal)
			}
			if p.ZhiShiDoor != tc.wantZhiShi || p.ZhiShiPalace != tc.wantZhiShiPal {
				t.Errorf("值使 = %s在%d宫, want %s在%d宫", p.ZhiShiDoor, p.ZhiShiPalace, tc.wantZhiShi, tc.wantZhiShiPal)
			}
			if p.KongWang != tc.wantKongWang {
				t.Errorf("空亡 = %s, want %s", p.KongWang, tc.wantKongWang)
			}
			if p.MaXing != tc.wantMaXing {
				t.Errorf("马星 = %s, want %s", p.MaXing, tc.wantMaXing)
			}

			// 同时运行 verify 自检链
			results := verify.Verify(p)
			if !verify.AllPassed(results) {
				t.Errorf("verify 失败: %+v", results)
			}
		})
	}
}

// TestBoundaryTimes 验证子时、午时、节气交界等边界时刻不崩溃且 verify 通过。
func TestBoundaryTimes(t *testing.T) {
	cases := []struct {
		name string
		tm   [6]int // year month day hour minute second
	}{
		{"子时交界", [6]int{2024, 8, 19, 0, 0, 0}},
		{"午时正中", [6]int{2024, 8, 19, 12, 0, 0}},
		{"立春交节", [6]int{2024, 2, 4, 16, 27, 0}}, // 2024立春约16:27
		{"夏至交节", [6]int{2024, 6, 21, 4, 51, 0}}, // 2024夏至约04:51
		{"冬至交节", [6]int{2023, 12, 22, 11, 27, 0}}, // 2023冬至约11:27
		{"跨年夜", [6]int{2023, 12, 31, 23, 59, 59}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tm := beijingTime(tc.tm[0], tc.tm[1], tc.tm[2], tc.tm[3], tc.tm[4], tc.tm[5])
			p, err := plate.BuildPlate(tm)
			if err != nil {
				t.Fatalf("BuildPlate failed: %v", err)
			}
			results := verify.Verify(p)
			if !verify.AllPassed(results) {
				t.Errorf("verify 失败: %+v", results)
			}
			fmt.Printf("%s: %s %s%d局 旬首%s 值符%s在%d宫\n",
				tc.name, p.JieQi, p.Dun, p.Bureau, p.XunShou, p.ZhiFuStar, p.ZhiFuPalace)
		})
	}
}
