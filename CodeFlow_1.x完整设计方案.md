# CodeFlow 1.x Workflow Pack 完整设计方案

> 版本线定位：CodeFlow 1.x = Workflow Pack  
> 核心目标：把 OpenSpec + Superpowers + Claude Code / Codex / Cursor 的 AI 开发规则、技能、协作模式和质量门以轻量模板包形式安装到项目中。  
> 当前稳定目标版本：CodeFlow 1.2 Skill System

---

## 1. 产品定位

CodeFlow 1.x 是轻量级 AI Coding Workflow Pack，不是完整平台。

它负责把以下能力安装到项目内：

- Claude Code adapter
- Codex adapter
- Cursor adapter
- OpenSpec + Superpowers 工作流规则
- Graphify 可选增强规则
- Collaborative Agents 协作代理规范
- Vue + TypeScript 前端工程 Skill
- Go 后端工程 Skill
- Context Budget / Quality Gates / Skill Learning
- 自动 Skill Routing
- `.codeflow/` 基础工作流文件

一句话：

```text
CodeFlow 1.x = 给项目安装 AI 开发规范、技能和工具适配的轻量规则包。
```

---

## 2. 1.x 与 2.x 的边界

### CodeFlow 1.x 负责

```text
1. 安装 Claude Code / Codex / Cursor 的配置文件
2. 安装 rules / commands / agents / skills / hooks
3. 安装 .codeflow/workflows
4. 安装工程技能包
5. 提供 Collaborative Agents 使用规范
6. 提供 Graphify 可选增强规则
7. 提供团队可提交的项目内 AI 规范
```

### CodeFlow 1.x 不负责

```text
1. 不提供 Go CLI 产品能力
2. 不提供 Web Dashboard
3. 不提供 SQLite + FTS5 全文索引
4. 不提供多项目管理
5. 不提供可视化需求/迭代管理
6. 不提供项目画像自动扫描
7. 不提供运行时数据库状态管理
```

---

## 3. 版本规划

### 1.0：Multi-Tool Workflow Pack

核心能力：

```text
1. Claude Code adapter
2. Codex adapter
3. Cursor adapter
4. .codeflow/ 通用状态目录
5. Graphify 可选增强
6. Git 安全规则
7. 简单需求 / 复杂需求工作流
8. OpenSpec + Superpowers 基础集成
```

### 1.1：Collaborative Agents

新增能力：

```text
1. Claude Code Agent Teams 规则和命令模板
2. Codex Subagent Workflows 规则和 skill
3. Cursor Parallel Agents 规则
4. Collaborative Agents 通用工作流
5. parallel-agent-review
6. collaborative-agent-development
7. agent team cleanup / plan approval / worktree 策略
```

### 1.2：Skill System

新增能力：

```text
1. codeflow-frontend-vue-ts 完整技能包
2. codeflow-backend-go 完整技能包
3. codeflow-context-budget
4. codeflow-quality-gates
5. codeflow-skill-learning
6. codeflow-eval-checkpoints
7. codeflow-auto-skill-routing
8. Cursor globs 自动加载规则
9. Codex AGENTS.md skill routing
10. Claude rules + skills routing
```

---

## 4. 安装后目录结构

```text
target-project/
├── CLAUDE.md
├── AGENTS.md
├── .codeflow/
│   ├── state.md
│   ├── active-change.md
│   ├── workflows/
│   └── prompts/
├── .claude/
│   ├── codeflow/
│   ├── rules/
│   ├── commands/
│   ├── agents/
│   ├── skills/
│   └── settings.json
├── .agents/
│   └── skills/
└── .cursor/
    └── rules/
```

---

## 5. 核心设计原则

```text
1. Skills-first：skills 是核心能力单元
2. Commands 只是快捷入口
3. Rules 是硬约束
4. Hooks 做自动化、安全门和上下文保存
5. Agents 做专项执行者
6. Collaborative Agents 用于复杂并行协作
7. 不自动 Git 写操作
8. 不自动执行 Graphify
9. 简单需求不使用多代理
10. 当前任务只加载相关 skill
```

---

## 6. Adapter 设计

### 6.1 Claude Code Adapter

安装内容：

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

Claude Code 专属能力：

```text
1. slash commands
2. agents
3. skills
4. hooks
5. Agent Teams guidance
```

### 6.2 Codex Adapter

安装内容：

```text
AGENTS.md
.agents/skills/
```

Codex 使用方式：

