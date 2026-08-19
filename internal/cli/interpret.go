package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"qimen-tool/internal/interpretation"
	"qimen-tool/internal/plate"
)

// RunInterpret 执行解释子命令
func RunInterpret(timeStr, topic, format string) error {
	t, err := parseTime(timeStr)
	if err != nil {
		return err
	}

	p, err := plate.BuildPlate(t)
	if err != nil {
		return err
	}

	report, err := interpretation.Interpret(p, topic)
	if err != nil {
		return err
	}

	switch format {
	case "json":
		data, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "markdown", "md":
		fmt.Println(InterpretToMarkdown(report))
	default:
		return fmt.Errorf("不支持的输出格式: %s", format)
	}
	return nil
}

// InterpretToMarkdown 将解释报告转为 Markdown
func sectionNumber(resources, threats, patterns int, veto bool) string {
	n := 4
	if resources > 0 {
		n++
	}
	if threats > 0 {
		n++
	}
	if patterns > 0 {
		n++
	}
	if veto {
		n++
	}
	n++ // 综合判断是最后一个章节，再 +1
	chinese := []string{"", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十"}
	if n <= 10 {
		return chinese[n]
	}
	return fmt.Sprintf("%d", n)
}

func InterpretToMarkdown(r *interpretation.Report) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# 奇门遁甲分析报告\n\n"))
	sb.WriteString(fmt.Sprintf("- 时间：%s\n", r.Meta.Time))
	sb.WriteString(fmt.Sprintf("- 问题类型：%s\n\n", r.Meta.Description))

	sb.WriteString("## 一、自身状态\n\n")
	sb.WriteString(fmt.Sprintf("- 用神：%s\n", r.Self.Symbol))
	sb.WriteString(fmt.Sprintf("- 落宫：%d宫\n", r.Self.Palace))
	sb.WriteString(fmt.Sprintf("- 状态：%s\n", r.Self.Level))
	sb.WriteString(fmt.Sprintf("- 说明：%s\n\n", strings.Join(r.Self.Notes, "，")))

	sb.WriteString("## 二、所问之事状态\n\n")
	sb.WriteString(fmt.Sprintf("- 用神：%s\n", r.Target.Symbol))
	sb.WriteString(fmt.Sprintf("- 落宫：%d宫\n", r.Target.Palace))
	sb.WriteString(fmt.Sprintf("- 状态：%s\n", r.Target.Level))
	sb.WriteString(fmt.Sprintf("- 说明：%s\n\n", strings.Join(r.Target.Notes, "，")))

	sb.WriteString("## 三、人与事的关系\n\n")
	sb.WriteString(fmt.Sprintf("- 关系：%s\n", r.Relationship.Type))
	sb.WriteString(fmt.Sprintf("- 结论：%s\n\n", r.Relationship.Description))

	if len(r.Resources) > 0 {
		sb.WriteString("## 四、资源与助力\n\n")
		for _, f := range r.Resources {
			sb.WriteString(fmt.Sprintf("- %s落%d宫（%s）：%s\n", f.Symbol, f.Palace, f.RelationToSelf, f.Note))
		}
		sb.WriteString("\n")
	}

	if len(r.Threats) > 0 {
		sb.WriteString("## 五、威胁与隐患\n\n")
		for _, f := range r.Threats {
			sb.WriteString(fmt.Sprintf("- %s落%d宫（%s）：%s\n", f.Symbol, f.Palace, f.RelationToSelf, f.Note))
		}
		sb.WriteString("\n")
	}

	if len(r.DetectedPatterns) > 0 {
		sb.WriteString("## 六、格局提示\n\n")
		for _, pat := range r.DetectedPatterns {
			sb.WriteString(fmt.Sprintf("- **%s**（%s）临%d宫：%s\n", pat.Name, pat.Nature, pat.Palace, pat.Description))
		}
		sb.WriteString("\n")
	}

	activeVetos := false
	for _, v := range r.VetoChecks {
		if v.Active {
			activeVetos = true
			break
		}
	}
	if activeVetos {
		sb.WriteString("## 七、一票否决\n\n")
		for _, v := range r.VetoChecks {
			if v.Active {
				sb.WriteString(fmt.Sprintf("- **%s**：%s\n", v.Name, v.Message))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("## %s、综合判断\n\n", sectionNumber(len(r.Resources), len(r.Threats), len(r.DetectedPatterns), activeVetos)))
	sb.WriteString(fmt.Sprintf("%s\n\n", r.Summary))
	sb.WriteString(fmt.Sprintf("> %s\n", r.Caution))

	return sb.String()
}
