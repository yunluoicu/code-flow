package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func InitProject(opt InitOptions) error {
	if len(opt.Tools) == 0 {
		opt.Tools = []string{"claude"}
	}
	if opt.Port == 0 {
		opt.Port = 4399
	}

	fmt.Println("CodeFlow init")
	fmt.Println("tools:", strings.Join(opt.Tools, ","))
	fmt.Println("dry-run:", opt.DryRun, "force:", opt.Force)

	if err := initCodeflowFiles(opt); err != nil {
		return err
	}
	for _, t := range opt.Tools {
		switch t {
		case "claude":
			if err := installClaude(opt); err != nil {
				return err
			}
		case "codex":
			if err := installCodex(opt); err != nil {
				return err
			}
		case "cursor":
			if err := installCursor(opt); err != nil {
				return err
			}
		default:
			return fmt.Errorf("不支持的工具: %s", t)
		}
	}

	if err := Profile(); err != nil {
		return err
	}
	if err := Index(); err != nil {
		return err
	}
	if err := Sync(SyncOptions{}); err != nil {
		return err
	}

	fmt.Println("CodeFlow init 完成")
	fmt.Println(GraphStatusMessage())
	return nil
}

func initCodeflowFiles(opt InitOptions) error {
	base := ".codeflow"
	dirs := []string{
		base, base + "/requirements", base + "/iterations", base + "/decisions",
		base + "/reviews", base + "/checks", base + "/workflows", base + "/tmp", base + "/logs",
	}
	for _, d := range dirs {
		if opt.DryRun {
			fmt.Println("[dry-run] mkdir", d)
		} else if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	manifest := map[string]any{
		"version":     Version,
		"tools":       opt.Tools,
		"installedAt": now(),
		"projectPath": cwd(),
	}
	if err := writeJSON(".codeflow/manifest.json", manifest, opt.DryRun); err != nil {
		return err
	}
	config := fmt.Sprintf("web:\n  host: 127.0.0.1\n  port: %d\n\ngraphify:\n  autoRun: false\n\nsecurity:\n  allowGitWrite: false\n", opt.Port)
	if !exists(".codeflow/config.yaml") || opt.Force {
		if err := writeFile(".codeflow/config.yaml", config, opt.DryRun); err != nil {
			return err
		}
	}
	if !exists(".codeflow/state.md") || opt.Force {
		if err := writeFile(".codeflow/state.md", "# CodeFlow State\n\nphase: idle\nstatus: idle\n", opt.DryRun); err != nil {
			return err
		}
	}
	if !exists(".codeflow/active-change.md") || opt.Force {
		if err := writeFile(".codeflow/active-change.md", "# Active Change\n\n当前没有活跃 change。\n", opt.DryRun); err != nil {
			return err
		}
	}
	workflows := map[string]string{
		"simple-requirement.md":   "brainstorming 轻量确认 → discovery → writing-plans → 询问执行方式 → TDD → Review\n",
		"complex-requirement.md":  "brainstorming → discovery → OpenSpec change → Spec Review → writing-plans → 执行 → TDD → Review → archive\n",
		"graphify-discovery.md":   "Graphify 可选。无 graphify-out/graph.json 提醒 /graphify .；可能过期提醒 /graphify . --update。\n",
		"collaborative-agents.md": "Collaborative Agents：Claude Agent Teams、Codex Subagent Workflows、Cursor Parallel Agents。简单需求不要使用多代理，高风险任务必须 plan approval。\n",
	}
	for name, c := range workflows {
		if err := writeFile(filepath.Join(".codeflow/workflows", name), "# "+name+"\n\n"+c, opt.DryRun); err != nil {
			return err
		}
	}
	return nil
}

func installClaude(opt InitOptions) error {
	fmt.Println("install Claude Code adapter")
	block := "<!-- CodeFlow start -->\n@.claude/codeflow/CLAUDE.md\n<!-- CodeFlow end -->\n"
	if err := appendBlock("CLAUDE.md", "<!-- CodeFlow start -->", "<!-- CodeFlow end -->", block, opt.DryRun); err != nil {
		return err
	}
	files := map[string]string{
		".claude/codeflow/CLAUDE.md":                    "# CodeFlow for Claude Code\n\n始终使用简体中文回复。读取 `.codeflow/` 状态。有代码逻辑改动必须 TDD 和 `/review`。Git 写操作必须用户确认。\n",
		".claude/rules/codeflow-core.md":                "# CodeFlow Core\n\n新需求先判断简单/复杂。TDD 和 /review 必须执行。Git 写操作必须确认。\n",
		".claude/rules/codeflow-git.md":                 "# Git Safety\n\n未经用户确认禁止 git add/commit/push/merge/rebase/reset/clean/restore/checkout/switch/stash/tag。\n",
		".claude/rules/codeflow-graphify.md":            "# Graphify\n\nGraphify 可选，不自动运行。缺失提醒 `/graphify .`，可能过期提醒 `/graphify . --update`。\n",
		".claude/commands/codeflow-new.md":              "# /codeflow-new\n\n开始新需求：判断简单/复杂，执行 brainstorming 和 Existing Capability Discovery。\n",
		".claude/commands/codeflow-status.md":           "# /codeflow-status\n\n读取 `.codeflow/manifest.json`、`.codeflow/state.md`、`.codeflow/active-change.md`。\n",
		".claude/commands/codeflow-review.md":           "# /codeflow-review\n\n审查本次需求改动或当前分支是否达到发布要求。\n",
		".claude/commands/codeflow-handoff.md":          "# /codeflow-handoff\n\n生成上下文交接摘要。\n",
		".claude/rules/codeflow-agent-teams.md":         "# CodeFlow Agent Teams\n\nAgent Teams 仅适用于 Claude Code。默认不启用，不生成项目级 .claude/teams/。创建前必须用户确认，高风险任务必须 plan approval，完成后必须 Team lead cleanup。\n",
		".claude/commands/codeflow-team-review.md":      "# /codeflow-team-review\n\n创建 Claude Code Agent Team 进行并行 Review。不要自动修改代码，输出 Must Fix / Should Fix / Test Gap / Risk / Evidence。\n",
		".claude/commands/codeflow-team-investigate.md": "# /codeflow-team-investigate\n\n创建 Claude Code Agent Team 进行竞争假设式问题调查。不要自动修改代码。\n",
		".claude/commands/codeflow-team-feature.md":     "# /codeflow-team-feature\n\n创建 Claude Code Agent Team 处理复杂功能。所有 teammates 先输出 plan，高风险任务需要 plan approval。\n",
		".claude/commands/codeflow-team-cleanup.md":     "# /codeflow-team-cleanup\n\n由 Team lead 清理 Agent Team。teammate 不允许自己 cleanup。\n",
		".claude/skills/codeflow-collaborative-agents/SKILL.md": `
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

` + "`" + `` + "`" + `` + "`" + `md
## Collaborative Agents Result

### Overall Conclusion

### Agent Findings

### Must Fix

### Should Fix

### Test Gap

### Risk

### Evidence

### Next Step
` + "`" + `` + "`" + `` + "`" + `
`,
		".claude/agents/codeflow-discovery-agent.md": "---\nname: codeflow-discovery-agent\ndescription: Existing Capability Discovery\n---\n\n读取 OpenSpec、代码、测试、Graphify 线索，输出复用决策。\n",
		".claude/agents/codeflow-code-reviewer.md":   "---\nname: codeflow-code-reviewer\ndescription: CodeFlow Review\n---\n\n检查需求、OpenSpec、计划、测试缺口和风险。\n",
		".claude/skills/codeflow-review/SKILL.md":    "---\nname: codeflow-review\ndescription: 审查需求改动、风险和测试缺口。\n---\n\n输出 Must Fix、Should Fix、Test Gap、Risk、Next Step。\n",
		".claude/skills/codeflow-discovery/SKILL.md": "---\nname: codeflow-discovery\ndescription: 发现已有能力和复用点。\n---\n\n用于新需求前的 Existing Capability Discovery。\n",
	}
	for p, c := range files {
		if err := writeFile(p, c, opt.DryRun); err != nil {
			return err
		}
	}
	settings := map[string]any{
		"hooks": map[string]any{
			"SessionStart": []any{map[string]any{"matcher": "", "hooks": []any{map[string]string{"type": "command", "command": "codeflow status"}}}},
		},
	}
	b, _ := json.MarshalIndent(settings, "", "  ")
	if !exists(".claude/settings.json") || opt.Force {
		if err := writeFile(".claude/settings.json", string(b)+"\n", opt.DryRun); err != nil {
			return err
		}
	}
	return nil
}

func installCodex(opt InitOptions) error {
	fmt.Println("install Codex adapter")
	block := "<!-- CodeFlow start -->\n## CodeFlow\n\n始终使用简体中文回复。新需求先判断简单/复杂；复杂需求先做 Existing Capability Discovery；代码逻辑改动默认 TDD；有代码改动必须 Review；Git 写操作必须用户确认。\n\n详细规则见 `.codeflow/workflows/`，Skills 位于 `.agents/skills/`。\n<!-- CodeFlow end -->\n"
	if err := appendBlock("AGENTS.md", "<!-- CodeFlow start -->", "<!-- CodeFlow end -->", block, opt.DryRun); err != nil {
		return err
	}
	files := map[string]string{
		".agents/skills/codeflow-discovery/SKILL.md": "---\nname: codeflow-discovery\ndescription: 发现已有能力和复用点。\n---\n\n读取 OpenSpec、代码、测试和 Graphify 线索。\n",
		".agents/skills/codeflow-review/SKILL.md":    "---\nname: codeflow-review\ndescription: 审查需求改动和风险。\n---\n\n输出 Must Fix、Should Fix、Test Gap、Risk。\n",
		".agents/skills/codeflow-collaborative-agents/SKILL.md": `
---
name: codeflow-codex-subagent-workflows
description: 用于 Codex 中通过 subagent workflows 进行并行探索、Review、测试、triage 和复杂任务拆解。
---

# CodeFlow Codex Subagent Workflows

Codex 使用 Subagent Workflows，不是 Claude Agent Teams。

## 使用规则

1. 只在用户明确要求时 spawn subagents。
2. 每个 subagent 必须有角色、范围、输出格式。
3. 默认用于 read-heavy 分析、测试、triage、review。
4. 并行写代码前必须获得用户确认。
5. 主线程必须等待所有 subagents 结果后再总结。
6. 并行写代码要谨慎，优先 worktree。
7. 输出必须包含 Must Fix / Should Fix / Test Gap / Risk / Evidence。

## 推荐 Review Prompt

请按 CodeFlow 使用 Codex subagents 审查当前分支。

Spawn one subagent per point:
1. Security
2. Code quality
3. Bugs
4. Race
5. Test flakiness
6. Maintainability

Wait for all agents, then summarize findings by:
- Must Fix
- Should Fix
- Test Gap
- Risk
- Evidence
`,
		".codeflow/prompts/codex-parallel-review.md": "请按 CodeFlow 使用 Codex subagents 审查当前分支。Spawn one subagent per point: Security, Code quality, Bugs, Race, Test flakiness, Maintainability。等待所有结果后汇总 Must Fix / Should Fix / Test Gap / Risk / Evidence。\n",
	}
	for p, c := range files {
		if err := writeFile(p, c, opt.DryRun); err != nil {
			return err
		}
	}
	return nil
}

func installCursor(opt InitOptions) error {
	fmt.Println("install Cursor adapter")
	block := "<!-- CodeFlow start -->\n## CodeFlow\n\nCursor 通过 `.cursor/rules/*.mdc` 获取 CodeFlow 规则。始终中文回复，Git 写操作必须确认。\n<!-- CodeFlow end -->\n"
	if err := appendBlock("AGENTS.md", "<!-- CodeFlow start -->", "<!-- CodeFlow end -->", block, opt.DryRun); err != nil {
		return err
	}
	rules := map[string]string{
		".cursor/rules/codeflow-core.mdc":     "---\ndescription: CodeFlow 核心规则\nalwaysApply: true\n---\n\n始终使用简体中文回复。新需求先判断简单/复杂。有代码逻辑改动必须 TDD 和 Review。\n",
		".cursor/rules/codeflow-git.mdc":      "---\ndescription: CodeFlow Git 安全规则\nalwaysApply: true\n---\n\n未经用户确认禁止 git add/commit/push/merge/rebase/reset/clean/restore/checkout/switch/stash/tag。\n",
		".cursor/rules/codeflow-graphify.mdc": "---\ndescription: Graphify 可选增强\nalwaysApply: false\n---\n\n无 graphify-out/graph.json 提醒 /graphify .；可能过期提醒 /graphify . --update。\n",
		".cursor/rules/codeflow-collaborative-agents.mdc": `
---
description: CodeFlow Cursor parallel agent collaboration rules
alwaysApply: false
---

# CodeFlow Cursor Collaborative Agents

Cursor 支持 Parallel Agents / Subagents / Cloud Agents，但 CodeFlow 不管理 Cursor runtime。

## 使用规则

1. 仅复杂任务使用 parallel agents。
2. 简单需求使用单 Agent。
3. 每个 agent 必须有明确角色、范围、输出格式。
4. 并行写代码前必须用户确认。
5. 不允许多个 agent 修改同一文件。
6. 默认先做 plan / review / investigation。
7. 主会话必须综合结果。
8. 推荐使用 worktrees 避免冲突。

## 推荐 Review Prompt

按 CodeFlow 并行 Agent 模式审查当前变更。

请将任务拆成多个独立 agent：
1. 安全和权限风险
2. 并发与性能风险
3. 测试覆盖和回归风险
4. 架构一致性和可维护性

每个 agent 只输出：
- Must Fix
- Should Fix
- Test Gap
- Risk
- Evidence

最后由主 Agent 综合成最终 Review Result。
不要自动修改代码。
`,
		".codeflow/prompts/cursor-parallel-review.md": "按 CodeFlow 并行 Agent 模式审查当前变更。拆成安全、性能、测试、维护性四个独立 agent，最后由主 Agent 综合。不要自动修改代码。\n",
	}
	for p, c := range rules {
		if err := writeFile(p, c, opt.DryRun); err != nil {
			return err
		}
	}
	return nil
}
