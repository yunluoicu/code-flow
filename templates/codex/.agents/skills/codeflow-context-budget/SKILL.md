---
name: codeflow-context-budget
description: Use to control context size, active skills, MCP/tools, agents, and long-session handoff.
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
