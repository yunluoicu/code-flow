
# CodeFlow Collaborative Agents

CodeFlow 多代理协作统一规范。定义 Claude Code / Codex / Cursor 三端的多代理协作通用规则、四种执行方式及推荐 Prompt 模板。简单需求禁止使用多代理，多代理写代码前必须用户确认。

CodeFlow Collaborative Agents 是 CodeFlow 对多 AI 工具协作代理能力的统一抽象。

它不是单一工具能力，而是三端适配：

| 工具 | CodeFlow 名称 | 原生能力 |
|---|---|---|
| Claude Code | Agent Teams | Team lead + Teammates + shared task list + mailbox |
| Codex | Subagent Workflows | 并行 specialized subagents + agent threads |
| Cursor | Parallel Agents | Subagents / Cloud Agents / worktrees / rules |

## 统一执行方式

CodeFlow 中复杂任务支持四种执行方式：

1. `executing-plans`
   - 当前会话按计划执行
   - 适合简单、小范围、顺序任务

2. `subagent-driven-development`
   - 拆给子代理做专项任务，只需要结果回报
   - 适合 bounded work、探索、测试、triage

3. `parallel-agent-review`
   - 多个代理并行审查同一变更的不同维度
   - 适合安全、性能、测试、维护性并行 Review

4. `collaborative-agent-development`
   - 多代理协作开发复杂功能或调查问题
   - 适合跨模块开发、架构评审、竞争假设排查

## 通用规则

1. 简单需求不要使用多代理。
2. 多代理只用于复杂 Review、跨模块开发、架构评审、疑难问题调查。
3. 多代理默认先做研究、审查、方案和测试建议。
4. 多代理写代码前必须用户确认。
5. 写代码前必须说明每个代理的文件范围。
6. 避免多个代理修改同一文件。
7. 并行写代码优先使用 git worktree。
8. 高风险任务必须先 plan approval。
9. 所有代理完成后必须输出 final report。
10. 主代理必须综合结果，不允许直接拼接各代理输出。
11. 多代理 token / 成本更高，创建前必须提醒用户。
12. 不自动执行 Git 写操作。

## Claude Code：Agent Teams

Claude Code 的 Agent Teams 是实验性功能，默认关闭。

规则：

1. 仅 Claude Code 支持 Agent Teams。
2. 默认不启用，只安装规则和命令模板。
3. 启用前需要检查 Claude Code 版本和环境变量。
4. 不生成项目级 `.claude/teams/`。
5. 创建 Agent Team 前必须用户确认。
6. 高风险任务必须 require plan approval。
7. 完成后必须由 Team lead cleanup。
8. 不允许多个 teammate 同时修改同一文件。
9. Agent Teams 不替代 subagents。

建议启用方式：

```bash
export CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
```

## Codex：Subagent Workflows

Codex 使用 Subagent Workflows，不叫 Agent Teams。

规则：

1. Codex 只在用户明确要求时 spawn subagents。
2. 每个 subagent 必须有清晰角色、范围和输出格式。
3. 默认用于 read-heavy 分析、测试、triage、review。
4. 并行写代码前必须获得用户确认。
5. 主线程必须等待所有 subagents 结果后再总结。
6. 输出必须包含 Must Fix / Should Fix / Test Gap / Risk / Evidence。

## Cursor：Parallel Agents

Cursor 使用 Parallel Agents / Subagents / Cloud Agents / Rules。

规则：

1. Cursor 不使用 Claude Agent Teams runtime。
2. 通过 `.cursor/rules/*.mdc` 和 prompt 模板约束协作行为。
3. 复杂任务可以使用 parallel agents。
4. 每个 agent 必须有明确角色、范围、输出格式。
5. 并行写代码前必须用户确认。
6. 推荐使用 worktrees 避免冲突。
7. 主 Agent 必须综合结果。

## 推荐 Prompt：并行 Review

```text
按 CodeFlow Collaborative Agents 模式审查当前变更。

请并行拆分以下审查角色：
1. Security Reviewer：安全、权限、敏感配置、危险操作
2. Performance Reviewer：性能、并发、资源使用
3. Test Reviewer：测试覆盖、边界场景、回归风险
4. Maintainability Reviewer：架构一致性、可维护性、重复实现

每个代理只输出：
- Must Fix
- Should Fix
- Test Gap
- Risk
- Evidence

主代理必须等待全部结果后，合成最终 Review Result。
不要自动修改代码。
```

## 推荐 Prompt：疑难问题调查

```text
按 CodeFlow Collaborative Agents 模式调查这个问题。

请使用竞争假设方式拆分多个代理：
1. Hypothesis A
2. Hypothesis B
3. Hypothesis C
4. Test / Evidence Checker

每个代理需要尝试证明自己的假设，也要指出其他假设可能不成立的证据。
最后由主代理输出最终共识、证据、风险和下一步。
不要自动修改代码。
```

## 推荐 Prompt：复杂功能开发

```text
按 CodeFlow Collaborative Agents 模式处理这个复杂功能。

请先拆分角色：
1. Architecture / Design
2. Backend / API
3. Test / QA
4. Review / Risk

要求：
1. 所有代理先输出 plan。
2. 高风险实现必须 plan approval。
3. 每个代理说明负责文件范围。
4. 避免多个代理修改同一文件。
5. 未经用户确认不要写代码。
```
