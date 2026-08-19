package verify

import (
	"fmt"

	"qimen-tool/internal/models"
)

// Result 单条验证结果
type Result struct {
	Category string `json:"category"` // format / logic / consistency
	Pass     bool   `json:"pass"`
	Message  string `json:"message"`
}

// Verify 对排好的盘面执行断言自检链
func Verify(p *models.Plate) []Result {
	var results []Result
	results = append(results, verifyFormat(p)...)
	results = append(results, verifyLogic(p)...)
	results = append(results, verifyConsistency(p)...)
	return results
}

// AllPassed 判断所有验证项是否通过
func AllPassed(results []Result) bool {
	for _, r := range results {
		if !r.Pass {
			return false
		}
	}
	return true
}

func verifyFormat(p *models.Plate) []Result {
	var rs []Result
	if p == nil {
		rs = append(rs, Result{"format", false, "盘面为 nil"})
		return rs
	}
	if p.RuleSetVersion == "" {
		rs = append(rs, Result{"format", false, "规则集版本未设置"})
	} else {
		rs = append(rs, Result{"format", true, fmt.Sprintf("规则集版本：%s", p.RuleSetVersion)})
	}
	if len(p.Palaces) != 9 {
		rs = append(rs, Result{"format", false, fmt.Sprintf("宫位数量应为 9，实际 %d", len(p.Palaces))})
	} else {
		rs = append(rs, Result{"format", true, "宫位数量为 9"})
	}
	for i, palace := range p.Palaces {
		if palace.Number != i+1 {
			rs = append(rs, Result{"format", false, fmt.Sprintf("宫位顺序错误：索引 %d 对应宫号 %d", i, palace.Number)})
		}
	}
	if len(rs) == 2 {
		rs = append(rs, Result{"format", true, "宫位顺序正确"})
	}
	return rs
}

func verifyLogic(p *models.Plate) []Result {
	var rs []Result
	if p == nil || len(p.Palaces) != 9 {
		return rs
	}

	// 八门：8 个不同门，中五宫为空
	doors := make(map[string]int)
	for _, palace := range p.Palaces {
		if palace.Door != "" {
			doors[palace.Door]++
		}
	}
	if len(doors) == 8 {
		rs = append(rs, Result{"logic", true, fmt.Sprintf("八门完整：%d 个不同门", len(doors))})
	} else {
		rs = append(rs, Result{"logic", false, fmt.Sprintf("八门不完整：%d 个不同门", len(doors))})
	}
	for door, count := range doors {
		if count > 1 {
			rs = append(rs, Result{"logic", false, fmt.Sprintf("门 %s 重复出现 %d 次", door, count)})
		}
	}

	// 九星：转盘法中五寄坤，天禽寄天芮，故通常只显示 8 个不同星名
	stars := make(map[string]int)
	centerEmpty := false
	for _, palace := range p.Palaces {
		if palace.Number == 5 && palace.Star == "" {
			centerEmpty = true
		}
		if palace.Star != "" {
			stars[palace.Star]++
		}
	}
	if p.CenterPalace == "中五宫寄坤二宫" {
		if len(stars) == 8 && centerEmpty {
			rs = append(rs, Result{"logic", true, fmt.Sprintf("九星符合转盘法：%d 个不同星名，天禽寄天芮（中五空）", len(stars))})
		} else {
			rs = append(rs, Result{"logic", false, fmt.Sprintf("九星异常：%d 个不同星名，中五宫应空", len(stars))})
		}
	} else if len(stars) == 9 {
		rs = append(rs, Result{"logic", true, fmt.Sprintf("九星完整：%d 个不同星", len(stars))})
	} else {
		rs = append(rs, Result{"logic", false, fmt.Sprintf("九星不完整：%d 个不同星", len(stars))})
	}

	// 八神：8 个不同神
	spirits := make(map[string]int)
	for _, palace := range p.Palaces {
		if palace.Spirit != "" {
			spirits[palace.Spirit]++
		}
	}
	if len(spirits) == 8 {
		rs = append(rs, Result{"logic", true, fmt.Sprintf("八神完整：%d 个不同神", len(spirits))})
	} else {
		rs = append(rs, Result{"logic", false, fmt.Sprintf("八神不完整：%d 个不同神", len(spirits))})
	}

	return rs
}

