# CodeFlow 使用指南

Commands、Agents、Skills、Rules 完整详解。

---

## 目录

- [一、Commands（命令）](#一commands命令)
- [二、Agents（专项代理）](#二agents专项代理)
- [三、Skills（核心能力）](#三skills核心能力)
- [四、Rules（硬约束）](#四rules硬约束)

---

## 一、Commands（命令）

Slash commands 是工作流的快捷入口，通过 `/codeflow-*` 触发。仅 Claude Code 支持。

### 基础命令

| 命令                   | 一句话        | 使用场景                 |
|----------------------|------------|----------------------|
| `/codeflow-new <描述>` | 启动新需求      | 接到新需求时，替代直接写代码       |
| `/codeflow-continue` | 继续上次未完成的需求 | 上次会话中断、上下文压缩后恢复      |
| `/codeflow-status`   | 查看当前需求进度   | 不记得进度时确认、同步状态前查看     |
| `/codeflow-review`   | 审查当前分支改动   | 代码改动完成后、提交前质量把关      |
| `/codeflow-finish`   | 完成需求，收尾归档  | 代码写完、review 通过、测试通过后 |
| `/codeflow-handoff`  | 保存上下文摘要    | 长任务切换、上下文压缩前、多人协作    |

### Agent Team 命令

需要 `export CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`，仅 Claude Code 支持。

| 命令                                | 一句话             | 使用场景                       |
|-----------------------------------|-----------------|----------------------------|
| `/codeflow-team-review <范围>`      | Agent Team 并行审查 | 重要 PR 多维度审查（安全/性能/测试/可维护性） |
| `/codeflow-team-feature <描述>`     | Agent Team 协作开发 | 跨模块大功能开发                   |
| `/codeflow-team-investigate <问题>` | Agent Team 疑难调查 | 难复现 Bug、生产诡异问题、多假设排查       |
| `/codeflow-team-cleanup`          | 清理 Agent Team   | 协作完成后确保无遗留活跃 agent         |

### 使用示例

```text
# 启动简单需求
/codeflow-new 修复订单列表分页参数丢失问题

# 启动复杂需求
/codeflow-new 新增用户角色权限管理系统

# 审查当前改动
/codeflow-review

# Agent Team 并行审查
/codeflow-team-review 审查当前分支对支付模块的改动

# Agent Team 调查问题
/codeflow-team-investigate 支付成功后订单状态仍为"待支付"
```

---

## 二、Agents（专项代理）

6 个专项 agent，通过 Agent Team 命令或手动 Task spawn 使用。**核心约束：只分析/审查，不默认修改代码。**

| Agent                          | 角色            | 职责                                    |
|--------------------------------|---------------|---------------------------------------|
| `codeflow-requirement-analyst` | 需求分析专家        | 需求边界澄清、功能点拆解、影响范围评估、风险识别              |
| `codeflow-code-reviewer`       | 代码审查专家        | 代码质量、安全漏洞、性能问题、可维护性，输出五维度报告           |
| `codeflow-openspec-reviewer`   | OpenSpec 审查专家 | Spec 完整性、接口契约一致性、边界场景覆盖               |
| `codeflow-discovery-agent`     | 能力发现专家        | 搜索已有实现、API、组件、工具函数，避免重复造轮子            |
| `codeflow-plan-reviewer`       | 计划审查专家        | 架构合理性、步骤可行性、风险遗漏、测试策略                 |
| `codeflow-test-reviewer`       | 测试审查专家        | 测试覆盖、边界场景、回归风险、table-driven tests 规范性 |

### 输出格式

所有 agent 遵循统一输出格式：

```md
## 审查结果

### Must Fix

### Should Fix

### Test Gap

### Risk

### Evidence
```

---

## 三、Skills（核心能力）

Skills 是 CodeFlow 架构的核心能力单元。Commands 只是 Skills 的快捷入口。

### Workflow Skills（工作流技能）

| Skill                           | 触发时机            | 做什么                                       |
|---------------------------------|-----------------|-------------------------------------------|
| `codeflow-tdd`                  | 有代码逻辑改动时        | 先写测试再写实现。Go: table-driven tests；Vue: 组件测试 |
| `codeflow-review`               | 代码改动完成后         | Review 流程、五维度输出规范                         |
| `codeflow-discovery`            | 新需求开始时          | 搜索已有实现/API/组件，避免重复                        |
| `codeflow-openspec`             | 复杂需求            | Propose → spec-review → archive           |
| `codeflow-graphify`             | 可选增强            | 知识图谱索引与查询，不自动执行                           |
| `codeflow-handoff`              | 长任务/跨会话         | 生成紧凑上下文摘要，供后续恢复                           |
| `codeflow-collaborative-agents` | 复杂 Review/开发/调查 | 定义四种执行方式与三端协作规则                           |
| `codeflow-context-budget`       | 始终生效            | 控制 skill 加载数量、代理数量，防止上下文膨胀                |
| `codeflow-quality-gates`        | 任务完成前           | 强制检查 TDD/Review/状态完整/gofmt/test           |
| `codeflow-eval-checkpoints`     | 关键节点            | 计划完成/子任务完成/Review 前/发布前检查点                |
| `codeflow-skill-learning`       | 发现重复问题时         | 提醒用户沉淀模式到 .codeflow/learnings/            |

### Engineering Skills（工程技能）

| Skill                      | 触发条件                                                  | 核心规则                                           |
|----------------------------|-------------------------------------------------------|------------------------------------------------|
| `codeflow-frontend-vue-ts` | `*.vue` / `*.ts` / `vite.config.ts` / `components/` 等 | 页面结构清晰、组件边界明确、类型完整、loading/empty/error 状态齐全    |
| `codeflow-backend-go`      | `*.go` / `go.mod` / `cmd/` / `internal/` 等            | 类型清晰、错误可追踪、context 贯穿、结构化日志、table-driven tests |

### 自动路由

Claude Code 会根据当前文件自动加载对应工程 Skill，无需手动指定：

- 涉及 Vue / TS / Vite / Web UI → 自动加载 `codeflow-frontend-vue-ts`
- 涉及 Go / go.mod / 后端服务 → 自动加载 `codeflow-backend-go`
- 全栈改动 → 同时加载两个 Skill

---

## 四、Rules（硬约束）

Rules 始终生效，无需手动触发。

### 核心约束

| Rule                | 约束内容                                              |
|---------------------|---------------------------------------------------|
| `codeflow-core`     | 始终中文；新需求判断简单/复杂；代码改动必须 TDD 和 Review               |
| `codeflow-git`      | 未经用户确认禁止任何 Git 写操作                                |
| `codeflow-review`   | 有代码逻辑改动必须 Review                                  |
| `codeflow-context`  | 统一读 `.codeflow/state.md`；长任务用 `/codeflow-handoff` |
| `codeflow-openspec` | 复杂需求必须创建 change-id；archive 前 Review 无 must-fix    |

### 工作流定义

| Rule               | 适用场景                 | 流程                                                                    |
|--------------------|----------------------|-----------------------------------------------------------------------|
| `codeflow-simple`  | 需求明确、改动小、不涉及核心链路     | brainstorming → discovery → plans → TDD → Review                      |
| `codeflow-complex` | 多模块、核心链路、接口/DB/权限/安全 | brainstorming → discovery → OpenSpec → plans → TDD → Review → archive |

### 集成约束

| Rule                            | 约束内容                                  |
|---------------------------------|---------------------------------------|
| `codeflow-superpowers`          | 绑定 Superpowers 标准流程到 CodeFlow 工作流     |
| `codeflow-graphify`             | Graphify 可选，不自动执行索引/更新，Hooks 拦截写入     |
| `codeflow-auto-skill-routing`   | 根据技术栈自动加载工程 Skill                     |
| `codeflow-collaborative-agents` | 多代理通用规则：简单需求禁用、写代码前确认、避免同文件冲突         |
| `codeflow-agent-teams`          | Agent Teams 专项：默认不启用、需环境变量、token 成本提醒 |
