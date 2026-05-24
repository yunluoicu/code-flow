
---
name: codeflow-collaborative-agents
description: 用于复杂 Review、疑难问题调查、跨模块开发和架构评审时，指导 Claude Agent Teams、Codex Subagent Workflows、Cursor Parallel Agents 的协作规则。
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
