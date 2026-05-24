<!-- CodeFlow start -->
## CodeFlow

本项目使用 CodeFlow 工作流。始终使用简体中文回复。新需求先判断简单/复杂；复杂需求先做 Existing Capability Discovery；代码逻辑改动默认 TDD；有代码改动必须 Review；所有 Git 写操作必须用户确认；Graphify 可选，不自动执行。

详细规则见 `.codeflow/workflows/`，Skills 位于 `.agents/skills/`。
<!-- CodeFlow end -->


## CodeFlow Codex Subagent Workflows

Codex 可在用户明确要求时使用 subagent workflows。复杂 review / investigation / test triage 可使用 one agent per point。并行写代码前必须用户确认，主线程必须等待所有 subagents 后综合结论。


## CodeFlow Auto Skill Routing

When working on Vue / TypeScript / Vite / Nuxt / Web UI files, use `.agents/skills/codeflow-frontend-vue-ts/SKILL.md`.

When working on Go backend files, go.mod, APIs, services, repositories, jobs, workers, or database models, use `.agents/skills/codeflow-backend-go/SKILL.md`.

For full-stack changes, use both skills and explicitly separate frontend impact, backend impact, API contract, and test plan.
