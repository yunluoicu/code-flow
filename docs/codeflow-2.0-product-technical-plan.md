# CodeFlow 2.0 产品与技术方案

> 版本定位：CodeFlow 2.0 = AI Development Workspace  
> 目标：把 AI Coding 项目的工具安装、项目画像、项目索引、需求迭代、OpenSpec/Superpowers 过程文档、Review/Check、SQLite 全文索引和 Web Dashboard 统一管理起来。

---

## 1. 背景与目标

### 1.1 背景

CodeFlow 1.x 已经定义为 Workflow Pack，主要解决：

- Claude Code / Codex / Cursor 的 AI 工作流规则安装
- OpenSpec + Superpowers 开发约束
- Graphify 可选增强规则
- 多 AI 工具 adapter 安装
- `.codeflow/` 基础状态文件

但 1.x 本质仍是“规则包 / 安装包”，能力集中在项目内 AI 开发约束，并不负责：

- 多项目管理
- 项目画像和模块索引
- 需求和迭代管理
- OpenSpec/Superpowers 文档聚合
- Review / Check 风险聚合
- Web 可视化
- SQLite 全文搜索

CodeFlow 2.0 升级为完整本地工作台，用 CLI + Web Dashboard 管理 AI Coding 项目的全生命周期。

---

### 1.2 产品定位

```text
CodeFlow 1.x = Workflow Pack
CodeFlow 2.0 = AI Development Workspace
```

CodeFlow 2.0 是一个面向 AI Coding 的本地项目工作台，负责把以下内容统一管理：

- 项目基本信息
- 项目技术栈
- 功能模块
- AI 工具 adapter
- OpenSpec specs / changes
- Superpowers 过程文档
- 需求
- 迭代
- Review 结果
- Check 风险
- Graphify 状态
- SQLite + FTS5 全文索引
- Web Dashboard

一句话：

```text
CodeFlow 2.0 = AI 项目工作流 + 项目知识索引 + Web 管理台
```

---

### 1.3 核心原则

```text
1. 使用 Go 实现
2. 必须包含 Web Dashboard
3. 使用 SQLite + FTS5
4. 保留并复用 CodeFlow 1.x templates
5. Web 默认端口 4399，同时支持自定义
6. 不自动执行 Git 写操作
7. 不自动执行 Graphify
8. 不替代 OpenSpec
9. 不替代 Superpowers
10. 不强依赖 LLM，静态分析必须可用
```

---

## 2. 版本边界

### 2.1 CodeFlow 1.x：Workflow Pack

CodeFlow 1.x 负责：

```text
1. 安装 Claude Code adapter
2. 安装 Codex adapter
3. 安装 Cursor adapter
4. 安装 .codeflow 基础状态目录
5. 安装 rules / commands / agents / skills / hooks
6. 提供 Graphify 可选增强规则
```

CodeFlow 1.x 不负责：

```text
1. Web Dashboard
2. SQLite 索引
3. 多项目管理
4. 项目画像
5. 需求管理
6. 迭代管理
7. OpenSpec/Superpowers 文档聚合展示
```

---

### 2.2 CodeFlow 2.0：AI Development Workspace

CodeFlow 2.0 负责：

```text
1. Go CLI
2. 多 AI adapter 安装管理
3. 项目画像生成
4. 项目模块索引
5. OpenSpec 文档导入
6. Superpowers 过程文档导入
7. 需求管理
8. 迭代管理
9. Review / Check Gate
10. Graphify 状态管理
11. SQLite + FTS5 全文索引
12. Web Dashboard
13. 多项目管理
14. 全文搜索
```

---

## 3. 用户与使用场景

### 3.1 目标用户

| 用户 | 诉求 |
|---|---|
| 技术负责人 | 统一 AI 开发流程，查看项目需求、方案、Review、风险 |
| Go 后端开发 | 通过 CLI 初始化项目 AI 工作流 |
| AI Coding 使用者 | 让 Claude Code / Codex / Cursor 按统一规范工作 |
| 项目维护者 | 查看已有 OpenSpec / Superpowers 文档和项目模块 |
| 团队成员 | 通过 Web Dashboard 查看项目状态和迭代进展 |

