#!/usr/bin/env bash
set -e
INPUT="$(cat)"
TOOL="$(python3 -c 'import json,sys; print(json.loads(sys.argv[1]).get("tool_name", ""))' "$INPUT" 2>/dev/null || true)"
COMMAND="$(python3 -c 'import json,sys; print((json.loads(sys.argv[1]).get("tool_input") or {}).get("command", ""))' "$INPUT" 2>/dev/null || true)"
FILE_PATH="$(python3 -c 'import json,sys; print((json.loads(sys.argv[1]).get("tool_input") or {}).get("file_path", ""))' "$INPUT" 2>/dev/null || true)"

json_deny() {
  REASON="$1" python3 -c 'import json,os; print(json.dumps({"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":os.environ.get("REASON","")}}, ensure_ascii=False))'
}

if [ "$TOOL" = "Bash" ]; then
  if echo "$COMMAND" | grep -Eqi "rm -rf|git reset|git clean|git push|git merge|git rebase|git stash|git tag|git add|git commit|git checkout|git switch|git restore"; then
    json_deny "CodeFlow 阻止了 Git 写操作或危险命令。请先向用户说明命令、目的、影响文件、风险，并获得明确确认。"
    exit 0
  fi
fi
if [ "$TOOL" = "Edit" ] || [ "$TOOL" = "Write" ] || [ "$TOOL" = "MultiEdit" ]; then
  if echo "$FILE_PATH" | grep -Eqi "\.env$|\.pem$|\.key$|id_rsa|secrets|credentials"; then
    json_deny "敏感文件禁止直接修改。需要用户明确确认。"
    exit 0
  fi
fi
echo "## CodeFlow PreToolUse 提醒"
echo "确认是否已判断需求类型、已有计划、Git 写操作是否获用户确认。"
exit 0
