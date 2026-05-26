---
name: codeflow-handoff
description: 上下文接力：生成紧凑摘要保存到 state，供下次会话或另一 agent 快速恢复。长任务或跨会话场景使用。
---

# CodeFlow Handoff

长任务上下文接力。生成紧凑的上下文摘要保存到 .codeflow/state.md，供下一次会话（/codeflow-continue）或另一个 agent 快速恢复。

**使用场景**：长任务中途切换、上下文即将压缩前、多人协作传递任务。

**工作流阶段**：TDD → Review 之后，verification 或 finishing 之前（可选）。

始终使用简体中文回复。

本 Skill 只提供流程、检查项和输出规范，不自动执行 Git 写操作。
