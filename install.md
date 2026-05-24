# CodeFlow Claude Code 安装说明

## 目标

将 CodeFlow 工作流安装到当前项目，支持 Claude Code。

安装后：

- 不覆盖项目根 `CLAUDE.md`
- 不破坏原有 `.claude/settings.json`
- 在 `.claude/codeflow/` 下安装流程规则、hooks、状态文件
- 支持简单需求和复杂需求两套工作流
- 始终使用简体中文回复
- 所有 Git 写操作必须用户确认

## 安装前检查

请先检查当前项目：

1. 是否存在 `CLAUDE.md`
2. 是否存在 `.claude/`
3. 是否存在 `.claude/settings.json`
4. 是否存在 `openspec/`
5. 是否已安装 OpenSpec / OPSX
6. 是否已安装 Superpowers

要求：缺少 `openspec/` 或 Superpowers 时只提示用户，不要自动初始化或安装；不要修改业务代码；不要创建 OpenSpec change。

## 推荐安装方式

```bash
python3 scripts/install_claude_code.py --target /path/to/target-project
```

如果已在目标项目根目录执行：

```bash
python3 /path/to/code-flow/scripts/install_claude_code.py --target .
```

## 手动安装步骤

1. 创建目录：

```bash
mkdir -p .claude/codeflow/workflows .claude/codeflow/hooks .claude/codeflow/state
```

2. 复制模板：

```text
templates/claude-code/.claude/codeflow/ → .claude/codeflow/
```

3. 更新项目根 `CLAUDE.md`：

```md
<!-- CodeFlow start -->
@.claude/codeflow/CLAUDE.md
<!-- CodeFlow end -->
```

要求：不覆盖原有 `CLAUDE.md`；已存在 CodeFlow 引用时不要重复添加。

4. 合并 `.claude/settings.json`：只合并 CodeFlow hooks，不删除用户已有配置。参考 `templates/claude-code/.claude/settings.codeflow.example.json`。

5. 初始化状态文件：若已存在则备份，不直接覆盖。

6. 设置 hook 执行权限：

```bash
chmod +x .claude/codeflow/hooks/*.sh
```

## 幂等安装要求

- 已存在 CodeFlow 引用，不重复添加
- 已存在 CodeFlow hooks，不重复添加
- 已存在状态文件，不直接覆盖
- 已存在 `.claude/codeflow/CLAUDE.md`，如内容不同必须先备份
- 不覆盖用户已有 `.claude/settings.json`

## 禁止事项

安装过程中禁止覆盖项目根 `CLAUDE.md`、删除已有 `.claude/settings.json`、修改业务代码、自动安装 OpenSpec、自动安装 Superpowers、创建 `openspec/changes/<change-id>`、执行 Git 写操作。
