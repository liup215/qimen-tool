package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"qimen-tool/internal/interpretation"
)

// RunRules 执行规则查询子命令
// query 可为：door/门名、star/星名、spirit/神名、pattern/格局、topic/用神主题、jiugong/宫名、sanqi/奇名
func RunRules(query, format string) error {
	rules := interpretation.Rules()
	q := strings.TrimSpace(query)

	result := make(map[string]interface{})
	found := false

	if v, ok := rules.Doors[q]; ok {
		result["type"] = "door"
		result["name"] = q
		result["data"] = v
		found = true
	} else if v, ok := rules.Stars[q]; ok {
		result["type"] = "star"
		result["name"] = q
		result["data"] = v
		found = true
	} else if v, ok := rules.Spirits[q]; ok {
		result["type"] = "spirit"
		result["name"] = q
		result["data"] = v
		found = true
	} else if v, ok := rules.YongShen[q]; ok {
		result["type"] = "topic"
		result["name"] = q
		result["data"] = v
		found = true
	} else if v, ok := rules.Patterns.NamedPatterns[q]; ok {
		result["type"] = "named_pattern"
		result["name"] = q
		result["data"] = v
		found = true
	} else if v, ok := rules.Patterns.JiugongXiangyi[q]; ok {
		result["type"] = "jiugong"
		result["name"] = q
		result["data"] = v
		found = true
	} else if v, ok := rules.Patterns.SanqiXiangyi[q]; ok {
		result["type"] = "sanqi"
		result["name"] = q
		result["data"] = v
		found = true
	} else if q == "topics" {
		topics := interpretation.YongShenTopics()
		sort.Strings(topics)
		result["type"] = "topic-list"
		result["topics"] = topics
		found = true
	} else if q == "patterns" || q == "格局" {
		result["type"] = "patterns"
		result["data"] = map[string]string{
			"空亡":    rules.Patterns.KongWangMeaning,
			"马星":    rules.Patterns.MaXingMeaning,
			"伏吟":    rules.Patterns.FuYin,
			"反吟":    rules.Patterns.FanYin,
			"旺相休囚": rules.Patterns.WangXiang.Description,
		}
		found = true
	} else if q == "named_patterns" || q == "格局列表" {
		names := make([]string, 0, len(rules.Patterns.NamedPatterns))
		for k := range rules.Patterns.NamedPatterns {
			names = append(names, k)
		}
		sort.Strings(names)
		result["type"] = "named-pattern-list"
		result["patterns"] = names
		found = true
	} else if q == "jiugong" || q == "九宫" {
		keys := make([]string, 0, len(rules.Patterns.JiugongXiangyi))
		for k := range rules.Patterns.JiugongXiangyi {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		result["type"] = "jiugong-list"
		result["gua"] = keys
		found = true
	} else if q == "sanqi" || q == "三奇" {
		keys := make([]string, 0, len(rules.Patterns.SanqiXiangyi))
		for k := range rules.Patterns.SanqiXiangyi {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		result["type"] = "sanqi-list"
		result["qi"] = keys
		found = true
	}

	if !found {
		return fmt.Errorf("未找到规则: %s。可查询：门（如 休）、星（如 天蓬）、神（如 值符）、用神主题（如 career）、topics、patterns、named_patterns、jiugong、sanqi", q)
	}

	switch format {
	case "json":
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "markdown", "md":
		fmt.Println(RulesToMarkdown(result))
	default:
		return fmt.Errorf("不支持的输出格式: %s", format)
	}
	return nil
}

// RulesToMarkdown 规则查询结果转 Markdown
func RulesToMarkdown(result map[string]interface{}) string {
	var sb strings.Builder
	t := result["type"].(string)
	switch t {
	case "door", "star", "spirit":
		name := result["name"].(string)
		data := result["data"].(interpretation.DoorStarSpirit)
		sb.WriteString(fmt.Sprintf("# %s（%s）\n\n", name, typeLabel(t)))
		if data.WuXing != "" {
			sb.WriteString(fmt.Sprintf("- 五行：%s\n", data.WuXing))
		}
		sb.WriteString(fmt.Sprintf("- 吉凶：%s\n", data.JiXiong))
		sb.WriteString(fmt.Sprintf("- 象义：%s\n", data.XiangYi))
	case "topic":
		data := result["data"].(interpretation.YongShenRule)
		sb.WriteString(fmt.Sprintf("# 用神主题：%s\n\n", data.Description))
		for _, line := range data.YongShen {
			sb.WriteString(fmt.Sprintf("- %s\n", line))
		}
	case "topic-list":
		topics := result["topics"].([]string)
		sb.WriteString("# 可用用神主题\n\n")
		for _, topic := range topics {
			sb.WriteString(fmt.Sprintf("- %s\n", topic))
		}
	case "patterns":
		sb.WriteString("# 通用格局含义\n\n")
		data := result["data"].(map[string]string)
		for k, v := range data {
			sb.WriteString(fmt.Sprintf("## %s\n\n%s\n\n", k, v))
		}
	case "named_pattern":
		name := result["name"].(string)
		data := result["data"].(interpretation.NamedPattern)
		sb.WriteString(fmt.Sprintf("# %s（%s）\n\n", name, data.Nature))
		sb.WriteString(fmt.Sprintf("- 天盘干：%s\n", data.Heaven))
		sb.WriteString(fmt.Sprintf("- 地盘干：%s\n", data.Earth))
		sb.WriteString(fmt.Sprintf("- 含义：%s\n", data.Description))
	case "named-pattern-list":
		names := result["patterns"].([]string)
		sb.WriteString("# 具名格局列表\n\n")
		for _, name := range names {
			sb.WriteString(fmt.Sprintf("- %s\n", name))
		}
	case "jiugong":
		name := result["name"].(string)
		data := result["data"].(string)
		sb.WriteString(fmt.Sprintf("# %s宫象意\n\n%s\n", name, data))
	case "jiugong-list":
		gua := result["gua"].([]string)
		sb.WriteString("# 九宫象意列表\n\n")
		for _, g := range gua {
			sb.WriteString(fmt.Sprintf("- %s\n", g))
		}
	case "sanqi":
		name := result["name"].(string)
		data := result["data"].(string)
		sb.WriteString(fmt.Sprintf("# %s奇象意\n\n%s\n", name, data))
	case "sanqi-list":
		qi := result["qi"].([]string)
		sb.WriteString("# 三奇象意列表\n\n")
		for _, q := range qi {
			sb.WriteString(fmt.Sprintf("- %s\n", q))
		}
	}
	return sb.String()
}

func typeLabel(t string) string {
	switch t {
	case "door":
		return "门"
	case "star":
		return "星"
	case "spirit":
		return "神"
	}
	return t
}
