# CodeFlow FAQ

## 安装

### Q: 安装会覆盖我的现有配置吗？

默认不会。CodeFlow 只追加 `<!-- CodeFlow start -->` / `<!-- CodeFlow end -->` 标记块到 `CLAUDE.md` 和 `AGENTS.md`。对其他文件使用 `--dry-run` 预览后再安装。如需强制覆盖 CodeFlow 管理文件，使用 `--force`。

### Q: 我只用 Claude Code，需要安装其他 adapter 吗？

不需要。安装时指定 `--tools claude` 即可。

### Q: 如何升级 CodeFlow？

```bash
python3 scripts/install_all.py --target . --tools claude --upgrade
```

---

## 使用

### Q: 简单需求和复杂需求有什么区别？

| | 简单需求 | 复杂需求 |
|---|---|---|
| 改动范围 | 小 | 大，跨模块 |
| 涉及链路 | 非核心 | 核心业务链路 |
| 接口/DB/权限 | 不涉及 | 涉及 |
| OpenSpec | 不创建 | 必须创建 change-id |
| 执行方式 | executing-plans | subagent-driven-development 或 Agent Teams |
| Archive | 可选 | 必须 archive |

### Q: 什么时候用 Agent Teams？

仅在以下场景使用 Agent Teams：

- 复杂 PR 需要多维度并行审查
- 跨模块大功能开发
- 疑难问题需要竞争假设排查

**不要**在简单需求、单文件修复、顺序强依赖任务中使用 Agent Teams。

### Q: Agent Teams 消耗多少 token？

比普通流程多 2-5 倍。创建团队前 CodeFlow 会提醒并等待确认。

### Q: /codeflow-continue 恢复不了状态怎么办？

检查 `.codeflow/state.md` 是否存在。如果文件被清理或状态丢失，直接用 `/codeflow-new` 重新开始并描述当前进度即可。

### Q: Graphify 是必须的吗？

不是。Graphify 是可选的增强能力，辅助 Existing Capability Discovery。CodeFlow 不会自动执行 Graphify 索引或更新。

---

## 安全

### Q: CodeFlow 会自动提交代码吗？

**绝对不会。** Git 写操作（add/commit/push/merge/rebase/reset）必须用户明确确认，Hooks 层会拦截。

### Q: CodeFlow 会修改我的业务代码吗？

安装过程只添加配置文件（rules/commands/agents/skills）和 `.codeflow/` 状态目录，不修改任何业务代码。

### Q: 敏感文件会被修改吗？

Hooks 层会阻止对 `.env`、`.pem`、`.key`、`id_rsa`、`secrets`、`credentials` 等文件的修改操作。

---

## 多工具

### Q: Claude Code、Codex、Cursor 的 CodeFlow 能力一样吗？

不完全一样。核心规则和 Skills 三端通用。Commands 和 Hooks 仅 Claude Code 支持，Agent Teams 仅 Claude Code 支持。

详细对比见 [适配器说明](adapters.md)。

### Q: 如何在 Codex 中使用 CodeFlow？

Codex 不支持 slash commands，通过自然语言触发：

```text
按 CodeFlow 开始新需求：<需求描述>
按 CodeFlow 使用 subagents 审查当前分支
```

### Q: 如何在 Cursor 中使用 CodeFlow？

Cursor 通过 `.cursor/rules/*.mdc` 自动加载 CodeFlow 规则，在 Agent Chat 中自然语言触发：

```text
按 CodeFlow 开始新需求：<需求描述>
按 CodeFlow 并行 Agent 模式审查当前变更
```
