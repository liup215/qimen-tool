package plate

import (
	"fmt"
	"time"

	"qimen-tool/internal/calendar"
)

// Bureau 表示定局结果
type Bureau struct {
	JieQi string // 当前节气
	Dun   string // 阳遁 / 阴遁
	Yuan  string // 上元 / 中元 / 下元
	Ju    int    // 局数 1-9
}

// 拆补法定局表
// 键：节气名，值：[上元局, 中元局, 下元局]
var yangBureau = map[string][3]int{
	"冬至": {1, 7, 4},
	"惊蛰": {1, 7, 4},
	"小寒": {2, 8, 5},
	"大寒": {3, 9, 6},
	"春分": {3, 9, 6},
	"雨水": {9, 6, 3},
	"清明": {4, 1, 7},
	"立夏": {4, 1, 7},
	"立春": {8, 5, 2},
	"谷雨": {5, 2, 8},
	"小满": {5, 2, 8},
	"芒种": {6, 3, 9},
}

var yinBureau = map[string][3]int{
	"夏至": {9, 3, 6},
	"白露": {9, 3, 6},
	"小暑": {8, 2, 5},
	"大暑": {7, 1, 4},
	"秋分": {7, 1, 4},
	"立秋": {2, 5, 8},
	"寒露": {6, 9, 3},
	"立冬": {6, 9, 3},
	"处暑": {1, 4, 7},
	"霜降": {5, 8, 2},
	"小雪": {5, 8, 2},
	"大雪": {4, 7, 1},
}

// DetermineBureau 根据时间确定阴阳遁、元、局数（拆补法）
// 使用二十四节气（JieQi）定局，定元使用日柱符头（甲、己）
func DetermineBureau(t time.Time) (*Bureau, error) {
	jie := calendar.CurrentJieQi(t)
	if jie == "" {
		return nil, fmt.Errorf("无法确定当前节气")
	}

	_, _, dayGanZhi, _ := calendar.SiZhu(t)
	fuTou := FuTou(dayGanZhi)
	yuan := YuanFromFuTou(fuTou)

	var ju int
	var dun string
	if arr, ok := yangBureau[jie]; ok {
		dun = "阳遁"
		switch yuan {
		case "上元":
			ju = arr[0]
		case "中元":
			ju = arr[1]
		case "下元":
			ju = arr[2]
		}
	} else if arr, ok := yinBureau[jie]; ok {
		dun = "阴遁"
		switch yuan {
		case "上元":
			ju = arr[0]
		case "中元":
			ju = arr[1]
		case "下元":
			ju = arr[2]
		}
	} else {
		return nil, fmt.Errorf("未知节令: %s", jie)
	}

	return &Bureau{
		JieQi: jie,
		Dun:   dun,
		Yuan:  yuan,
		Ju:    ju,
	}, nil
}

// FuTou 根据日干支计算拆补法符头（往前推到最近的甲日或己日）
func FuTou(dayGanZhi string) string {
	idx := calendar.GanZhiIndex(dayGanZhi)
	if idx < 0 {
		return ""
	}
	for i := 0; i < 60; i++ {
		gz := calendar.IndexToGanZhi(idx - i)
		gan, _ := calendar.SplitGanZhi(gz)
		if gan == "甲" || gan == "己" {
			return gz
		}
	}
	return ""
}

// YuanFromFuTou 根据符头地支确定元
func YuanFromFuTou(fuTou string) string {
	_, zhi := calendar.SplitGanZhi(fuTou)
	switch zhi {
	case "子", "午", "卯", "酉":
		return "上元"
	case "寅", "申", "巳", "亥":
		return "中元"
	case "辰", "戌", "丑", "未":
		return "下元"
	}
	return ""
}

// XunShou 根据干支计算旬首
func XunShou(ganZhi string) string {
	idx := calendar.GanZhiIndex(ganZhi)
	if idx < 0 {
		return ""
	}
	xunShouIdx := idx - (idx % 10)
	return calendar.IndexToGanZhi(xunShouIdx)
}

// KongWang 根据旬首计算空亡
func KongWang(xunShou string) string {
	switch xunShou {
	case "甲子":
		return "戌亥"
	case "甲戌":
		return "申酉"
	case "甲申":
		return "午未"
	case "甲午":
		return "辰巳"
	case "甲辰":
		return "寅卯"
	case "甲寅":
		return "子丑"
	}
	return ""
}

// MaXing 根据地支三合局计算驿马。
// 口诀：申子辰马在寅，寅午戌马在申，巳酉丑马在亥，亥卯未马在巳。
// 时家奇门主流以时支取马星，亦可以日支取马星作辅助参考。
func MaXing(zhi string) string {
	switch zhi {
	case "申", "子", "辰":
		return "寅"
	case "寅", "午", "戌":
		return "申"
	case "巳", "酉", "丑":
		return "亥"
	case "亥", "卯", "未":
		return "巳"
	}
	return ""
}
