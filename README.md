# CodeFlow

CodeFlow 是一套面向 AI Coding 工具的 OpenSpec + Superpowers 标准开发工作流包。

当前版本：`CodeFlow 1.2 Skill System Workflow Pack`

支持：

- Claude Code
- Codex
- Cursor

## 设计目标

CodeFlow 不替代 OpenSpec，也不替代 Superpowers。它负责把以下能力组织成稳定可执行的项目级 AI 工作流：

```text
需求澄清 → 已有能力发现 → OpenSpec change → 计划拆解 → 执行 → TDD → Review → 验证 → 收尾 → 归档
```

## 多工具能力对照

| 能力 | Claude Code | Codex | Cursor |
|---|---|---|---|
| 核心规则 | `.claude/rules/` | `AGENTS.md` | `.cursor/rules/*.mdc` |
| 命令入口 | `.claude/commands/` | 自然语言触发 | Agent Chat 触发 |
| Agents | `.claude/agents/` | 不安装 | 不安装 |
| Skills | `.claude/skills/` | `.agents/skills/` | 转为 rules / workflows |
| Hooks | `.claude/settings.json` | 不支持 | 不支持 |
| 状态文件 | `.codeflow/` | `.codeflow/` | `.codeflow/` |
| Graphify | 可选 | 可选 | 可选 |

## 安装方式：手动安装 / AI 代装

### 手动安装

```bash
git clone https://github.com/yunluoicu/code-flow /tmp/code-flow
cd /path/to/your-project
python3 /tmp/code-flow/scripts/install_all.py --target . --tools claude
```

安装 Codex：

```bash
python3 /tmp/code-flow/scripts/install_all.py --target . --tools codex
```

安装 Cursor：

```bash
python3 /tmp/code-flow/scripts/install_all.py --target . --tools cursor
```

三端都安装：

```bash
python3 /tmp/code-flow/scripts/install_all.py --target . --tools claude,codex,cursor
```

预览安装，不实际写入：

```bash
python3 /tmp/code-flow/scripts/install_all.py --target . --tools claude,codex,cursor --dry-run
```

升级：

```bash
python3 /tmp/code-flow/scripts/install_all.py --target . --tools claude,codex,cursor --upgrade
```

强制覆盖 CodeFlow 管理文件：

```bash
python3 /tmp/code-flow/scripts/install_all.py --target . --tools claude,codex,cursor --force
```

### AI 代装

打开目标项目的 Claude Code / Codex / Cursor，然后发送：

```text
请安装 CodeFlow 到当前项目。

仓库：
https://github.com/yunluoicu/code-flow

安装说明：
https://raw.githubusercontent.com/yunluoicu/code-flow/main/install.md

我要安装的工具：
Claude Code / Codex / Cursor

请先阅读安装说明，然后列出：
1. 准备执行的步骤
2. 会创建哪些文件
3. 会修改哪些文件
4. 是否会修改 CLAUDE.md / AGENTS.md / .claude/settings.json / .cursor/rules
5. 是否会执行 git / python / chmod 等命令
6. 风险点

等我确认后再执行安装。
```

一句话：手动安装是你开车；AI 代装是 AI 开车，你坐副驾踩刹车。

## 使用方式

Claude Code：使用 `/codeflow-new`、`/codeflow-status`、`/codeflow-review` 等命令。

Codex：直接说 `按 CodeFlow 开始新需求：<需求内容>`。

Cursor：在 Agent Chat 里说 `按 CodeFlow 开始新需求：<需求内容>`。

## Graphify 可选增强

CodeFlow 不自动安装 Graphify，也不自动执行 `/graphify .` 或 `/graphify . --update`。

如果没有 `graphify-out/graph.json`，会提醒用户可执行：

```text
/graphify .
```

如果图谱可能落后，会提醒用户可执行：

```text
/graphify . --update
```

Graphify 只作为项目理解线索，最终必须以真实代码、OpenSpec specs 和测试文件为准。

## Git 安全规则

所有 Git 写操作必须用户确认，包括但不限于：`git add`、`git commit`、`git push`、`git merge`、`git rebase`、`git reset`、`git clean`、`git restore`、`git checkout`、`git switch`、`git stash`、`git tag`。


## Collaborative Agents

CodeFlow 1.2 继续包含跨工具 Collaborative Agents 能力：

- Claude Code：Agent Teams
- Codex：Subagent Workflows
- Cursor：Parallel Agents / Rules

详见 `docs/collaborative-agents.md`。


## Skill System

CodeFlow 1.2 新增 Engineering Skills 与 Skill System：Vue + TypeScript、Go、Context Budget、Quality Gates、Skill Learning、Eval Checkpoints、Auto Skill Routing。详见 `docs/skill-system.md`。
