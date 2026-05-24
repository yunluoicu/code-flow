# CodeFlow

本项目使用 CodeFlow 工作流。

核心规则：

- 始终使用简体中文回复
- 新需求先判断简单 / 复杂
- 复杂需求先执行 Existing Capability Discovery
- 复杂需求需要 OpenSpec change
- 代码逻辑改动默认 TDD
- 有代码改动必须 Review
- 所有 Git 写操作必须用户确认
- Graphify 可选，不自动执行

详细规则：

- `.codeflow/workflows/simple-requirement.md`
- `.codeflow/workflows/complex-requirement.md`
- `.codeflow/workflows/existing-capability-discovery.md`
- `.codeflow/workflows/graphify-discovery.md`
