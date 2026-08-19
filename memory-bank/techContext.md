# 技术上下文

## 确定技术栈
- **核心排盘与 CLI**：**Go**（单一二进制，零依赖运行）。
  - 命令行工具：提供 `plate`、`interpret`、`rules`、`verify` 四个子命令。
  - 输出支持 **JSON**（供 agent/程序解析）和 **Markdown/文本**（供人阅读）。
- **历法计算**：
  - Go 标准库 `time` 处理公历日期。
  - 农历/干支/节气使用 `github.com/6tail/lunar-go`，保证准确且无需维护复杂历法算法。
- **无 Web 前端**：
  - 纯 CLI 工具，无需 Gin、HTTP、浏览器。
  - 可视化通过 Markdown 九宫格表格实现。
- **数据存储**：**不使用数据库**。
  - 规则库、格局库、用神说明用 **JSON 静态文件** + `go:embed` 内嵌。
  - 案例保存可导出为 JSON 文件，由用户本地管理。
- **Skill 包装**：
  - 用户级 pi skill：`C:/Users/22569/.agents/skills/qimen-tool/SKILL.md`。
  - 项目内副本：`C:/Users/22569/Documents/Study/QimenDj/SKILL.md`。
  - 全局可执行文件：`C:/Users/22569/.local/bin/qimen.exe`（已加入 PATH）。
  - 通过 `qimen` 命令直接调用，无需进入项目目录。

## 约束
- 需要准确的节气算法；传统节气按太阳黄经 15° 划分。
- 干支、空亡、马星等规则需与经典文献一致。
- 盘面可视化需要支持中文字符清晰展示。

## 开发工具
- 版本控制：Git
- 测试：对照已出版的奇门遁甲实例进行回归测试
- 文档：Markdown（本记忆库）
