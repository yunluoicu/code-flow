# CodeFlow 2.1

CodeFlow 2.1 是一个面向 AI Coding 的本地项目工作台：Go CLI + 多 AI Adapter + 项目画像 + 项目索引 + OpenSpec/Superpowers 导入 + SQLite/FTS5 + Web Dashboard。

## 定位

```text
CodeFlow 1.x = Workflow Pack
CodeFlow 2.1 = AI Development Workspace
```

2.0 不是 MVP，包含完整可用能力：

- Go CLI：`codeflow`
- Claude Code / Codex / Cursor adapter 安装
- `.codeflow/` 项目工作区
- 项目画像 `project.yaml / project.md`
- 模块索引 `modules.yaml / index.json`
- OpenSpec 导入
- Superpowers 过程文档导入
- 需求管理
- 迭代管理
- Review / Check Gate
- Graphify 状态检查
- SQLite + FTS5 全文搜索
- Web Dashboard，默认端口 `4399`

## 安装

### go install

```bash
go install github.com/yunluoicu/code-flow/cmd/codeflow@latest
```

### 源码编译

```bash
git clone https://github.com/yunluoicu/code-flow
cd code-flow
go build -o codeflow ./cmd/codeflow
./codeflow version
```

## 快速开始

```bash
# 初始化当前项目，安装三端 adapter，并生成项目画像和索引
codeflow init --tools claude,codex,cursor

# 查看状态
codeflow status

# 重新生成项目画像
codeflow profile

# 重新生成模块索引
codeflow index

# 同步 OpenSpec / Superpowers / docs / code 到 SQLite
codeflow sync

# 启动 Web Dashboard
codeflow web --port 4399
```

## 常用命令

```bash
codeflow version

codeflow init
codeflow init --tools claude,codex,cursor
codeflow init --dry-run
codeflow init --force

codeflow doctor
codeflow status
codeflow upgrade
codeflow uninstall --force

codeflow profile
codeflow index
codeflow sync

codeflow requirement new --title "需求标题"
codeflow requirement list
codeflow requirement show <id>

codeflow iteration new --name "迭代名称"
codeflow iteration list
codeflow iteration show <id>

codeflow changes list
codeflow changes show <change-id>
codeflow changes check <change-id>

codeflow check
codeflow review

codeflow graph status
codeflow graph suggest

codeflow web
codeflow web --port 4399
codeflow web --workspace ~/projects
```

## Web Dashboard

默认：

```bash
codeflow web
```

访问：

```text
http://127.0.0.1:4399
```

支持：

- 项目列表
- 项目详情
- 项目画像
- 模块列表
- 需求列表和详情
- 迭代列表和详情
- OpenSpec Specs / Changes
- Superpowers 过程记录
- Review 结果
- Check 风险
- Graphify 状态
- 全文搜索

## 安全边界

CodeFlow 不会自动执行：

- `git add`
- `git commit`
- `git push`
- `git merge`
- `git rebase`
- `git reset`
- `git clean`
- `/graphify .`
- `/graphify . --update`

Graphify 只做状态检查和建议，不自动运行。

## 目录

```text
.codeflow/
├── manifest.json
├── config.yaml
├── project.yaml
├── project.md
├── modules.yaml
├── index.json
├── state.md
├── active-change.md
├── requirements/
├── iterations/
├── decisions/
├── reviews/
├── checks/
├── workflows/
├── tmp/
└── logs/
```

全局数据库：

```text
~/.codeflow/codeflow.db
```

## Web Dashboard 视觉升级

CodeFlow 2.1 Web Dashboard 不是简单 HTML 页面，而是正式的本地开发工作台。

当前发布包包含：

```text
Go JSON API
嵌入式 Web Dashboard
左侧 Sidebar
顶部搜索
项目 KPI
项目卡片
项目详情
需求 / 迭代 / OpenSpec / Review / Graphify / 搜索页面
```

运行方式不变：

```bash
codeflow web --port 4399
```

前端源码位于：

```text
web/
```

内置静态发布资源位于：

```text
internal/core/webui/dist/
```


## Collaborative Agents

CodeFlow 2.1 新增 Collaborative Agents 模块：

- Claude Code：Agent Teams readiness / commands / hooks guidance
- Codex：Subagent Workflows + `.agents/skills`
- Cursor：Parallel Agents + `.cursor/rules`

新增 CLI：

```bash
codeflow agents status
codeflow agents suggest "复杂 Review"
```

Dashboard 增加 Collaborative Agents 页面，展示三端能力、规则、启用建议和风险提示。
