package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func AgentsStatus() error {
	fmt.Println("CodeFlow Collaborative Agents Status")
	fmt.Println()

	fmt.Println("Claude Code:")
	if _, err := exec.LookPath("claude"); err == nil {
		fmt.Println("- claude CLI: found")
	} else {
		fmt.Println("- claude CLI: not found")
	}
	if os.Getenv("CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS") == "1" {
		fmt.Println("- Agent Teams: enabled by CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1")
	} else {
		fmt.Println("- Agent Teams: disabled or unknown")
	}
	fmt.Println("- Runtime: managed by Claude Code")
	fmt.Println("- CodeFlow role: rules / commands / hooks / readiness guidance")
	fmt.Println()

	fmt.Println("Codex:")
	if _, err := exec.LookPath("codex"); err == nil {
		fmt.Println("- codex CLI: found")
	} else {
		fmt.Println("- codex CLI: not found")
	}
	fmt.Println("- Subagent Workflows: supported when explicitly requested")
	fmt.Println("- CodeFlow role: AGENTS.md / .agents/skills / prompts")
	fmt.Println()

	fmt.Println("Cursor:")
	fmt.Println("- Parallel Agents / Subagents: supported by Cursor runtime when available")
	fmt.Println("- CodeFlow role: .cursor/rules / prompts / worktree guidance")
	fmt.Println()

	if !exists(".codeflow/workflows/collaborative-agents.md") {
		fmt.Println("warning: missing .codeflow/workflows/collaborative-agents.md")
	}
	if !exists(".claude/commands/codeflow-team-review.md") {
		fmt.Println("warning: missing Claude team commands")
	}
	if !exists(".agents/skills/codeflow-collaborative-agents/SKILL.md") {
		fmt.Println("warning: missing Codex collaborative agents skill")
	}
	if !exists(".cursor/rules/codeflow-collaborative-agents.mdc") {
		fmt.Println("warning: missing Cursor collaborative agents rule")
	}
	return nil
}

func AgentsSuggest(task string) error {
	task = strings.TrimSpace(strings.ToLower(task))
	fmt.Println("CodeFlow Collaborative Agents Suggest")
	fmt.Println()
	if task == "" {
		fmt.Println("请提供任务描述，或参考：")
		fmt.Println("- 简单需求：executing-plans")
		fmt.Println("- 复杂审查：parallel-agent-review")
		fmt.Println("- 疑难问题：collaborative-agent-investigation")
		fmt.Println("- 跨模块功能：collaborative-agent-development")
		fmt.Println("- 同文件修改：不要使用多代理")
		return nil
	}
	if strings.Contains(task, "review") || strings.Contains(task, "审查") || strings.Contains(task, "发布") {
		fmt.Println("建议：parallel-agent-review")
		fmt.Println("角色：security / performance / test / maintainability")
		return nil
	}
	if strings.Contains(task, "调查") || strings.Contains(task, "排查") || strings.Contains(task, "bug") || strings.Contains(task, "问题") {
		fmt.Println("建议：collaborative-agent-investigation")
		fmt.Println("方式：competing hypotheses，每个代理负责一个假设，最终由主代理综合。")
		return nil
	}
	if strings.Contains(task, "跨模块") || strings.Contains(task, "架构") || strings.Contains(task, "复杂") || strings.Contains(task, "feature") {
		fmt.Println("建议：collaborative-agent-development")
		fmt.Println("要求：先 plan approval，明确每个代理文件范围，优先 worktree。")
		return nil
	}
	fmt.Println("建议：executing-plans 或 subagent-driven-development")
	fmt.Println("简单需求不要使用多代理。")
	return nil
}
