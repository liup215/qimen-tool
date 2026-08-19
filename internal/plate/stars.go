package plate

// 九星顺序（转盘法，按洛书数字顺序）
var starOrder = []string{"天蓬", "天任", "天冲", "天辅", "天英", "天芮", "天柱", "天心"}

// 宫位对应的原始九星（天禽寄天芮）
var palaceStar = map[int]string{
	1: "天蓬",
	8: "天任",
	3: "天冲",
	4: "天辅",
	9: "天英",
	2: "天芮", // 天禽寄此
	5: "天禽", // 中五宫
	7: "天柱",
	6: "天心",
}

// starIndex 返回星在 starOrder 中的索引
func starIndex(star string) int {
	for i, s := range starOrder {
		if s == star {
			return i
		}
	}
	return -1
}

// BuildStars 排列九星
// xunShouPalace: 旬首宫
// hourStemPalace: 时干宫（值符落宫）
// yin: 是否阴遁
func BuildStars(xunShouPalace, hourStemPalace int, yin bool) map[int]string {
	stars := make(map[int]string)
	// 中五宫寄坤二宫
	if xunShouPalace == 5 {
		xunShouPalace = 2
	}
	zhiFuStar := palaceStar[xunShouPalace]
	if zhiFuStar == "天禽" {
		zhiFuStar = "天芮"
	}
	idx := starIndex(zhiFuStar)
	if idx < 0 {
		return stars
	}

	p := hourStemPalace
	for i := 0; i < 8; i++ {
		if p == 5 {
			p = nextPalace(p, yin)
		}
		stars[p] = starOrder[idx]
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
	return stars
}

// ZhiFuStar 返回值符星
func ZhiFuStar(xunShouPalace int) string {
	star := palaceStar[xunShouPalace]
	if star == "天禽" {
		return "天芮"
	}
	return star
}