---

### 3.2 核心场景

#### 场景 1：初始化已有项目

```bash
codeflow init --tools claude,codex,cursor
```

执行内容：

```text
1. 安装 AI adapter
2. 扫描 openspec/
3. 扫描 Superpowers 过程文档
4. 扫描 docs / README
5. 扫描源码结构
6. 检测 Graphify 状态
7. 生成 project.yaml / project.md / modules.yaml
8. 建立 SQLite 索引
```

#### 场景 2：生成项目画像

```bash
codeflow profile
```

输出：

```text
.codeflow/project.yaml
.codeflow/project.md
```

#### 场景 3：生成模块索引

```bash
codeflow index
```

输出：

```text
.codeflow/modules.yaml
.codeflow/index.json
```

#### 场景 4：同步已有资料

```bash
codeflow sync
codeflow sync --openspec
codeflow sync --superpowers
codeflow sync --docs
codeflow sync --graphify
```

#### 场景 5：管理需求

```bash
codeflow requirement new
codeflow requirement list
codeflow requirement show REQ-20260524-001
```

#### 场景 6：管理迭代

```bash
codeflow iteration new
codeflow iteration list
codeflow iteration show ITER-202605
```

#### 场景 7：启动 Dashboard

```bash
codeflow web
codeflow web --port 4399
codeflow web --port 18080
codeflow web --workspace ~/projects
```

---

## 4. 总体架构

### 4.1 架构分层

```text
CodeFlow 2.0
├── CLI Layer
│   ├── init
│   ├── profile
│   ├── index
│   ├── sync
│   ├── requirement
│   ├── iteration
│   ├── check
│   └── web
├── Adapter Layer
│   ├── Claude Code Adapter
│   ├── Codex Adapter
│   └── Cursor Adapter
├── Project Intelligence Layer
│   ├── OpenSpec Importer
│   ├── Superpowers Importer
│   ├── Docs Scanner
│   ├── Source Scanner
│   └── Graphify Status
├── Storage Layer
│   ├── .codeflow files
│   └── SQLite + FTS5
└── Web Dashboard
    ├── Project Catalog
    ├── Project Detail
    ├── Requirements
    ├── Iterations
    ├── OpenSpec
    ├── Superpowers
    ├── Review / Check
    └── Search
```

---

### 4.2 数据流

```text
项目目录
├── openspec/
├── Superpowers 过程文档
├── docs/
├── README.md
├── 源码
├── graphify-out/
└── .codeflow/
        ↓
codeflow sync / index / profile
        ↓
.codeflow/*.yaml / *.md / *.json
        ↓
SQLite + FTS5
        ↓
Web Dashboard
        ↓
项目画像 / 需求 / 迭代 / OpenSpec / Review / Graphify 状态
```

---

## 5. 技术选型

### 5.1 后端与 CLI

| 模块 | 技术 |
|---|---|
| 语言 | Go |
| CLI | cobra |
| 配置 | yaml / json |
| Web Server | net/http + chi 或 gin |
| DB | SQLite + FTS5 |
| 静态资源 | go:embed |
| 模板文件 | go:embed |
| 项目扫描 | 文件系统扫描 + AST 可选 |
| 日志 | slog 或 zap |

---

### 5.2 Go 项目结构

```text
code-flow/
├── go.mod
├── cmd/
│   └── codeflow/
│       └── main.go
├── internal/
│   ├── app/
│   ├── cli/
│   ├── config/
│   ├── installer/
│   ├── adapters/
│   │   ├── claude/
│   │   ├── codex/
│   │   └── cursor/
│   ├── project/
│   ├── scanner/
│   ├── openspec/
│   ├── superpowers/
│   ├── graphify/
│   ├── storage/
│   ├── indexer/
│   ├── web/
│   └── check/
├── templates/
│   ├── core/
│   ├── claude-code/
│   ├── codex/
│   └── cursor/
├── web/
│   ├── src/
│   └── dist/
└── docs/
```

---

### 5.3 安装方式

#### go install

