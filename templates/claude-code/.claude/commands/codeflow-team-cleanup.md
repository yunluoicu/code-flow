---
description: 清理 Agent Team：要求所有 teammates 汇报状态，确认无活跃 teammate 后由 Team lead 清理
---

# /codeflow-team-cleanup

清理 Claude Code Agent Team，要求所有 teammates 汇报最终状态，确认无活跃 teammate 后由 Team lead 执行 cleanup 并输出最终总结。

```text
/codeflow-team-cleanup
```

前置条件：仅 Claude Code 支持；需设置 CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1。Agent Team 协作完成后必须执行，避免遗留活跃 agent。

要求：
1. 要求所有 teammates 汇报最终状态。
2. 确认没有活跃 teammate。
3. 只能由 Team lead cleanup。
4. 输出 cleanup 后的最终总结。