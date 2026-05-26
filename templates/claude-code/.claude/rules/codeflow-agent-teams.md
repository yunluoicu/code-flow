
# CodeFlow Agent Teams

Claude Code Agent Teams 专项约束。Agent Teams 仅适用于 Claude Code，默认不启用（需设置 CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1），token 成本更高，创建前必须用户确认。

Agent Teams 仅适用于 Claude Code。

## 默认策略

- 默认不启用 Agent Teams
- 默认只安装规则和命令模板
- 用户明确开启后才使用 Agent Teams
- 创建团队前必须用户确认
- token 成本更高，创建前必须提醒用户

## 禁止事项

- 不生成项目级 `.claude/teams/`
- 不让 teammate 自己 cleanup
- 不允许多个 teammate 修改同一文件
- 不允许简单需求使用 Agent Teams
- 不允许未 plan approval 就实施高风险任务

## 适合场景

- complex-review
- cross-layer-feature
- competing-hypothesis-debug
- architecture-investigation

## 不适合场景

- simple-change
- same-file-editing
- sequential-task
- tight-coupled-editing

## 执行方式

writing-plans 后询问：

1. executing-plans
2. subagent-driven-development
3. agent-team-driven-development

用户确认前不得创建团队。

## Cleanup

完成后必须要求所有 teammates 汇报最终状态。
确认无活跃 teammate 后，由 Team lead 执行 cleanup。
