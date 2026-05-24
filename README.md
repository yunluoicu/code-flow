# CodeFlow

CodeFlow 是一套面向 **Claude Code** 的 OpenSpec + Superpowers 工作流安装包。

它不替代 OpenSpec，也不替代 Superpowers，而是把两者组织成一套稳定的项目级 AI 开发流程。

## 特性

- 不污染项目根 `CLAUDE.md`
- 安装到 `.claude/codeflow/`
- 支持 Claude Code hooks
- 支持简单需求 / 复杂需求两套流程
- 支持 OpenSpec change 生命周期
- 支持 Superpowers 官方执行链路
- 支持状态文件，避免上下文丢失
- 始终使用简体中文回复
- 所有 Git 写操作必须用户确认

## 第一版范围

当前只支持 Claude Code。

暂不支持：Codex、Cursor、MCP Server、Dashboard、数据库状态管理。

## 安装

推荐把 `install.md` 发给 Claude Code，让它按文档安装。

也可以使用脚本：

```bash
python3 scripts/install_claude_code.py --target /path/to/target-project
```

## 安装后的项目结构

```text
目标项目/
├── CLAUDE.md
└── .claude/
    ├── settings.json
    └── codeflow/
        ├── CLAUDE.md
        ├── workflows/
        ├── hooks/
        └── state/
```

项目根 `CLAUDE.md` 只追加：

```md
<!-- CodeFlow start -->
@.claude/codeflow/CLAUDE.md
<!-- CodeFlow end -->
```

## 核心流程

### 简单需求

```text
brainstorming 轻量确认
→ Existing Capability Discovery 轻量检查
→ writing-plans 轻量计划
→ 询问用户选择 executing-plans / subagent-driven-development
→ 按用户选择执行
→ test-driven-development 默认必须执行
→ /review 默认必须执行
→ 询问是否 requesting-code-review
→ 询问是否 verification-before-completion
→ 询问是否 finishing-a-development-branch
```

### 复杂需求

```text
brainstorming
→ Existing Capability Discovery
→ /opsx:propose <change-id>
→ Spec Review
→ writing-plans
→ 询问用户选择 executing-plans / subagent-driven-development
→ 按用户选择执行
→ test-driven-development 默认必须执行
→ /review 默认必须执行
→ 询问是否 requesting-code-review
→ 询问是否 verification-before-completion
→ 询问是否 finishing-a-development-branch
→ 用户确认后 /opsx:archive <change-id>
```

## Git 操作规则

未经用户明确确认，禁止执行任何会改变 Git 状态的命令，包括但不限于 `git add`、`git commit`、`git push`、`git merge`、`git rebase`、`git reset`、`git clean`、`git restore`、`git checkout`、`git switch`、`git stash`、`git tag`。
