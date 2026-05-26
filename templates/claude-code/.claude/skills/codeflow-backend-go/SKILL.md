---
name: codeflow-backend-go
description: Go 后端工程规范：类型清晰、错误可追踪、context 贯穿、接口小而明确、结构化日志、table-driven tests。触发条件：*.go / go.mod / cmd/ / internal/ / api/ / service/ 等。
---

# CodeFlow Backend Go Skill

## 触发条件

当任务涉及以下内容时必须使用本技能：

- `go.mod`
- `*.go`
- `cmd/`
- `internal/`
- `pkg/`
- `api/`
- `service/`
- `repository/`
- `model/`
- `middleware/`
- `job/`
- `worker/`
- `cron/`
- Go API / 后端服务 / 数据库 / 队列 / 定时任务 / 监控 / 测试

## 核心目标

输出生产级、可维护、符合 Go idiom 的后端代码。

要求：

1. 类型清晰。
2. 错误可追踪。
3. context 贯穿 IO 边界。
4. 接口小而明确。
5. 依赖方向清晰。
6. 日志结构化。
7. 测试可执行。
8. 不过度抽象。
9. 不写 Java/PHP 味 Go。

## 使用方式

开始前先检查：

- `references/go-style.md`
- `references/error-handling.md`
- `references/context.md`
- `references/interface-design.md`
- `references/concurrency.md`
- `references/api-service-repository.md`
- `references/security.md`
- `references/testing.md`
- `checklists/go-review.md`

## 输出要求

Go 实现完成后必须说明：

- 修改的 package / API / service / repository
- 错误处理方式
- context 和并发退出方式
- 日志字段
- 测试或替代验证方式
- 发布/迁移/配置影响