func verifyConsistency(p *models.Plate) []Result {
	var rs []Result
	if p == nil || len(p.Palaces) != 9 {
		return rs
	}

	// 值符星应等于值符落宫的天盘星；中五寄坤时按坤二宫校验
	zfCheckPalace := p.ZhiFuPalace
	if zfCheckPalace == 5 {
		zfCheckPalace = 2
	}
	zfPalace := p.Palaces[zfCheckPalace-1]
	if zfPalace.Star == p.ZhiFuStar {
		rs = append(rs, Result{"consistency", true, fmt.Sprintf("值符星 %s 与落宫 %d 星一致（中五寄坤）", p.ZhiFuStar, p.ZhiFuPalace)})
	} else {
		rs = append(rs, Result{"consistency", false, fmt.Sprintf("值符星 %s 与落宫 %d 星 %s 不一致", p.ZhiFuStar, p.ZhiFuPalace, zfPalace.Star)})
	}

	// 值使门应等于值使落宫的门；中五寄坤时按坤二宫校验
	zsCheckPalace := p.ZhiShiPalace
	if zsCheckPalace == 5 {
		zsCheckPalace = 2
	}
	zsPalace := p.Palaces[zsCheckPalace-1]
	if zsPalace.Door == p.ZhiShiDoor {
		rs = append(rs, Result{"consistency", true, fmt.Sprintf("值使门 %s 与落宫 %d 门一致（中五寄坤）", p.ZhiShiDoor, p.ZhiShiPalace)})
	} else {
		rs = append(rs, Result{"consistency", false, fmt.Sprintf("值使门 %s 与落宫 %d 门 %s 不一致", p.ZhiShiDoor, p.ZhiShiPalace, zsPalace.Door)})
	}

	// 反查索引一致性
	for _, palace := range p.Palaces {
		if palace.Door != "" {
			if idx, ok := p.DoorIndex[palace.Door]; !ok || idx != palace.Number {
				rs = append(rs, Result{"consistency", false, fmt.Sprintf("door_index[%s]=%d 与宫位 %d 不一致", palace.Door, idx, palace.Number)})
			}
		}
		if palace.Star != "" {
			if idx, ok := p.StarIndex[palace.Star]; !ok || idx != palace.Number {
				rs = append(rs, Result{"consistency", false, fmt.Sprintf("star_index[%s]=%d 与宫位 %d 不一致", palace.Star, idx, palace.Number)})
			}
		}
		if palace.Spirit != "" {
			if idx, ok := p.SpiritIndex[palace.Spirit]; !ok || idx != palace.Number {
				rs = append(rs, Result{"consistency", false, fmt.Sprintf("spirit_index[%s]=%d 与宫位 %d 不一致", palace.Spirit, idx, palace.Number)})
			}
		}
		if palace.HeavenStem != "" {
			if idx, ok := p.StemIndex[palace.HeavenStem]; !ok || idx != palace.Number {
				rs = append(rs, Result{"consistency", false, fmt.Sprintf("stem_index[%s]=%d 与宫位 %d 不一致", palace.HeavenStem, idx, palace.Number)})
			}
		}
	}
	if len(rs) == 2 {
		rs = append(rs, Result{"consistency", true, "反查索引与宫位一致"})
	}

	// 日干/时干落宫一致性
	if p.DayStemIndex > 0 && p.DayStemIndex <= 9 {
		dayStemPalace := p.Palaces[p.DayStemIndex-1]
		rs = append(rs, Result{"consistency", true, fmt.Sprintf("日干落宫 %d，天盘干为 %s", p.DayStemIndex, dayStemPalace.HeavenStem)})
	}
	if p.HourStemIndex > 0 && p.HourStemIndex <= 9 {
		hourStemPalace := p.Palaces[p.HourStemIndex-1]
		rs = append(rs, Result{"consistency", true, fmt.Sprintf("时干落宫 %d，天盘干为 %s", p.HourStemIndex, hourStemPalace.HeavenStem)})
	}

	return rs
}
