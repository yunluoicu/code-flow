#!/usr/bin/env bash
set -e
ROOT="${CLAUDE_PROJECT_DIR:-$(pwd)}"
STATE="$ROOT/.claude/codeflow/state/state.md"
ACTIVE="$ROOT/.claude/codeflow/state/active-change.md"
echo "## CodeFlow 会话启动提醒"
echo "始终使用简体中文回复。"
echo "新需求必须先判断：简单需求 / 复杂需求。"
echo "复杂需求流程：brainstorming → discovery → /opsx:propose → Spec Review → writing-plans → 选择执行方式 → TDD → /review → 询问 verification / finishing → /opsx:archive"
echo "简单需求流程：brainstorming → discovery → writing-plans → 选择执行方式 → TDD → /review → 询问 verification / finishing"
echo "Git 写操作必须用户确认。"
[ -f "$STATE" ] && { echo; echo "## 当前 CodeFlow State"; cat "$STATE"; }
[ -f "$ACTIVE" ] && { echo; echo "## 当前 Active Change"; cat "$ACTIVE"; }
