# 奇门遁甲 CLI 工具

基于 **时家奇门转盘法**、**拆补法定局** 的 Go CLI 工具。规则集版本 `mainline-cn-v1`：转盘法、拆补法、中五宫寄坤二宫、北京时间默认时区、值符值使用时柱旬首。

无数据库、无 Web、单一二进制，适合 Agent/Skill 调用。

## 功能

- `qimen plate`：根据公历时间排盘，输出 JSON 或 Markdown。
- `qimen interpret`：按主题（事业、财运、婚姻、健康等）输出通用规则解释报告。
- `qimen rules`：查询八门、九星、八神、用神主题、具名格局、九宫象意、三奇象意。
- `qimen verify`：对排好的盘面执行格式、逻辑、一致性自检。

## 安装

```bash
go build -o qimen.exe ./cmd/qimen
```

## 快速开始

```bash
# 排盘（JSON，方便程序解析）
./qimen.exe plate -t "2024-08-19 14:30" -f json

# 排盘（Markdown，方便阅读）
./qimen.exe plate -t "2024-08-19 14:30" -f markdown

# 当前时间排盘
./qimen.exe plate -t now -f markdown

# 按主题解释
./qimen.exe interpret -t "2024-08-19 14:30" -topic career -f markdown

# 查询规则
./qimen.exe rules -q 青龙返首 -f markdown

# 排盘自检
./qimen.exe verify -t "2024-08-19 14:30" -f markdown
```

## 用法

### `qimen plate`

根据公历时间生成奇门遁甲盘面。

```bash
qimen plate -t "2024-08-19 14:30" -f json
```

参数：
- `-t`：时间，支持 `now`、`YYYY-MM-DD HH:MM`、`YYYY-MM-DD HH:MM:SS`、`YYYY-MM-DD`。默认北京时间。
- `-f`：输出格式，`json` 或 `markdown`。默认 `json`。

输出包含：
- 四柱、节气、阴阳遁、局数
- 九宫的天地盘干、八门、九星、八神、空亡、马星
- 规则集版本、中宫寄坤说明
- 反查索引：`door_index`、`star_index`、`spirit_index`、`stem_index`
- 日/时干落宫、`stem_relation`（天地盘干生克关系）

### `qimen interpret`

输出通用规则解释报告，不替代具体事件判断。

```bash
qimen interpret -t "2024-08-19 14:30" -topic career -f markdown
```

支持的 topic：`general`、`wealth`、`career`、`marriage`、`health`、`travel`、`lawsuit`、`lost`、`study`。

报告结构：
1. 自身状态（日干）
2. 所问之事状态（用神）
3. 人与事的关系（宫位五行生克）
4. 资源与助力（值符、值使、吉神等）
5. 威胁与隐患（白虎、玄武、螣蛇等）
6. 格局提示（青龙返首、白虎猖狂等）
7. 一票否决（五不遇时、三奇入墓、天网四张、全局伏吟/反吟）
8. 综合判断

### `qimen rules`

查询规则库。

```bash
qimen rules -q 天蓬 -f markdown           # 九星
qimen rules -q 休 -f markdown             # 八门
qimen rules -q 值符 -f markdown           # 八神
qimen rules -q career -f markdown         # 用神主题
qimen rules -q patterns -f markdown       # 通用格局含义
qimen rules -q named_patterns -f markdown # 具名格局列表
qimen rules -q 青龙返首 -f markdown       # 具体格局
qimen rules -q 坎 -f markdown             # 九宫象意
qimen rules -q 乙 -f markdown             # 三奇象意
qimen rules -q topics                     # 所有用神主题
```

### `qimen verify`

对指定时间的盘面执行自检链。

```bash
qimen verify -t "2024-08-19 14:30" -f markdown
```

覆盖三类断言：
- **format**：规则集版本、宫位数量与顺序
- **logic**：八门、九星、八神完整性
- **consistency**：值符星/值使门与落宫一致、反查索引一致、日/时干落宫一致

## 技术说明

- **历法**：使用 `github.com/6tail/lunar-go` 计算四柱和节气。
- **定局**：拆补法，按十二节（立春、惊蛰、清明等）换局，上元子午卯酉、中元寅申巳亥、下元辰戌丑未。
- **值符值使**：使用时柱旬首。
- **排盘**：转盘法，中五宫寄坤二宫，天禽星寄天芮星。
- **解释**：基于规则库输出结构化报告，具体事件结论由 LLM 根据盘面自行组织。
- **规则库**：JSON 静态文件 + `go:embed` 内嵌，便于无依赖分发。
- **测试**：`tests/` 下包含单元测试、逻辑测试、回归测试（golden cases）和边界时刻测试。

## 项目结构

```
.
├── cmd/qimen/              CLI 入口
├── internal/
│   ├── assets/             规则库 JSON
│   ├── calendar/           历法、干支、节气
│   ├── cli/                子命令实现
│   ├── interpretation/     解释引擎
│   ├── models/             数据模型
│   ├── plate/              排盘核心
│   └── verify/             排盘自检
├── tests/                  测试用例
├── memory-bank/            项目记忆库
└── README.md
```

## 测试

```bash
go test ./... -v
```

## 验证与限制

- 本工具采用 **转盘法**，不是飞盘法。
- 定局采用 **拆补法**，不是置闰法。
- 中五宫寄坤二宫。
- 已覆盖子时、午时、节气交界、跨年夜等边界时刻，均通过 `verify` 自检。
- 具体事件结论需要结合用神、月令、五行生克等进一步分析，本工具只提供排盘与通用规则。

## 免责声明

本工具仅供学习参考，不替代专业判断。疾病、法律、财务等重大决策请咨询现实专业人士。

## 许可证

MIT
