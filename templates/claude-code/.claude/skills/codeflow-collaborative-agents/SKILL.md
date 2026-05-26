
---
name: codeflow-collaborative-agents
description: 多代理协作规范：定义四种执行方式（executing-plans / subagent-driven-development / parallel-agent-review / collaborative-agent-development），统一 Claude Code / Codex / Cursor 三端协作规则。简单需求禁用，写代码前必须用户确认。
---

# CodeFlow Collaborative Agents Skill

## 什么时候使用

- 复杂 Review
- 跨模块功能
- 架构调查
- 疑难问题排查
- 竞争假设验证
- 大型测试/风险分析

## 什么时候不要使用

- 简单需求
- 单文件小修复
- 顺序强依赖任务
- 多代理会同时修改同一文件的任务

## 执行方式

1. executing-plans
2. subagent-driven-development
3. parallel-agent-review
4. collaborative-agent-development

## 统一规则

- 多代理写代码前必须用户确认
- 每个代理必须有明确角色、范围和输出格式
- 高风险任务必须先 plan approval
- 并行写代码优先使用 git worktree
- 主代理必须综合结果
- 不自动执行 Git 写操作

## 输出格式

```md
## Collaborative Agents Result

### Overall Conclusion

### Agent Findings

### Must Fix

### Should Fix

### Test Gap

### Risk

### Evidence

### Next Step
```
