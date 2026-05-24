# CodeFlow for Claude Code

## 0. 语言规则

始终使用简体中文回复。除非用户明确要求英文，否则所有解释、计划、审查、错误说明、完成报告都必须使用简体中文。

## 1. 工作流目标

本项目使用 CodeFlow 约束 Claude Code 开发流程。CodeFlow 只组织流程，不替代 OpenSpec，也不替代 Superpowers。

默认使用：OpenSpec、Superpowers、Claude Code hooks。

## 2. 新需求进入规则

收到任何开发类需求后，必须先判断简单需求或复杂需求。禁止未分类直接修改代码。

每次新需求开始前必须先确认：是否已有活跃 change；本需求是新需求还是继续已有需求；是否需要读取 `openspec/specs/`；是否需要 Existing Capability Discovery；是否需要创建 `openspec/changes/<change-id>`；是否需要 Superpowers brainstorming。

如果不确定，先询问用户，或按复杂需求处理。

## 3. 简单需求流程

简单需求不强制创建 OpenSpec change，但仍必须遵守 CodeFlow。

```text
Superpowers brainstorming 轻量确认
→ Existing Capability Discovery 轻量检查
→ Superpowers writing-plans 轻量计划
→ 询问用户选择 executing-plans / subagent-driven-development
→ 按用户选择执行
→ test-driven-development 默认必须执行
→ /review 默认必须执行
→ 询问是否执行 requesting-code-review
→ 询问是否执行 verification-before-completion
→ 询问是否执行 finishing-a-development-branch
```

如果简单需求过程中发现影响范围扩大，必须升级为复杂需求流程。

## 4. 复杂需求流程

复杂需求必须执行：

```text
Superpowers brainstorming
→ Existing Capability Discovery
→ /opsx:propose <change-id>
→ Spec Review
→ Superpowers writing-plans
→ 询问用户选择 executing-plans / subagent-driven-development
→ 按用户选择执行
→ test-driven-development 默认必须执行
→ /review 默认必须执行
→ 询问是否执行 requesting-code-review
→ 询问是否执行 verification-before-completion
→ 询问是否执行 finishing-a-development-branch
→ 用户确认后执行 /opsx:archive <change-id>
```

## 5. 执行方式选择规则

完成 writing-plans 后，不允许直接开始实现。必须询问用户选择 executing-plans 或 subagent-driven-development。用户确认前不得进入实现阶段。

推荐：简单需求默认建议 executing-plans；复杂需求默认建议 subagent-driven-development；最终以用户选择为准。

## 6. subagent 上下文规则

如果用户选择 subagent-driven-development，Claude 必须先为每个子任务准备上下文包。上下文包必须包含 change-id、当前任务目标、相关 OpenSpec change、已确认实施计划、Existing Capability Discovery 结论、当前任务范围、明确不调整内容、测试要求、完成后必须输出的结果。子 agent 如果发现上下文不足，必须返回 `NEEDS_CONTEXT`，不允许自行猜测。

## 7. TDD 默认规则

只要涉及代码逻辑修改，默认必须执行 test-driven-development。执行要求：优先 RED 测试或最小验证用例，再 GREEN 最小实现，最后必要时 REFACTOR。无法测试时必须说明原因，并给出手动验证步骤。

以下情况可以不执行代码级 TDD，但必须提供验证方式：纯文档修改、纯注释修改、纯提示词修改、纯配置说明修改、不影响运行逻辑的格式调整。

## 8. /review 审查规则

无论简单需求还是复杂需求，只要有代码改动，实现完成后都必须执行 `/review` 审查本次需求改动。`/review` 的目标不是重新讨论需求，也不是自动修复代码，而是检查本次实现是否符合需求、计划、OpenSpec 约束和测试要求。如果 `/review` 发现 must-fix，必须先修复，再重新执行测试和 `/review`。

## 9. requesting-code-review 规则

requesting-code-review 是阶段性审查，不自动执行。每完成一个阶段任务后，必须询问用户是否执行 requesting-code-review。用户同意则执行，用户拒绝则继续下一步任务。

## 10. verification-before-completion 规则

verification-before-completion 不自动执行。实现完成后必须询问用户是否执行。用户拒绝时，不得声明“完整验证通过”，只能说明“实现已完成，但未执行完整完成前验证”。如果用户拒绝 verification-before-completion，允许继续后续流程，但 archive 说明中必须记录“未执行完整完成前验证”。

## 11. finishing-a-development-branch 规则

finishing-a-development-branch 不自动执行。准备结束当前需求前必须询问用户是否执行。用户同意则执行，用户拒绝则输出当前完成内容、测试情况、风险和后续建议。

## 12. Archive Gate

执行 `/opsx:archive` 前必须确认：`/review` 无 must-fix；`tasks.md` 勾选状态真实；测试或替代验证已记录；OpenSpec spec delta 与最终实现一致；没有未处理的用户确认问题；没有未说明的风险。

## 13. Git 写操作规则

未经用户明确确认，禁止执行任何会改变 Git 状态的命令，包括但不限于 git add、git commit、git push、git merge、git rebase、git reset、git clean、git restore、git checkout、git switch、git stash、git tag。

执行前必须说明准备执行的命令、为什么执行、会影响哪些文件或分支、是否可回滚、风险是什么。用户明确确认后才能执行。

## 14. 状态管理规则

复杂需求必须维护 `.claude/codeflow/state/state.md` 和 `.claude/codeflow/state/active-change.md`。允许的 phase：idle、brainstorming、discovery、proposing、spec-review、planning、waiting-execution-mode、executing、tdd、reviewing、verification-pending、finishing-pending、archive-ready、archived、blocked。每次开始任务前必须读取状态文件。

## 15. 完成回复规则

完成回复必须包含修改内容、修改文件、测试结果、未验证内容、风险和待确认点、是否执行 `/review`、是否可以进入下一步。没有验证结果，不能说完整完成。
