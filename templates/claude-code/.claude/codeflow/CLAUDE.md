# CodeFlow for Claude Code

CodeFlow Claude Code Adapter 入口。安装后在项目 CLAUDE.md 中通过 `@.claude/codeflow/CLAUDE.md` 引入。

始终使用简体中文回复。读取 `.codeflow/state.md`、`.codeflow/active-change.md` 和 `.codeflow/workflows/`。

## Slash Commands

| 命令 | 作用 |
|---|---|
| /codeflow-new | 启动新需求工作流 |
| /codeflow-continue | 继续上次未完成的需求 |
| /codeflow-status | 查看当前需求进度 |
| /codeflow-review | 审查当前分支改动 |
| /codeflow-finish | 完成需求，收尾归档 |
| /codeflow-handoff | 保存上下文摘要供后续恢复 |
| /codeflow-team-review | Agent Team 并行审查 |
| /codeflow-team-feature | Agent Team 协作开发 |
| /codeflow-team-investigate | Agent Team 疑难调查 |
| /codeflow-team-cleanup | 清理 Agent Team |

所有 Git 写操作必须用户确认。
