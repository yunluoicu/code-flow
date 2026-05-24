# CodeFlow Skill System

CodeFlow Skill System 借鉴 Claude Code 生态的 skills-first 设计方式：

- Skills：核心能力单元。
- Commands：触发 skill 的快捷入口。
- Rules：始终遵守的硬约束。
- Hooks：自动化、安全门、状态保存。
- Agents：专项执行者。

## Workflow Skills

- codeflow-openspec
- codeflow-discovery
- codeflow-tdd
- codeflow-review
- codeflow-graphify
- codeflow-handoff
- codeflow-collaborative-agents
- codeflow-context-budget
- codeflow-quality-gates
- codeflow-skill-learning
- codeflow-eval-checkpoints

## Engineering Skills

- codeflow-frontend-vue-ts
- codeflow-backend-go

## 自动路由

- Vue / TS / Vite / Nuxt / UI → codeflow-frontend-vue-ts
- Go / go.mod / API / service / repository / worker → codeflow-backend-go
- 长任务 / compact → codeflow-context-budget + handoff
- Review / 发布前 → codeflow-quality-gates + review

## 上下文预算

不要把所有 skill 默认 alwaysApply。只让核心安全、Git、上下文规则常驻；工程 skill 通过项目类型、文件类型和任务内容触发。
