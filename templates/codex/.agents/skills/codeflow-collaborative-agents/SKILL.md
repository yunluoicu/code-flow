
---
name: codeflow-codex-subagent-workflows
description: 用于 Codex 中通过 subagent workflows 进行并行探索、Review、测试、triage 和复杂任务拆解。
---

# CodeFlow Codex Subagent Workflows

Codex 使用 Subagent Workflows，不是 Claude Agent Teams。

## 使用规则

1. 只在用户明确要求时 spawn subagents。
2. 每个 subagent 必须有角色、范围、输出格式。
3. 默认用于 read-heavy 分析、测试、triage、review。
4. 并行写代码前必须获得用户确认。
5. 主线程必须等待所有 subagents 结果后再总结。
6. 并行写代码要谨慎，优先 worktree。
7. 输出必须包含 Must Fix / Should Fix / Test Gap / Risk / Evidence。

## 推荐 Review Prompt

请按 CodeFlow 使用 Codex subagents 审查当前分支。

Spawn one subagent per point:
1. Security
2. Code quality
3. Bugs
4. Race
5. Test flakiness
6. Maintainability

Wait for all agents, then summarize findings by:
- Must Fix
- Should Fix
- Test Gap
- Risk
- Evidence
