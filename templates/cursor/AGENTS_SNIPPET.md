<!-- CodeFlow start -->
## CodeFlow

本项目使用 CodeFlow 工作流。Cursor 主要通过 `.cursor/rules/*.mdc` 获取规则。始终中文，Git 写操作必须确认，Graphify 可选不自动执行。
<!-- CodeFlow end -->


## CodeFlow Cursor Parallel Agents

Cursor 可通过 parallel agents / subagents / cloud agents 进行复杂任务协作。CodeFlow 通过 `.cursor/rules/codeflow-collaborative-agents.mdc` 约束使用方式，不管理 Cursor runtime。