```text
按 CodeFlow 开始新需求：<需求内容>
按 CodeFlow 复杂需求流程处理：<需求内容>
按 CodeFlow 使用 subagents 审查当前分支
```

Codex 不安装 hooks，不支持 `/codeflow-*` commands。

### 6.3 Cursor Adapter

安装内容：

```text
.cursor/rules/*.mdc
AGENTS.md
```

Cursor 使用方式：

```text
按 CodeFlow 开始新需求：<需求内容>
按 CodeFlow 并行 Agent 模式审查当前变更
```

Cursor 不安装 hooks，不支持 `/codeflow-*` commands。

---

## 7. 工作流设计

### 7.1 简单需求流程

适用：

```text
1. 需求明确
2. 改动范围小
3. 不涉及核心链路
4. 不涉及接口 / DB / 权限 / 安全
5. 可以快速验证
```

流程：

```text
brainstorming 轻量确认
→ Existing Capability Discovery 轻量检查
→ writing-plans 轻量计划
→ 询问执行方式
→ TDD
→ Review
→ verification，可选
→ finishing，可选
```

默认执行方式：

```text
executing-plans
```

### 7.2 复杂需求流程

适用：

```text
1. 涉及多模块
2. 涉及核心业务链路
3. 涉及接口、数据库、配置、权限、安全
4. 涉及并发、数据一致性、发布风险
5. 需要正式验收标准
```

流程：

```text
brainstorming
→ Existing Capability Discovery
→ /opsx:propose <change-id>
→ Spec Review
→ writing-plans
→ 询问执行方式
→ TDD
→ Review
→ requesting-code-review，可选
→ verification-before-completion，可选
→ finishing-a-development-branch，可选
→ 用户确认后 archive
```

---

## 8. 执行方式设计

CodeFlow 1.x 支持四种执行方式：

```text
1. executing-plans
2. subagent-driven-development
3. parallel-agent-review
4. collaborative-agent-development
```

| 执行方式 | 适用场景 |
|---|---|
| executing-plans | 简单需求、小范围修改、顺序执行 |
| subagent-driven-development | 专项任务拆分，只需要结果回报 |
| parallel-agent-review | 多角色并行审查，如安全、性能、测试 |
| collaborative-agent-development | 跨模块开发、架构调查、疑难问题排查 |

---

## 9. Collaborative Agents 设计

### 9.1 统一概念

CodeFlow 统一称为：

```text
Collaborative Agents
```

三端适配：

| 工具 | 原生能力 | CodeFlow 适配 |
|---|---|---|
| Claude Code | Agent Teams | Team lead + Teammates |
| Codex | Subagent Workflows | specialized subagents |
| Cursor | Parallel Agents / Subagents / Cloud Agents | rules + prompt guidance |

### 9.2 通用规则

```text
1. 简单需求不要使用多代理
2. 多代理只用于复杂 Review、跨模块开发、架构评审、疑难问题调查
3. 多代理默认先做研究、审查、方案和测试建议
4. 多代理写代码前必须用户确认
5. 写代码前必须说明每个代理的文件范围
6. 避免多个代理修改同一文件
7. 并行写代码优先使用 git worktree
8. 高风险任务必须先 plan approval
9. 所有代理完成后必须输出 final report
10. 主代理必须综合结果，不允许直接拼接输出
11. 多代理 token / 成本更高，创建前必须提醒用户
```

### 9.3 Claude Code Agent Teams

规则：

```text
1. 仅 Claude Code 支持 Agent Teams
2. 默认不启用，只安装规则和命令模板
3. 不生成项目级 .claude/teams/
4. 创建团队前必须用户确认
5. 高风险任务必须 plan approval
6. 完成后必须由 Team lead cleanup
7. 不允许多个 teammate 修改同一文件
8. Agent Teams 不替代 subagents
```

Claude commands：

```text
.claude/commands/codeflow-team-review.md
.claude/commands/codeflow-team-investigate.md
.claude/commands/codeflow-team-feature.md
.claude/commands/codeflow-team-cleanup.md
```

### 9.4 Codex Subagent Workflows

安装：

```text
.agents/skills/codeflow-collaborative-agents/
.agents/skills/codeflow-codex-subagent-workflows/
```

规则：

```text
1. Codex 只在用户明确要求时 spawn subagents
2. 每个 subagent 必须有角色、范围、输出格式
3. 默认用于 read-heavy 分析、测试、triage、review
4. 并行写代码前必须获得用户确认
5. 主线程必须等待所有 subagents 结果后再总结
```

### 9.5 Cursor Parallel Agents

