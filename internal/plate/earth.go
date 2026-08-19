package plate

import "qimen-tool/internal/calendar"

// 三奇六仪顺序
var yiQiList = []string{"戊", "己", "庚", "辛", "壬", "癸", "丁", "丙", "乙"}

// 洛书九宫圆周顺序（顺时针）
var luoShuClockwise = []int{1, 8, 3, 4, 9, 2, 7, 6, 1}

// 洛书九宫圆周顺序（逆时针）
var luoShuCounterClockwise = []int{1, 6, 7, 2, 9, 4, 3, 8, 1}

// nextPalace 按洛书九宫圆周顺序返回下一宫
// 阳遁顺（顺时针）：1->8->3->4->9->2->7->6->1
// 阴遁逆（逆时针）：1->6->7->2->9->4->3->8->1
// 中五宫寄坤二宫，因此从 5 宫出发下一宫为 2 宫
func nextPalace(palace int, yin bool) int {
	if palace == 5 {
		return 2
	}
	order := luoShuClockwise
	if yin {
		order = luoShuCounterClockwise
	}
	for i := 0; i < len(order)-1; i++ {
		if order[i] == palace {
			return order[i+1]
		}
	}
	return palace
}

// prevPalace 与 nextPalace 相反
func prevPalace(palace int, yin bool) int {
	return nextPalace(palace, !yin)
}

// nextPalaceNumeric 按宫号数字顺序返回下一宫，用于地盘排布
// 阳遁顺：1->2->3->4->5->6->7->8->9->1
// 阴遁逆：1->9->8->7->6->5->4->3->2->1
func nextPalaceNumeric(palace int, yin bool) int {
	if yin {
		p := palace - 1
		if p < 1 {
			p = 9
		}
		return p
	}
	p := palace + 1
	if p > 9 {
		p = 1
	}
	return p
}

// prevPalaceNumeric 与 nextPalaceNumeric 相反
func prevPalaceNumeric(palace int, yin bool) int {
	return nextPalaceNumeric(palace, !yin)
}

// stepsBetween 计算从 from 到 to 沿圆周转盘方向的步数
func stepsBetween(from, to int, yin bool) int {
	if from == to {
		return 0
	}
	steps := 0
	p := from
	for {
		p = nextPalace(p, yin)
		steps++
		if p == to {
			return steps
		}
		if steps > 9 {
			return -1
		}
	}
}

// BuildEarthPlate 排地盘三奇六仪
// 返回 map[宫号]天干
// 地盘按宫号数字顺序阳顺阴逆布排
func BuildEarthPlate(ju int, yin bool) map[int]string {
	plate := make(map[int]string)
	p := ju
	for _, stem := range yiQiList {
		plate[p] = stem
		p = nextPalaceNumeric(p, yin)
	}
	return plate
}

// palaceOfStem 查找某天干在地盘上的宫位
func palaceOfStem(earth map[int]string, stem string) int {
	for p, s := range earth {
		if s == stem {
			return p
		}
	}
	return -1
}

// palaceOfZhi 查找地支所在宫位
func palaceOfZhi(zhi string) int {
	for p, branch := range calendar.PalaceBranches {
		for _, b := range []rune(branch) {
			if string(b) == zhi {
				return p + 1
			}
		}
	}
	return -1
}
