---
description: Agent Team 疑难调查：竞争假设方式，每个 teammate 验证一个假设并反驳其他假设
argument-hint: <问题描述>
---

# /codeflow-team-investigate

创建 Claude Code Agent Team 进行疑难问题调查，使用"竞争假设"方式：每个 teammate 负责验证一个假设并尝试反驳其他假设。

```text
/codeflow-team-investigate <问题描述>
```

适用场景：难以复现的 Bug、生产环境诡异问题、需要从多角度排除根因。前置条件：仅 Claude Code 支持；需设置 CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1。

要求：
1. 使用 competing hypotheses。
2. 每个 teammate 负责一个假设。
3. 每个 teammate 尝试反驳其他假设。
4. 不自动修改代码。
5. Team lead 输出最终共识、证据、风险和下一步。