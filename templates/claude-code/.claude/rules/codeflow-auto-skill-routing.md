# CodeFlow Auto Skill Routing

自动 Skill 路由规则。根据当前任务涉及的技术栈，自动加载对应的工程 Skill，无需手动触发。

## Claude Code

当任务涉及 Vue / TypeScript / Vite / Nuxt / Web UI 时，必须优先使用：

- `codeflow-frontend-vue-ts`

当任务涉及 Go / go.mod / 后端服务 / API / DB / Worker / Cron 时，必须优先使用：

- `codeflow-backend-go`

如果同一需求同时涉及前端和后端，必须同时使用两个 skill，并在计划中拆分：

1. Frontend impact
2. Backend impact
3. API contract
4. Test plan

## Skills-first

- Skills 是核心能力单元。
- Commands 只是快捷入口。
- Rules 是硬约束。
- Hooks 做自动化、安全门、状态保存。
- Agents 是专项执行者。
