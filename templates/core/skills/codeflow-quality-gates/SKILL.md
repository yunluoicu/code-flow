---
name: codeflow-quality-gates
description: Use to enforce test, review, security, and verification checkpoints before completion.
---

# CodeFlow Quality Gates Skill

## 通用质量门

1. 有代码逻辑改动必须 TDD 或替代验证。
2. 有代码改动必须 Review。
3. Review 输出必须包含 Must Fix / Should Fix / Test Gap / Risk / Evidence。
4. 不能验证时必须说明原因和替代验证方式。
5. Git 写操作必须用户确认。

## 前端质量门

- loading / empty / error / disabled / selected 状态完整。
- 无明显 UI 反模式。
- API contract 类型完整。

## Go 质量门

- gofmt。
- go test。
- error wrapping。
- context。
- structured logging。
- 并发退出和资源释放。
