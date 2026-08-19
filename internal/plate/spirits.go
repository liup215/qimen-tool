package plate

// 八神顺序（转盘法）
var spiritOrder = []string{"值符", "螣蛇", "太阴", "六合", "白虎", "玄武", "九地", "九天"}

// BuildSpirits 排列八神
// zhiFuPalace: 值符落宫（即时干宫）
// yin: 是否阴遁
func BuildSpirits(zhiFuPalace int, yin bool) map[int]string {
	spirits := make(map[int]string)
	currentPalace := zhiFuPalace
	for i := 0; i < 8; i++ {
		if currentPalace == 5 {
			currentPalace = nextPalace(currentPalace, yin)
		}
		spirits[currentPalace] = spiritOrder[i]
		currentPalace = nextPalace(currentPalace, yin)
	}
	return spirits
}
