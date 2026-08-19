package calendar

// 天干
var GanList = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

// 地支
var ZhiList = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// 六十甲子
var JiaZi60 []string

// 十二节气（节），用于定月干支
var JieList = []string{
	"立春", "惊蛰", "清明", "立夏", "芒种", "小暑",
	"立秋", "白露", "寒露", "立冬", "大雪", "小寒",
}

// 二十四节气
var JieQi24 = []string{
	"冬至", "小寒", "大寒",
	"立春", "雨水", "惊蛰",
	"春分", "清明", "谷雨",
	"立夏", "小满", "芒种",
	"夏至", "小暑", "大暑",
	"立秋", "处暑", "白露",
	"秋分", "寒露", "霜降",
	"立冬", "小雪", "大雪",
}

// 节气对应的太阳黄经
var JieQiLongitude = []float64{
	270, 285, 300,
	315, 330, 345,
	0, 15, 30,
	45, 60, 75,
	90, 105, 120,
	135, 150, 165,
	180, 195, 210,
	225, 240, 255,
}

// 九宫对应的地支（按坎一宫到离九宫）
var PalaceBranches = []string{"子", "未申", "卯", "辰巳", "", "戌亥", "酉", "丑寅", "午"}

// 九宫对应的八卦
var PalaceGua = []string{"坎", "坤", "震", "巽", "中", "乾", "兑", "艮", "离"}

// 九宫对应的五行
var PalaceElements = []string{"水", "土", "木", "木", "土", "金", "金", "土", "火"}

func init() {
	// 生成六十甲子
	for i := 0; i < 60; i++ {
		gan := GanList[i%10]
		zhi := ZhiList[i%12]
		JiaZi60 = append(JiaZi60, gan+zhi)
	}
}

// ganIndex 返回天干索引
func ganIndex(gan string) int {
	for i, g := range GanList {
		if g == gan {
			return i
		}
	}
	return -1
}

// zhiIndex 返回地支索引
func zhiIndex(zhi string) int {
	for i, z := range ZhiList {
		if z == zhi {
			return i
		}
	}
	return -1
}

// GanZhiIndex 返回干支在六十甲子中的索引（0=甲子）
func GanZhiIndex(gz string) int {
	for i, s := range JiaZi60 {
		if s == gz {
			return i
		}
	}
	return -1
}

// IndexToGanZhi 将索引转为干支
func IndexToGanZhi(idx int) string {
	idx = ((idx % 60) + 60) % 60
	return JiaZi60[idx]
}

// SplitGanZhi 拆分干支
func SplitGanZhi(gz string) (gan, zhi string) {
	runes := []rune(gz)
	if len(runes) != 2 {
		return "", ""
	}
	return string(runes[0]), string(runes[1])
}

// XunShou 返回指定干支所在旬的旬首（如甲戌日返回甲戌）
func XunShou(gz string) string {
	idx := GanZhiIndex(gz)
	if idx < 0 {
		return ""
	}
	xunIdx := (idx / 10) * 10
	return IndexToGanZhi(xunIdx)
}

// DunGan 返回六甲旬首所遁六仪（甲子戊、甲戌己、甲申庚、甲午辛、甲辰壬、甲寅癸）
func DunGan(xunShou string) string {
	mapping := map[string]string{
		"甲子": "戊",
		"甲戌": "己",
		"甲申": "庚",
		"甲午": "辛",
		"甲辰": "壬",
		"甲寅": "癸",
	}
	return mapping[xunShou]
}
