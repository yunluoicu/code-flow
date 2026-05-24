# 卸载 CodeFlow

删除 `.codeflow/`，移除 `CLAUDE.md` / `AGENTS.md` 的 CodeFlow block，删除 `.claude/codeflow`、`.claude/rules/codeflow-*`、`.claude/commands/codeflow-*`、`.claude/agents/codeflow-*`、`.claude/skills/codeflow-*`、`.agents/skills/codeflow-*`、`.cursor/rules/codeflow-*`，并从 `.claude/settings.json` 移除 codeflow_hooks.py 相关 hooks。
