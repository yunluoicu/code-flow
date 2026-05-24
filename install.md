# CodeFlow 1.0.2 安装说明

## 支持工具

- Claude Code
- Codex
- Cursor

## 推荐命令

```bash
python3 scripts/install_all.py --target . --tools claude,codex,cursor --dry-run
python3 scripts/install_all.py --target . --tools claude
python3 scripts/install_all.py --target . --tools codex
python3 scripts/install_all.py --target . --tools cursor
python3 scripts/install_all.py --target . --tools claude,codex,cursor
python3 scripts/install_all.py --target . --tools claude,codex,cursor --upgrade
python3 scripts/install_all.py --target . --tools claude,codex,cursor --force
```

## AI 代装提示词

```text
请从以下仓库安装 CodeFlow 到当前项目：

https://github.com/yunluoicu/code-flow

请先阅读安装说明：

https://raw.githubusercontent.com/yunluoicu/code-flow/main/install.md

我要安装的工具：Claude Code / Codex / Cursor

安装前请先输出：准备执行的步骤、会创建哪些文件、会修改哪些文件、风险点。
等我确认后再执行安装。
```

## 安装内容

通用 Core：`.codeflow/manifest.json`、`.codeflow/state.md`、`.codeflow/active-change.md`、`.codeflow/workflows/`。

Claude Code：`.claude/codeflow/`、`.claude/rules/`、`.claude/commands/`、`.claude/agents/`、`.claude/skills/`、`.claude/settings.json`、`CLAUDE.md` 引用。

Codex：`AGENTS.md`、`.agents/skills/`。

Cursor：`.cursor/rules/*.mdc`、可选 `AGENTS.md` 摘要。

## Graphify

CodeFlow 不自动安装 Graphify。如果没有 `graphify-out/graph.json`，会提醒 `/graphify .`；如果可能过期，会提醒 `/graphify . --update`。

## 禁止事项

安装过程禁止修改业务代码、自动安装 OpenSpec/Superpowers/Graphify、自动执行 Graphify、自动执行任何 Git 写操作。


## Collaborative Agents 安装说明

默认安装规则、skills、prompts 和 Claude commands，但不默认启用 Claude Agent Teams 实验能力。

Claude Code 如需启用 Agent Teams，请确认 Claude Code 版本和环境变量：

```bash
export CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
```

CodeFlow 不生成项目级 `.claude/teams/`，不管理 Agent Teams runtime。
