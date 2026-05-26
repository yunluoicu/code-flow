# CodeFlow

CodeFlow 是一套面向 AI Coding 工具的 OpenSpec + Superpowers 标准开发工作流包。

支持 Claude Code、Codex、Cursor 三端。

## 一句话概括

```text
CodeFlow = 给项目安装 AI 开发规范、技能和工具适配的轻量规则包。
```

## 核心能力

- **标准工作流**：需求澄清 → 已有能力发现 → OpenSpec change → 计划拆解 → 执行 → TDD → Review → 验证 → 收尾 → 归档
- **10 个 Slash Commands**：`/codeflow-new`、`/codeflow-review`、`/codeflow-finish` 等
- **6 个专项 Agents**：需求分析、代码审查、OpenSpec 审查、能力发现、计划审查、测试审查
- **13 个 Skills**：含前端 Vue+TS 工程规范、Go 后端工程规范、TDD、Review、上下文预算、质量门等
- **Collaborative Agents**：Agent Teams 并行协作（仅 Claude Code）
- **Hooks 安全门**：阻止未经确认的 Git 写操作、危险命令、敏感文件修改
- **自动 Skill 路由**：根据技术栈自动加载对应工程 Skill

## 多工具能力对照

| 能力       | Claude Code                  | Codex             | Cursor                |
|----------|------------------------------|-------------------|-----------------------|
| 核心规则     | `.claude/rules/`             | `AGENTS.md`       | `.cursor/rules/*.mdc` |
| 命令入口     | `/codeflow-*` slash commands | 自然语言触发            | Agent Chat 触发         |
| Agents   | 6 个专项 agent                  | 不安装               | 不安装                   |
| Skills   | `.claude/skills/`            | `.agents/skills/` | 转为 rules / workflows  |
| Hooks    | 安全门 + 状态注入                   | 不支持               | 不支持                   |
| 状态文件     | `.codeflow/`                 | `.codeflow/`      | `.codeflow/`          |
| Graphify | 可选                           | 可选                | 可选                    |

## 快速开始

### 1. 安装

```bash
git clone https://github.com/yunluoicu/code-flow /tmp/code-flow
cd /path/to/your-project

# 预览安装（推荐先执行）
python3 /tmp/code-flow/scripts/install_all.py --target . --tools claude --dry-run

# 安装 Claude Code adapter
python3 /tmp/code-flow/scripts/install_all.py --target . --tools claude

# 或三端全装
python3 /tmp/code-flow/scripts/install_all.py --target . --tools claude,codex,cursor
```

也可让 AI 代装：在目标项目中发送"请安装 CodeFlow 到当前项目"，指向 `https://github.com/yunluoicu/code-flow` 和
`install.md`。

### 2. 开启第一个需求

```
/codeflow-new 用户管理模块增加角色权限控制功能
```

CodeFlow 会自动判断需求复杂度，走对应的简单/复杂工作流。

### 3. 继续之前的需求

```
/codeflow-continue
```

### 4. 审查改动

```
/codeflow-review
```

### 5. 收尾归档

```
/codeflow-finish
```

## 文档导航

| 文档                             | 内容                                |
|--------------------------------|-----------------------------------|
| [README](README.md)            | 项目介绍与快速开始（本文件）                    |
| [使用指南](docs/usage-guide.md)    | Commands、Agents、Skills、Rules 完整详解 |
| [FAQ](docs/faq.md)             | 常见问题                              |
| [建议指南](docs/best-practices.md) | 什么场景用什么工作流、执行方式选择指南               |
| [安装说明](install.md)             | 详细安装与升级说明                         |
| [卸载说明](docs/uninstall.md)      | 卸载步骤                              |
| [适配器说明](docs/adapters.md)      | 三端 adapter 差异与设计                  |
| [设计方案](CodeFlow_1.x完整设计方案.md)  | 完整架构设计文档                          |

## 环境要求（Agent Teams）

Claude Code 的 Agent Teams 是实验性功能，默认不启用。如需使用 `/codeflow-team-*` 命令：

```bash
export CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
```

## 设计原则

```text
1. Skills-first：skills 是核心能力单元
2. Commands 只是快捷入口
3. Rules 是硬约束
4. Hooks 做自动化、安全门和上下文保存
5. Agents 做专项执行者
6. Collaborative Agents 用于复杂并行协作
7. 不自动 Git 写操作
8. 不自动执行 Graphify
9. 简单需求不使用多代理
10. 当前任务只加载相关 skill
```