安装：

```text
.cursor/rules/codeflow-collaborative-agents.mdc
.cursor/rules/codeflow-agent-prompts.mdc
```

规则：

```text
1. Cursor 不使用 Claude Agent Teams runtime
2. 通过 .cursor/rules 和 prompt 模板约束
3. 复杂任务可以使用 parallel agents
4. 每个 agent 必须有明确角色、范围、输出格式
5. 并行写代码前必须用户确认
6. 推荐使用 worktrees 避免冲突
```

---

## 10. Skill System 设计

### 10.1 Skills-first 架构

```text
Skills = 核心能力单元
Commands = 触发 skill 的快捷入口
Rules = 始终遵守的硬约束
Hooks = 自动化、安全门、状态保存
Agents = 专项执行者
```

### 10.2 Skill 分类

```text
Workflow Skills
├── codeflow-openspec
├── codeflow-discovery
├── codeflow-tdd
├── codeflow-review
├── codeflow-graphify
├── codeflow-handoff
├── codeflow-collaborative-agents
├── codeflow-context-budget
├── codeflow-quality-gates
├── codeflow-skill-learning
└── codeflow-eval-checkpoints

Engineering Skills
├── codeflow-frontend-vue-ts
└── codeflow-backend-go

Project Skills
├── codeflow-project-patterns
├── codeflow-api-contract
└── codeflow-release-check
```

---

## 11. Frontend Vue + TypeScript Skill

### 11.1 触发条件

```text
package.json
vite.config.ts
nuxt.config.ts
src/**/*.vue
src/**/*.ts
src/**/*.tsx
components/
pages/
layouts/
stores/
router/
tailwind.config.*
uno.config.*
```

### 11.2 技能包结构

```text
codeflow-frontend-vue-ts/
├── SKILL.md
├── references/
│   ├── ui-layout.md
│   ├── design-antipatterns.md
│   ├── vue-component-patterns.md
│   ├── state-management.md
│   ├── api-contract.md
│   ├── security.md
│   └── testing.md
├── checklists/
│   └── frontend-review.md
└── examples/
    ├── dashboard-page.md
    ├── detail-page.md
    └── table-filter-page.md
```

### 11.3 核心规则

```text
1. 页面结构清晰
2. 组件边界明确
3. TypeScript 类型完整
4. 不使用 any
5. 状态管理可控
6. API 层独立封装
7. loading / empty / error / disabled 状态完整
8. 不写临时 Demo 式 UI
9. 不 card 套 card
10. 不在前端硬编码 token / secret
```

---

## 12. Backend Go Skill

### 12.1 触发条件

```text
go.mod
*.go
cmd/
internal/
pkg/
api/
service/
repository/
model/
middleware/
job/
worker/
cron/
```

### 12.2 技能包结构

```text
codeflow-backend-go/
├── SKILL.md
├── references/
│   ├── go-style.md
│   ├── error-handling.md
│   ├── context.md
│   ├── interface-design.md
│   ├── concurrency.md
│   ├── api-service-repository.md
│   ├── security.md
│   └── testing.md
├── checklists/
│   └── go-review.md
└── examples/
    ├── service.md
    ├── repository.md
    ├── handler.md
    └── table-driven-test.md
```

### 12.3 核心规则

```text
1. 使用领域类型表达业务含义
2. 让非法状态不可表示
3. context 贯穿 IO 边界
4. 错误使用 %w 包装
5. 日志结构化
6. interface 定义在使用方
7. handler / service / repository 分层清楚
8. goroutine 生命周期明确
9. table-driven tests
10. 不打印 secret / token
11. API 字段小驼峰 json tag
12. 不写 Java/PHP 味 Go
```

---

## 13. 自动 Skill Routing

### 13.1 Claude Code

```text
.claude/rules/codeflow-auto-skill-routing.md
```

规则：

```text
涉及 Vue / TypeScript / Vite / Nuxt / Web UI → codeflow-frontend-vue-ts
涉及 Go / go.mod / API / service / repository / job → codeflow-backend-go
全栈改动 → 同时使用两个 skill
```

### 13.2 Codex

```text
AGENTS.md
.agents/skills/
```

规则：

```text
Vue / TS / Web UI → .agents/skills/codeflow-frontend-vue-ts
Go backend → .agents/skills/codeflow-backend-go
Full-stack → 同时使用两个 skill
```

### 13.3 Cursor

使用 `.mdc` globs：

```text
.cursor/rules/codeflow-frontend-vue-ts.mdc
.cursor/rules/codeflow-backend-go.mdc
```

