---
description: Agent Team 协作开发复杂功能：拆分 Architecture / Backend / Test / Review 角色并行
argument-hint: <复杂功能描述>
---

# /codeflow-team-feature

创建 Claude Code Agent Team 处理复杂功能开发，拆分 Architecture / Backend / Test / Review 角色并行协作。

```text
/codeflow-team-feature <复杂功能描述>
```

适用场景：跨模块大功能开发（如"新增全文搜索+索引服务+前端搜索框"）。前置条件：仅 Claude Code 支持；需设置 CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1。

要求：
1. 先要求所有 teammates 输出 plan。
2. 高风险任务必须 plan approval。
3. 每个 teammate 声明负责文件范围。
4. 避免多个 teammate 修改同一文件。
5. 未经用户确认不要写代码。
