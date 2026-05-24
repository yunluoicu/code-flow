# CodeFlow 2.0 发布说明

## 发布范围

- Go CLI
- Claude Code / Codex / Cursor adapter
- .codeflow 项目工作区
- 项目画像与索引
- 需求 / 迭代管理
- OpenSpec / Superpowers 导入
- SQLite + FTS5 索引
- Web Dashboard

## 构建

```bash
go mod tidy
go build -o codeflow ./cmd/codeflow
```

## 本地验证

```bash
codeflow version
codeflow init --tools claude,codex,cursor
codeflow doctor
codeflow profile
codeflow index
codeflow sync
codeflow web --port 4399
```