---

## 14. Context Budget

新增 skill：

```text
codeflow-context-budget
```

规则：

```text
1. 当前任务只加载相关 skill
2. 简单需求禁止加载过多 workflow skill
3. 多代理默认不超过最小必要数量
4. MCP 默认不超过必要数量
5. Review 和执行分开上下文
6. 长任务必须 handoff
7. compact 前必须生成上下文摘要
```

---

## 15. Quality Gates / Eval Checkpoints

新增：

```text
codeflow-quality-gates
codeflow-eval-checkpoints
```

检查：

```text
1. 前端改动必须检查 loading / empty / error / disabled / selected
2. Go 改动必须检查 gofmt / go test / error wrapping / context / logging
3. Review 输出必须有 Must Fix / Should Fix / Test Gap / Risk / Evidence
4. 每个任务完成后必须说明验证命令
5. 不能验证时必须说明原因和替代验证方式
```

---

## 16. Skill Learning

新增：

```text
codeflow-skill-learning
```

生成目录：

```text
.codeflow/learnings/
├── frontend.md
├── backend-go.md
├── project-patterns.md
└── review-findings.md
```

规则：

```text
1. 不自动写入 learnings
2. 发现重复问题时提醒用户是否沉淀
3. 写入前要求用户确认
4. learning 必须包含场景、错误模式、正确做法、示例
```

---

## 17. 安全边界

### 前端

```text
1. 不在前端硬编码 token / secret
2. API 错误信息不泄露内部信息
3. v-html 必须谨慎
4. 表单输入必须校验
5. 权限按钮不能只靠前端控制
```

### Go

```text
1. 不打印 secret / token
2. SQL / ES query / Mongo filter 防注入
3. HTTP handler 参数校验
4. 权限校验不放在外层临时判断
5. goroutine / channel / context 避免泄漏
6. 日志脱敏
```

---

## 18. Hooks 设计

Claude Code hooks：

```text
SessionStart
UserPromptSubmit
PreToolUse
SubagentStop
Stop
TaskCreated
TaskCompleted
TeammateIdle
```

用途：

```text
1. 注入 state
2. 提醒 skill routing
3. 阻止 Git 写操作
4. 阻止 Graphify 自动执行
5. 提醒 TDD / Review
6. Agent Teams 任务质量门
7. Stop 时提醒是否沉淀 learning
```

---

## 19. 团队提交建议

建议提交：

```text
.claude/rules/
.claude/commands/
.claude/agents/
.claude/skills/
.agents/skills/
.cursor/rules/
.codeflow/workflows/
.codeflow/prompts/
AGENTS.md
CLAUDE.md
```

不建议提交：

```text
.codeflow/state.md
.codeflow/active-change.md
.codeflow/tmp/
.codeflow/logs/
graphify-out/
```

---

## 20. 验收标准

### 20.1 Adapter 验收

```text
1. Claude Code adapter 可安装
2. Codex adapter 可安装
3. Cursor adapter 可安装
4. Git 安全规则三端一致
5. Graphify 提醒三端一致
```

### 20.2 Collaborative Agents 验收

```text
1. Claude Agent Teams 命令模板存在
2. Codex Subagent Workflows skill 存在
3. Cursor Parallel Agents rule 存在
4. 多代理通用规则存在
5. 不生成 .claude/teams/
6. 默认不启用 Claude Agent Teams
```

### 20.3 Skill System 验收

```text
1. Frontend Vue+TS skill 完整存在
2. Backend Go skill 完整存在
3. Context Budget skill 存在
4. Quality Gates skill 存在
5. Skill Learning skill 存在
6. Auto Skill Routing 规则存在
7. Cursor globs 正确
8. Codex AGENTS.md routing 正确
```

---

## 21. 不做范围

CodeFlow 1.x 不做：

```text
1. 不做 Web Dashboard
2. 不做 SQLite
3. 不做多项目管理
4. 不做运行时 agent 管理
5. 不自动安装 Graphify
6. 不自动执行 Git 操作
7. 不自动写入 learnings
```

---

## 22. 最终结论

CodeFlow 1.x 是稳定可发布的 Workflow Pack 版本线。

最新完整能力由 1.2 承载：

```text
CodeFlow 1.2 = Multi-Tool Workflow Pack + Collaborative Agents + Skill System
```

适合团队直接提交到项目仓库中，统一 Claude Code / Codex / Cursor 的 AI 开发规范、前后端工程技能、并行协作模式和质量门。
