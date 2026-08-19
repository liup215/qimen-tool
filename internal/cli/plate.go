package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"qimen-tool/internal/models"
	"qimen-tool/internal/plate"
)

// parseTime 解析用户输入的时间字符串
// 支持格式：2006-01-02 15:04, 2006-01-02 15:04:05, now
func parseTime(input string) (time.Time, error) {
	loc := time.FixedZone("Beijing", 8*60*60)
	now := time.Now().In(loc)

	if strings.TrimSpace(input) == "" || strings.ToLower(input) == "now" {
		return now, nil
	}

	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, f := range formats {
		if t, err := time.ParseInLocation(f, input, loc); err == nil {
			return t, nil
		}
	}

	return now, fmt.Errorf("无法解析时间: %s，支持格式如 2024-08-19 14:30 或 now", input)
}

// PlateToMarkdown 将盘面转为 Markdown 表格
func PlateToMarkdown(p *models.Plate) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# 奇门遁甲盘面\n\n"))
	sb.WriteString(fmt.Sprintf("- 阳历：%s\n", p.SolarTime))
	sb.WriteString(fmt.Sprintf("- 农历：%s\n", p.LunarTime))
	sb.WriteString(fmt.Sprintf("- 四柱：%s 年 %s 月 %s 日 %s 时\n", p.YearGanZhi, p.MonthGanZhi, p.DayGanZhi, p.HourGanZhi))
	sb.WriteString(fmt.Sprintf("- 节气：%s\n", p.JieQi))
	sb.WriteString(fmt.Sprintf("- 定局：%s %d 局\n", p.Dun, p.Bureau))
	sb.WriteString(fmt.Sprintf("- 旬首：%s，空亡：%s，马星：%s（时支马），日支马：%s\n", p.XunShou, p.KongWang, p.MaXing, p.DayMaXing))
	sb.WriteString(fmt.Sprintf("- 值符：%s 在 %d 宫，值使：%s 在 %d 宫\n", p.ZhiFuStar, p.ZhiFuPalace, p.ZhiShiDoor, p.ZhiShiPalace))
	sb.WriteString(fmt.Sprintf("- 日干 %s 在 %d 宫，时干 %s 在 %d 宫\n", dayStemOf(p.DayGanZhi), p.DayStemIndex, hourStemOf(p.HourGanZhi), p.HourStemIndex))
	sb.WriteString(fmt.Sprintf("- 规则集：%s，%s\n\n", p.RuleSetVersion, p.CenterPalace))

	// 九宫格 Markdown 表格（按洛书位置）
	sb.WriteString("| 巽四 | 离九 | 坤二 |\n")
	sb.WriteString("|------|------|------|\n")
	sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", palaceMarkdown(p, 4), palaceMarkdown(p, 9), palaceMarkdown(p, 2)))
	sb.WriteString("| 震三 | 中五 | 兑七 |\n")
	sb.WriteString("|------|------|------|\n")
	sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", palaceMarkdown(p, 3), palaceMarkdown(p, 5), palaceMarkdown(p, 7)))
	sb.WriteString("| 艮八 | 坎一 | 乾六 |\n")
	sb.WriteString("|------|------|------|\n")
	sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", palaceMarkdown(p, 8), palaceMarkdown(p, 1), palaceMarkdown(p, 6)))

	sb.WriteString("\n> 本盘面仅供学习参考，疾病、法律、财务等重大决策请咨询现实专业人士。\n")

	return sb.String()
}

func palaceMarkdown(p *models.Plate, num int) string {
	for _, palace := range p.Palaces {
		if palace.Number == num {
			rel := ""
			if palace.StemRelation != "" {
				rel = fmt.Sprintf("(%s)", palace.StemRelation)
			}
			return fmt.Sprintf("%s%s%s<br>%s/%s<br>%s %s %s", palace.HeavenStem, rel, palace.EarthStem, palace.Gua, palace.Branch, palace.Star, palace.Door, palace.Spirit)
		}
	}
	return ""
}

func dayStemOf(gz string) string {
	runes := []rune(gz)
	if len(runes) > 0 {
		return string(runes[0])
	}
	return ""
}

func hourStemOf(gz string) string {
	runes := []rune(gz)
	if len(runes) > 0 {
		return string(runes[0])
	}
	return ""
}

// RunPlate 执行排盘子命令
func RunPlate(timeStr, format string) error {
	t, err := parseTime(timeStr)
	if err != nil {
		return err
	}

	p, err := plate.BuildPlate(t)
	if err != nil {
		return err
	}

	switch format {
	case "json":
		data, err := json.MarshalIndent(p, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "markdown", "md":
		fmt.Println(PlateToMarkdown(p))
	default:
		return fmt.Errorf("不支持的输出格式: %s", format)
	}
	return nil
}