```bash
go install github.com/yunluoicu/code-flow/cmd/codeflow@latest
```

指定版本：

```bash
go install github.com/yunluoicu/code-flow/cmd/codeflow@v2.0.0
```

验证：

```bash
codeflow version
```

#### 源码编译

```bash
git clone https://github.com/yunluoicu/code-flow
cd code-flow
go build -o codeflow ./cmd/codeflow
./codeflow version
```

#### 二进制发布

后续提供：

```text
codeflow-darwin-arm64
codeflow-darwin-amd64
codeflow-linux-amd64
codeflow-linux-arm64
codeflow-windows-amd64.exe
```

---

## 6. CLI 设计

### 6.1 命令总览

```bash
codeflow version

codeflow init
codeflow init --tools claude,codex,cursor
codeflow init --dry-run
codeflow init --force

codeflow doctor
codeflow status
codeflow upgrade
codeflow uninstall

codeflow profile
codeflow index
codeflow sync

codeflow requirement new
codeflow requirement list
codeflow requirement show <id>

codeflow iteration new
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

---

### 6.2 `codeflow init`

#### 作用

初始化当前项目的 CodeFlow 工作区。

#### 功能

```text
1. 生成 .codeflow/
2. 生成 manifest.json
3. 生成 config.yaml
4. 安装 AI adapter
5. 扫描 openspec/
6. 扫描 Superpowers 过程文档
7. 扫描 docs / README
8. 扫描源码结构
9. 生成 project.yaml / project.md / modules.yaml
10. 建立 SQLite 索引
```

#### 参数

```bash
codeflow init --tools claude
codeflow init --tools codex
codeflow init --tools cursor
codeflow init --tools claude,codex,cursor
codeflow init --dry-run
codeflow init --force
codeflow init --port 4399
```

#### 行为要求

```text
1. 不修改业务代码
2. 不自动执行 Git 写操作
3. 不自动运行 Graphify
4. 不覆盖已有 CLAUDE.md / AGENTS.md，只追加 CodeFlow 片段
5. 不删除已有 .claude/settings.json，只合并 hooks
```

---

### 6.3 `codeflow doctor`

#### 作用

检查 CodeFlow 自身安装是否正常。

#### 检查项

```text
1. codeflow 版本
2. Go runtime 信息
3. SQLite 是否可用
4. FTS5 是否可用
5. 全局目录 ~/.codeflow 是否存在
6. 全局 DB 是否可访问
7. 模板是否存在
8. Web 默认端口是否可用
```

---

### 6.4 `codeflow check`

#### 作用

检查当前项目 AI 开发流程风险。

#### 检查项

```text
1. 是否已安装 CodeFlow
2. manifest.json 是否存在
3. adapter 是否完整
4. openspec/ 是否存在
5. OpenSpec change 是否未完成
6. tasks.md 是否未完成
7. Superpowers verification 是否缺失
8. Review 是否存在 must-fix
9. Graphify 是否不存在或可能过期
10. project.md / modules.yaml 是否落后
```

---

### 6.5 `codeflow profile`

#### 作用

生成或更新项目画像。

#### 输出

```text
.codeflow/project.yaml
.codeflow/project.md
```

#### 数据来源

```text
1. openspec/specs/
2. openspec/changes/
3. Superpowers 过程文档
4. docs/ / README.md
5. 源码扫描
6. Graphify 可选
```

---

### 6.6 `codeflow index`

#### 作用

生成项目模块索引。

#### 输出

```text
.codeflow/modules.yaml
.codeflow/index.json
```

#### 扫描内容

```text
1. 语言
2. 框架
3. 模块目录
4. 入口文件
5. 配置文件
6. 测试文件
7. OpenSpec specs
8. OpenSpec changes
9. Graphify 状态
```

---

### 6.7 `codeflow sync`

#### 作用

重新导入项目资料并更新 SQLite / FTS5 索引。

#### 参数

```bash
codeflow sync
codeflow sync --openspec
codeflow sync --superpowers
codeflow sync --docs
codeflow sync --graphify
```

---

### 6.8 `codeflow requirement`

#### 命令

```bash
codeflow requirement new
codeflow requirement list
codeflow requirement show <id>
```

#### 需求状态

```text
draft
planned
executing
review
done
archived
```

---

### 6.9 `codeflow iteration`

#### 命令

```bash
codeflow iteration new
codeflow iteration list
codeflow iteration show <id>
```

#### 迭代状态

```text
planning
active
release-ready
released
archived
```

---

### 6.10 `codeflow web`

#### 命令

```bash
codeflow web
codeflow web --port 4399
codeflow web --port 18080
codeflow web --workspace ~/projects
```

#### 端口优先级

```text
命令行参数 > .codeflow/config.yaml > ~/.codeflow/config.yaml > 默认 4399
```

---

## 7. 项目目录设计

### 7.1 项目内 `.codeflow/`

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

---

### 7.2 全局目录

```text
~/.codeflow/
├── config.yaml
├── codeflow.db
├── logs/
└── cache/
```

---

### 7.3 Git 提交建议

| 文件 | 建议 |
|---|---|
| `.codeflow/config.yaml` | 可以提交 |
| `.codeflow/project.yaml` | 可以提交 |
| `.codeflow/project.md` | 可以提交 |
| `.codeflow/modules.yaml` | 可以提交 |
| `.codeflow/index.json` | 可以提交 |
| `.codeflow/workflows/` | 可以提交 |
| `.codeflow/requirements/` | 可以提交 |
| `.codeflow/iterations/` | 可以提交 |
| `.codeflow/decisions/` | 可以提交 |
| `.codeflow/reviews/` | 可以提交 |
| `.codeflow/checks/` | 可以提交 |
| `.codeflow/state.md` | 默认不提交 |
| `.codeflow/active-change.md` | 默认不提交 |
| `.codeflow/tmp/` | 不提交 |
| `.codeflow/logs/` | 不提交 |
| `graphify-out/` | 默认不提交 |

---

## 8. 项目画像设计

### 8.1 `project.yaml`

```yaml
name: customer-server
type: backend-service
description: AI 客服后台服务
language:
  - go
