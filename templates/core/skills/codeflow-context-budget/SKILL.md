---
name: codeflow-context-budget
description: 上下文预算控制：限制当前任务加载的 skill 数量、MCP 工具、代理数量，避免撑爆 AI 上下文。长任务必须 handoff，compact 前生成摘要。
---

# CodeFlow Context Budget Skill

## 目标

控制上下文成本，避免 CodeFlow 自己把 AI 上下文撑爆。

## 规则

1. 当前任务只加载相关 skill。
2. 简单需求禁止加载过多 workflow skill。
3. 多代理默认不超过最小必要数量。
4. MCP 默认不超过必要数量。
5. Review 和执行分开上下文。
6. 长任务必须 handoff。
7. compact 前必须生成上下文摘要。

## 输出

- Active Skills
- Skipped Skills
- Context Risk
- Suggested Handoff
