package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"qimen-tool/internal/plate"
	"qimen-tool/internal/verify"
)

// RunVerify 执行盘面自检
func RunVerify(timeStr, format string) error {
	t, err := parseTime(timeStr)
	if err != nil {
		return err
	}

	p, err := plate.BuildPlate(t)
	if err != nil {
		return err
	}

	results := verify.Verify(p)

	switch format {
	case "json":
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "markdown", "md":
		fmt.Println(verifyToMarkdown(p, results))
	default:
		return fmt.Errorf("不支持的输出格式: %s", format)
	}
	return nil
}

func verifyToMarkdown(p interface{}, results []verify.Result) string {
	var sb strings.Builder
	sb.WriteString("# 排盘自检报告\n\n")
	passed := 0
	failed := 0
	for _, r := range results {
		status := "✅"
		if !r.Pass {
			status = "❌"
			failed++
		} else {
			passed++
		}
		sb.WriteString(fmt.Sprintf("- %s **%s** (%s)：%s\n", status, r.Category, passLabel(r.Pass), r.Message))
	}
	sb.WriteString(fmt.Sprintf("\n总计：%d 项通过，%d 项失败\n", passed, failed))
	return sb.String()
}

func passLabel(pass bool) string {
	if pass {
		return "通过"
	}
	return "失败"
}