framework:
  - gin
database:
  - mongodb
  - redis
  - elasticsearch
aiTools:
  - claude
  - codex
  - cursor
workflow:
  openspec: true
  superpowers: true
  graphify: optional
modules:
  - session
  - faq
  - ai-reply
  - rag-search
  - statistics
  - permission
```

---

### 8.2 `project.md`

```md
# 项目介绍

## 项目定位

## 技术栈

## 核心模块

## 主要业务流程

## OpenSpec 概览

## Superpowers 过程文档概览

## AI 开发注意事项
```

---

### 8.3 `modules.yaml`

```yaml
modules:
  session:
    path:
      - service/session*
      - models/session*
    responsibilities:
      - 会话状态
      - 消息轮次
      - 转人工策略
    relatedSpecs:
      - openspec/specs/session/spec.md
    relatedChanges:
      - openspec/changes/adjust-human-handoff-rules
```

---

## 9. 数据源优先级

`codeflow profile / index / sync` 必须按以下优先级处理：

```text
P0 openspec/specs/
P0 openspec/changes/
P0 Superpowers 过程文档
P1 .codeflow/requirements/
P1 .codeflow/iterations/
P1 docs/ / README.md
P2 源码扫描
P2 Graphify
```

核心原则：

```text
OpenSpec 说明系统承诺了什么
Superpowers 说明需求怎么讨论、计划、执行、验证
源码说明系统实际实现了什么
Graphify 说明代码可能怎么关联
CodeFlow 负责统一整理、索引和展示
```

---

## 10. OpenSpec 导入逻辑

### 10.1 扫描路径

```text
openspec/specs/
openspec/changes/
```

### 10.2 解析文件

```text
proposal.md
design.md
tasks.md
spec.md
```

### 10.3 导入结果

生成：

```text
.codeflow/index.json
SQLite openspec_specs
SQLite openspec_changes
SQLite documents_fts
```

Dashboard 展示：

```text
Specs 当前能力
Active Changes 进行中变更
Archived Changes 已归档变更
Tasks 完成情况
Spec Review 状态
```

### 10.4 状态识别

| 状态 | 判断方式 |
|---|---|
| active | `openspec/changes/<change-id>/` 存在且未归档 |
| task-incomplete | `tasks.md` 存在未勾选项 |
| design-missing | 缺少 `design.md` |
| proposal-missing | 缺少 `proposal.md` |
| spec-delta-missing | 缺少 `specs/*/spec.md` |
| archive-ready | tasks 完成、review 无 must-fix、验证记录存在 |

---

## 11. Superpowers 导入逻辑

### 11.1 支持的过程文档类型

```text
brainstorming
writing-plans
executing-plans
subagent-driven-development
test-driven-development
requesting-code-review
verification-before-completion
finishing-a-development-branch
```

### 11.2 导入后对应关系

| Superpowers 产物 | CodeFlow 数据 |
|---|---|
| brainstorming | requirement.discovery / requirement.scope |
| writing-plans | implementation plan |
| executing-plans | execution events |
| subagent-driven-development | subtask events |
| test-driven-development | test result |
| requesting-code-review | stage review |
| verification-before-completion | verification result |
| finishing-a-development-branch | finish report |

### 11.3 Dashboard 展示

```text
需求讨论
实施计划
任务执行记录
TDD 结果
阶段审查
完成前验证
分支收尾报告
```

---

## 12. Graphify 集成

### 12.1 原则

```text
Graphify 是可选增强
Graphify 不作为事实来源
不自动执行 Graphify
不自动提交 graphify-out/
```

### 12.2 状态检查

```bash
codeflow graph status
```

检查：

```text
1. 是否存在 graphify-out/graph.json
2. graphify-out/graph.json 是否可能过期
3. 是否存在 GRAPH_REPORT.md
```

### 12.3 缺失提示

如果不存在：

```text
graphify-out/graph.json
```

提示：

```text
当前项目未发现 graphify-out/graph.json。
如需增强项目已有能力发现，可以执行：

/graphify .
```

### 12.4 过期提示

如果最近代码文件修改时间晚于 `graphify-out/graph.json`：

```text
检测到部分代码文件修改时间晚于 graphify-out/graph.json，图谱可能已落后。
如需更新，可以执行：

/graphify . --update
```

---

## 13. SQLite + FTS5 设计

### 13.1 数据库位置

```text
~/.codeflow/codeflow.db
```

### 13.2 表结构

```sql
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    type TEXT,
    description TEXT,
    language TEXT,
    framework TEXT,
    ai_tools TEXT,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE project_modules (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    path TEXT,
    responsibilities TEXT,
    related_specs TEXT,
    related_changes TEXT,
    updated_at DATETIME
);

CREATE TABLE requirements (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    type TEXT,
    status TEXT,
    change_id TEXT,
    iteration_id TEXT,
    source TEXT,
    file_path TEXT,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE iterations (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    status TEXT,
    release_date TEXT,
    file_path TEXT,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE openspec_specs (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    capability TEXT NOT NULL,
    file_path TEXT NOT NULL,
    summary TEXT,
    updated_at DATETIME
);

CREATE TABLE openspec_changes (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    change_id TEXT NOT NULL,
    status TEXT,
    proposal_path TEXT,
    design_path TEXT,
    tasks_path TEXT,
    specs_path TEXT,
    task_total INTEGER,
    task_done INTEGER,
    updated_at DATETIME
);

CREATE TABLE superpowers_records (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    requirement_id TEXT,
    change_id TEXT,
    record_type TEXT,
    title TEXT,
    file_path TEXT,
    summary TEXT,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE reviews (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    requirement_id TEXT,
    change_id TEXT,
    review_type TEXT,
    conclusion TEXT,
    must_fix_count INTEGER,
    should_fix_count INTEGER,
    test_gap_count INTEGER,
    risk_level TEXT,
    file_path TEXT,
    created_at DATETIME
);

CREATE TABLE checks (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    check_type TEXT,
    status TEXT,
    message TEXT,
    detail TEXT,
    created_at DATETIME
);

CREATE TABLE documents (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    doc_type TEXT,
    title TEXT,
    file_path TEXT,
    content TEXT,
    updated_at DATETIME
);

CREATE VIRTUAL TABLE documents_fts USING fts5(
    title,
    content,
    doc_type,
    project_id UNINDEXED,
    file_path UNINDEXED
);
```

### 13.3 全文搜索范围

```text
项目介绍
模块说明
需求
迭代
OpenSpec specs
OpenSpec changes
Superpowers 过程记录
Review 结果
Check 结果
docs / README
```

---

## 14. Web Dashboard 设计

### 14.1 启动方式

```bash
codeflow web
codeflow web --port 4399
codeflow web --port 18080
codeflow web --workspace ~/projects
```

默认：

```text
host: 127.0.0.1
port: 4399
```

### 14.2 页面清单

```text
1. 项目列表
2. 项目详情
3. 项目画像
4. 模块列表
5. 需求列表
6. 需求详情
7. 迭代列表
8. 迭代详情
9. OpenSpec Specs
10. OpenSpec Changes
11. Superpowers 过程记录
12. Review 结果
13. Check 风险
14. Graphify 状态
15. 全文搜索
16. 设置页
```

### 14.3 项目列表页

展示：

```text
项目名称
项目路径
项目类型
技术栈
AI 工具
OpenSpec 状态
Graphify 状态
未完成需求数
最近更新时间
```

### 14.4 项目详情页

展示：

```text
项目介绍
技术栈
功能模块
OpenSpec specs
active changes
需求列表
迭代列表
Review 状态
Graphify 状态
```

### 14.5 需求详情页

展示：

```text
需求背景
目标
In Scope
Out of Scope
验收标准
关联 OpenSpec change
关联实施计划
关联 TDD
关联 Review
关联 Verification
当前状态
```

### 14.6 OpenSpec 页面

展示：

```text
Specs 当前能力
Changes 变更记录
proposal.md
design.md
tasks.md
spec delta
任务完成度
archive 状态
```

### 14.7 Superpowers 页面

展示：

```text
brainstorming
writing-plans
executing-plans
subagent-driven-development
test-driven-development
requesting-code-review
verification-before-completion
finishing-a-development-branch
```

### 14.8 Review / Check 页面

展示：

```text
Must Fix
Should Fix
Test Gap
Risk
Next Step
状态趋势
```

### 14.9 Graphify 页面

展示：

```text
graphify-out 是否存在
graph.json 更新时间
是否可能过期
GRAPH_REPORT.md 摘要
建议命令
```

### 14.10 搜索页

支持全文搜索：

```text
项目
模块
需求
迭代
OpenSpec
Superpowers
Review
Check
docs
```

---

## 15. AI Adapter 设计

### 15.1 Claude Code

安装：

```text
.claude/
├── codeflow/
├── rules/
├── commands/
├── agents/
├── skills/
└── settings.json
```

根 `CLAUDE.md` 追加：

```md
<!-- CodeFlow start -->
@.claude/codeflow/CLAUDE.md
<!-- CodeFlow end -->
```

### 15.2 Codex

安装：

```text
AGENTS.md
.agents/skills/
```

不安装 hooks，不提供 `/codeflow-*` commands。

使用方式：

```text
按 CodeFlow 开始新需求：<需求内容>
```

### 15.3 Cursor

安装：

```text
.cursor/rules/*.mdc
AGENTS.md
```

不安装 hooks，不提供 `/codeflow-*` commands。

使用方式：

```text
按 CodeFlow 开始新需求：<需求内容>
```

---

## 16. Check Gate 设计

### 16.1 `codeflow check`

检查：

```text
1. CodeFlow 是否安装
2. AI adapter 是否完整
3. OpenSpec 是否存在
4. OpenSpec changes 是否未完成
5. tasks.md 是否未完成
6. Superpowers verification 是否缺失
7. Review 是否存在 must-fix
8. Graphify 是否缺失或可能过期
9. project.md / modules.yaml 是否可能落后
10. SQLite 索引是否需要刷新
```

### 16.2 输出格式

```text
Check Result

Status: warning

Risks:
1. openspec/changes/adjust-ai-human-handoff/tasks.md 未完成
2. Review 存在 1 个 must-fix
3. graphify-out/graph.json 可能已过期

Suggested Actions:
1. 完成 tasks.md
2. 修复 must-fix 后重新 review
3. 执行 /graphify . --update
```

---

## 17. 安全与边界

### 17.1 不自动执行 Git 写操作

禁止自动执行：

```text
git add
git commit
git push
git merge
git rebase
git reset
git clean
git restore
git checkout
git switch
git stash
git tag
```

### 17.2 不自动执行 Graphify

禁止自动执行：

```text
/graphify .
/graphify . --update
```

### 17.3 Web 安全边界

2.0 Web Dashboard 默认：

```text
host: 127.0.0.1
port: 4399
无登录
本地访问
```

不做：

```text
权限系统
用户系统
远程团队协作
公网访问
```

---

## 18. 验收标准

### 18.1 CLI 验收

```text
1. go install 后可执行 codeflow
2. codeflow version 可显示版本
3. codeflow init 可安装 Claude / Codex / Cursor adapter
4. codeflow init --dry-run 不写入文件
5. codeflow doctor 可检查本机环境
6. codeflow status 可查看当前项目状态
7. codeflow profile 可生成 project.yaml / project.md
8. codeflow index 可生成 modules.yaml / index.json
9. codeflow sync 可导入 OpenSpec / Superpowers / docs / code
10. codeflow requirement new/list/show 可用
11. codeflow iteration new/list/show 可用
12. codeflow changes list/show/check 可用
13. codeflow check 可输出项目风险
14. codeflow graph status/suggest 可用
15. codeflow web 可启动 Dashboard
```

### 18.2 Web 验收

```text
1. 默认端口 4399 可访问
2. 支持 --port 自定义端口
3. 可查看所有项目
4. 可查看项目详情
5. 可查看项目画像
6. 可查看模块列表
7. 可查看需求列表和详情
8. 可查看迭代列表和详情
9. 可查看 OpenSpec specs / changes
10. 可查看 Superpowers 过程记录
11. 可查看 Review 结果
12. 可查看 Check 风险
13. 可查看 Graphify 状态
14. 可全文搜索项目、需求、方案、OpenSpec、Review
```

### 18.3 数据验收

```text
1. .codeflow/manifest.json 存在
2. .codeflow/config.yaml 存在
3. .codeflow/project.yaml 存在
4. .codeflow/project.md 存在
5. .codeflow/modules.yaml 存在
6. .codeflow/index.json 存在
7. SQLite 数据库可创建
8. FTS5 可搜索
9. OpenSpec 文档可导入
10. Superpowers 过程文档可导入
```

### 18.4 Adapter 验收

```text
1. Claude Code adapter 可安装
2. Claude Code commands / rules / agents / skills / hooks 可用
3. Codex adapter 可安装
4. Codex AGENTS.md + .agents/skills 可用
5. Cursor adapter 可安装
6. Cursor .cursor/rules/*.mdc 可用
7. 三端 Git 安全规则一致
8. 三端 Graphify 提醒一致
```

---

## 19. 不做范围

CodeFlow 2.0 不做：

```text
1. 不做远程团队协作
2. 不做登录和权限
3. 不做公网部署
4. 不做 MCP Server
5. 不做 CI/CD
6. 不自动安装 OpenSpec
7. 不自动安装 Superpowers
8. 不自动安装 Graphify
9. 不自动执行 Git 操作
10. 不自动执行 Graphify
```

---

## 20. 最终结论

CodeFlow 2.0 定义为完整产品，而不是 MVP。

必须同时具备：

```text
Go CLI
AI Adapter 安装
项目画像
项目索引
OpenSpec / Superpowers 导入
需求管理
迭代管理
Review / Check
Graphify 状态
SQLite + FTS5
Web Dashboard
多项目管理
全文搜索
```

CodeFlow 1.x 继续作为轻量 Workflow Pack 维护。

CodeFlow 2.0 作为完整 AI Development Workspace 独立维护，并复用 1.x templates。
