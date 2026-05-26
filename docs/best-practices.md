# CodeFlow 建议指南

什么场景用什么工作流、执行方式怎么选、多代理什么时候用。

---

## 工作流选择

### 决策树

```text
接到需求
  ├─ 改动小、不涉及核心链路/接口/DB/权限/安全？
  │   └─ 是 → 简单流程
  │       流程：brainstorming（轻量）→ discovery（轻量）→ plans（轻量）
  │             → 询问执行方式 → TDD → Review
  │       建议执行方式：executing-plans
  │
  └─ 涉及多模块、核心链路、接口/DB/配置/权限/安全？
      └─ 是 → 复杂流程
          流程：brainstorming → discovery → /opsx:propose
                → spec-review → plans → 询问执行方式
                → TDD → Review → verification → finishing → archive
          建议执行方式：subagent-driven-development
```

### 简单需求示例

- 修改按钮文案、调整样式
- 修复单个函数的边界条件 Bug
- 给已有组件增加一个小 prop
- 更新文档/注释

### 复杂需求示例

- 新增用户角色权限系统
- 数据库表结构变更 + API 改造 + 前端适配
- 新增全文搜索功能（索引服务 + API + 前端）
- 第三方服务集成

---

## 执行方式选择

CodeFlow 支持四种执行方式，按场景选择：

| 执行方式 | 适用场景 | 不适合 | Token 成本 |
|---|---|---|---|
| **executing-plans** | 简单需求、小范围修改、顺序执行 | 跨模块、需要并行 | 低 |
| **subagent-driven-development** | 专项任务拆分、探索、测试、triage | 紧密耦合修改 | 中 |
| **parallel-agent-review** | 多角色并行审查（安全/性能/测试/可维护性） | 简单 review | 中高 |
| **collaborative-agent-development** | 跨模块开发、架构调查、竞争假设排查 | 简单需求、单文件 | 高 |

### 选择原则

1. **默认用 executing-plans**：大部分需求走这个就够了
2. **任务可独立拆分时用 subagent-driven-development**：子任务间无共享状态
3. **重要 PR 用 parallel-agent-review**：需要多维度把关
4. **需要多方协作才用 collaborative-agent-development**：token 成本最高

---

## 多代理使用指南

### 什么时候用

- 复杂 PR 需要安全、性能、测试、可维护性多维度并行审查
- 跨模块功能需要 Architecture / Backend / Test / Review 角色分工
- 疑难 Bug 需要多个假设并行验证

### 什么时候不用

- 简单需求 — 杀鸡用牛刀
- 单文件小修复 — 一个人就能看清楚
- 需要改同一文件的多个任务 — 会产生冲突
- 顺序强依赖的任务 — 并行没有收益

### 通用规则

1. 多代理写代码前必须用户确认
2. 每个代理必须有明确角色、范围和输出格式
3. 高风险任务必须先 plan approval
4. 避免多个代理修改同一文件
5. 并行写代码优先使用 git worktree
6. 所有代理完成后必须输出 final report
7. 主代理综合结果，禁止直接拼接各代理输出

### Claude Code Agent Teams

额外规则：
- 仅 Claude Code 支持
- 默认不启用，需 `export CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`
- 创建团队前必须用户确认
- 完成后必须 Team lead cleanup
- 不允许多个 teammate 同时修改同一文件

---

## Code Review 指南

### 何时 Review

**有代码逻辑改动就必须 Review。** 仅限改文案、修注释可以跳过。

### Review 输出

```md
## Review Result

### Must Fix（必须修复）
- 安全问题、可能导致线上故障的逻辑错误

### Should Fix（建议修复）
- 性能优化、代码可读性改进

### Test Gap（测试缺口）
- 未覆盖的边界条件、回归风险

### Risk（风险点）
- 发布风险、兼容性风险、性能风险

### Evidence（证据）
- 截图、日志、测试输出
```

### 审查维度

| 维度 | 检查内容 |
|---|---|
| 安全 | SQL 注入、XSS、权限校验、敏感信息泄露 |
| 性能 | N+1 查询、不必要的大循环、内存泄漏 |
| 测试 | 正常路径、边界条件、错误路径 |
| 可维护性 | 重复代码、过度抽象、命名清晰度 |
| 状态完整 | loading / empty / error / disabled / selected |

---

## 上下文管理指南

### 何时 Handoff

- 任务执行超过 5 分钟
- 上下文即将压缩前
- 需要在另一个会话继续
- 多人协作传递任务

### 如何减少上下文消耗

1. 当前任务只加载相关 skill
2. 简单需求不加载复杂需求的 workflow skill
3. 多代理不超过最小必要数量
4. Review 和执行分开上下文
5. 长任务中途执行 `/codeflow-handoff`

---

## 团队协作指南

### 建议提交到仓库的文件

```
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

### 不建议提交的文件

```
.codeflow/state.md
.codeflow/active-change.md
.codeflow/tmp/
.codeflow/logs/
graphify-out/
```

### 新人入职

1. Clone 项目 → CodeFlow 配置已在仓库中
2. 发送 `/codeflow-new <第一个需求>` 即可上手
3. 不需要额外安装或配置
