---
name: codeflow-graphify
description: Graphify 知识图谱集成：可选增强能力，不自动执行索引/更新。图谱缺失时提醒 /graphify .，代码更新后提醒 /graphify . --update。
---

# CodeFlow Graphify

Graphify 知识图谱可选增强。将项目代码索引为知识图谱，辅助 Existing Capability Discovery 和代码理解。

**使用场景**：项目知识图谱查询（可选），不作为必须步骤。

**重要**：Graphify 是可选增强，禁止未经用户确认自动执行 /graphify . 或 /graphify . --update。Hooks 层会拦截写入操作。

始终使用简体中文回复。

本 Skill 只提供流程、检查项和输出规范，不自动执行 Git 写操作。
