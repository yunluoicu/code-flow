# codeflow-context

上下文管理约束。所有 CodeFlow 流程统一从 .codeflow/state.md 和 .codeflow/active-change.md 读写状态。长任务或跨会话场景必须使用 /codeflow-handoff 保存上下文摘要。
