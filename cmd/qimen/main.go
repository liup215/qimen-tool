package main

import (
	"flag"
	"fmt"
	"os"

	"qimen-tool/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "plate":
		plateFlags := flag.NewFlagSet("plate", flag.ExitOnError)
		timeStr := plateFlags.String("t", "now", "时间，如 2024-08-19 14:30 或 now")
		format := plateFlags.String("f", "json", "输出格式：json 或 markdown")
		plateFlags.Parse(os.Args[2:])
		if err := cli.RunPlate(*timeStr, *format); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	case "interpret":
		interpretFlags := flag.NewFlagSet("interpret", flag.ExitOnError)
		timeStr := interpretFlags.String("t", "now", "时间，如 2024-08-19 14:30 或 now")
		topic := interpretFlags.String("topic", "general", "用神主题，如 general/wealth/career/marriage/health 等")
		format := interpretFlags.String("f", "json", "输出格式：json 或 markdown")
		interpretFlags.Parse(os.Args[2:])
		if err := cli.RunInterpret(*timeStr, *topic, *format); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	case "rules":
		rulesFlags := flag.NewFlagSet("rules", flag.ExitOnError)
		query := rulesFlags.String("q", "topics", "查询项：门名、星名、神名、用神主题、topics、patterns")
		format := rulesFlags.String("f", "json", "输出格式：json 或 markdown")
		rulesFlags.Parse(os.Args[2:])
		if err := cli.RunRules(*query, *format); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	case "verify":
		verifyFlags := flag.NewFlagSet("verify", flag.ExitOnError)
		timeStr := verifyFlags.String("t", "now", "时间，如 2024-08-19 14:30 或 now")
		format := verifyFlags.String("f", "json", "输出格式：json 或 markdown")
		verifyFlags.Parse(os.Args[2:])
		if err := cli.RunVerify(*timeStr, *format); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`奇门遁甲 CLI 工具

用法:
  qimen plate [选项]
  qimen interpret [选项]
  qimen rules [选项]
  qimen verify [选项]

plate 选项:
  -t string    排盘时间，默认 now（当前北京时间）
               支持格式: 2006-01-02 15:04:05, 2006-01-02 15:04, 2006-01-02
  -f string    输出格式: json（默认）或 markdown

interpret 选项:
  -t string    排盘时间，默认 now
  -topic string 用神主题: general（默认）/ wealth / career / marriage / health / travel / lawsuit / lost / study
  -f string    输出格式: json（默认）或 markdown

rules 选项:
  -q string    查询项：门名/星名/神名/用神主题/topics/patterns，默认 topics
  -f string    输出格式: json（默认）或 markdown

verify 选项:
  -t string    排盘时间，默认 now
  -f string    输出格式：json（默认）或 markdown

示例:
  qimen plate -t "2024-08-19 14:30" -f markdown
  qimen interpret -t "2024-08-19 14:30" -topic career -f markdown
  qimen rules -q 天蓬 -f markdown
  qimen verify -t "2024-08-19 14:30"`)
}
