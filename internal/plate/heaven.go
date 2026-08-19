package plate

// BuildHeavenPlate 排天盘三奇六仪
// 天盘干由九星携带：某宫所落九星的原始宫位之地盘干，即该宫的天盘干
// 中五宫天盘干取地盘乙（天禽所带）
func BuildHeavenPlate(earth map[int]string, stars map[int]string) map[int]string {
	heaven := make(map[int]string)
	for p := 1; p <= 9; p++ {
		star := stars[p]
		if star == "" {
			continue
		}
		// 找到该星原始的宫位
		originalPalace := -1
		for op, s := range palaceStar {
			if s == star {
				originalPalace = op
				break
			}
		}
		if originalPalace > 0 {
			heaven[p] = earth[originalPalace]
		}
	}
	// 中五宫天禽星所带地盘干
	if v, ok := earth[5]; ok {
		heaven[5] = v
	}
	return heaven
}
