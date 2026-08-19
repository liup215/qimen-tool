package plate

import "qimen-tool/internal/calendar"

// 八门顺序（转盘法）
var doorOrder = []string{"休", "生", "伤", "杜", "景", "死", "惊", "开"}

// 宫位对应的原始八门
var palaceDoor = map[int]string{
	1: "休", // 坎
	8: "生", // 艮
	3: "伤", // 震
	4: "杜", // 巽
	9: "景", // 离
	2: "死", // 坤
	7: "惊", // 兑
	6: "开", // 乾
}

// doorIndex 返回门在 doorOrder 中的索引
func doorIndex(door string) int {
	for i, d := range doorOrder {
		if d == door {
			return i
		}
	}
	return -1
}

// BuildDoors 排列八门
// xunShouPalace: 旬首宫
// xunShou: 旬首干支（用于取旬首时支）
// hourZhi: 当前时支
// yin: 是否阴遁
// BuildDoors 排列八门，返回值使门落宫
// xunShouPalace: 旬首宫
// xunShou: 旬首干支（用于取旬首时支）
// hourZhi: 当前时支
// yin: 是否阴遁
func BuildDoors(xunShouPalace int, xunShou, hourZhi string, yin bool) (map[int]string, int) {
	doors := make(map[int]string)
	// 中五宫寄坤二宫
	if xunShouPalace == 5 {
		xunShouPalace = 2
	}
	zhiShiDoor := palaceDoor[xunShouPalace]
	idx := doorIndex(zhiShiDoor)
	if idx < 0 {
		return doors, -1
	}

	// 值使门落宫：从旬首时支到当前时支，每时辰走一步，按宫号数字顺序阳顺阴逆
	steps := zhiStepsBetween(xunShou, hourZhi)
	p := xunShouPalace
	for i := 0; i < steps; i++ {
		p = nextPalaceNumeric(p, yin)
	}
	zhiShiPalace := p

	// 排列八门：以值使门落宫为起点，按洛书圆周顺序阳顺阴逆
	p = zhiShiPalace
	for i := 0; i < 8; i++ {
		if p == 5 {
			p = nextPalace(p, yin)
		}
		doors[p] = doorOrder[idx]
		p = nextPalace(p, yin)
		if yin {
			idx--
			if idx < 0 {
				idx += 8
			}
		} else {
			idx++
			if idx >= 8 {
				idx -= 8
			}
		}
	}
	return doors, zhiShiPalace
}

// zhiStepsBetween 计算从旬首时支到当前时支经过的时辰数（步数）
// 阳遁宫号递增，阴遁宫号递减，步数均按地支顺数方向计算
func zhiStepsBetween(xunShou, hourZhi string) int {
	_, xunZhi := calendar.SplitGanZhi(xunShou)
	xunIdx := zhiOrderIndex(xunZhi)
	hourIdx := zhiOrderIndex(hourZhi)
	if xunIdx < 0 || hourIdx < 0 {
		return 0
	}
	return (hourIdx - xunIdx + 12) % 12
}

// zhiOrderIndex 返回地支在十二地支中的索引
func zhiOrderIndex(zhi string) int {
	for i, z := range calendar.ZhiList {
		if z == zhi {
			return i
		}
	}
	return -1
}

// ZhiShiDoor 返回值使门
func ZhiShiDoor(xunShouPalace int) string {
	return palaceDoor[xunShouPalace]
}
