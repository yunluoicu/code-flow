# CodeFlow 2.0 团队使用指南

## 安装

```bash
go install github.com/yunluoicu/code-flow/cmd/codeflow@latest
```

## 初始化项目

```bash
codeflow init --tools claude,codex,cursor
```

## 查看 Dashboard

```bash
codeflow web --port 4399
```

## AI 工具使用

### Claude Code

使用 `/codeflow-new`、`/codeflow-status`、`/codeflow-review`。

### Codex / Cursor

直接输入：

```text
按 CodeFlow 开始新需求：<需求内容>
```
