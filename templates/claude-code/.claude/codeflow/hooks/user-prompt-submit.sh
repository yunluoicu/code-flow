#!/usr/bin/env bash
set -e
INPUT="$(cat)"
PROMPT="$(python3 -c 'import json,sys; print((json.loads(sys.argv[1]).get("prompt") or "") if sys.argv[1] else "")' "$INPUT" 2>/dev/null || true)"
echo "## CodeFlow 输入处理提醒"
echo "始终使用简体中文回复。"
echo "先判断：新需求 / 继续旧需求 / 审查 / 验证 / 归档。"
echo "先判断：简单需求 / 复杂需求。"
echo "开发类需求不要直接写代码，先 brainstorming 和 Existing Capability Discovery。"
if echo "$PROMPT" | grep -Eqi "实现|新增|调整|修复|优化|重构|需求|开发|补齐|改造"; then echo "检测到可能是开发需求：复杂需求需 /opsx:propose，简单需求可走轻流程。"; fi
if echo "$PROMPT" | grep -Eqi "review|审查|是否可以发布|是否达到发布要求|检查当前分支"; then echo "检测到可能是审查：优先使用 Superpowers /review，只审查，不修改代码。"; fi
