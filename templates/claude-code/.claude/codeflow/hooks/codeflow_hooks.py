#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import json, os, re, sys
from pathlib import Path

def load():
    raw=sys.stdin.read()
    try: return json.loads(raw) if raw.strip() else {}
    except Exception: return {"raw":raw}

def root(): return Path(os.environ.get('CLAUDE_PROJECT_DIR', os.getcwd()))

def deny(reason):
    print(json.dumps({"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":reason}}, ensure_ascii=False))

def graph_msg(r):
    g=r/'graphify-out/graph.json'
    if not g.exists(): return '未发现 graphify-out/graph.json。如需增强项目理解，可执行：/graphify .'
    newest=0
    for ext in ('*.go','*.js','*.ts','*.tsx','*.vue','*.py','*.java','*.md'):
        for p in r.rglob(ext):
            if 'graphify-out' in str(p) or '.git' in str(p): continue
            try: newest=max(newest,p.stat().st_mtime)
            except Exception: pass
    try:
        if newest>g.stat().st_mtime: return '检测到部分代码文件晚于 graphify-out/graph.json，图谱可能已落后。可执行：/graphify . --update'
    except Exception: pass
    return '已检测到 graphify-out/graph.json，可作为 Existing Capability Discovery 线索。'

def main():
    data=load(); r=root(); event=data.get('hook_event_name') or data.get('hookEventName') or data.get('event') or ''
    if not event:
        if data.get('tool_name'): event='PreToolUse'
        elif 'prompt' in data: event='UserPromptSubmit'
    if event=='SessionStart':
        print('## CodeFlow SessionStart')
        print('始终使用简体中文回复。新需求先判断简单/复杂。有代码逻辑改动必须 TDD 和 /review。Git 写操作必须用户确认。')
        print(graph_msg(r))
        for rel in ['.codeflow/manifest.json','.codeflow/state.md','.codeflow/active-change.md']:
            p=r/rel
            if p.exists(): print('\n## '+rel+'\n'+p.read_text(encoding='utf-8'))
        return
    if event=='UserPromptSubmit':
        prompt=data.get('prompt','') or data.get('raw','')
        print('## CodeFlow UserPromptSubmit：先判断新需求/继续旧需求/审查/验证/归档。')
        if re.search(r'实现|新增|调整|修复|优化|重构|需求|开发|补齐|改造', prompt, re.I): print('检测到可能是开发需求：先 brainstorming，再 Existing Capability Discovery。'+graph_msg(r))
        return
    if event=='PreToolUse':
        tool=data.get('tool_name',''); ti=data.get('tool_input') or {}; cmd=ti.get('command','') or ''; fp=ti.get('file_path','') or ''
        if tool=='Bash':
            if re.search(r'\bgit\s+(add|commit|push|merge|rebase|reset|clean|restore|checkout|switch|stash|tag)\b', cmd, re.I): return deny('CodeFlow 阻止了 Git 写操作。请先说明命令、目的、影响文件或分支、风险，并获得用户确认。')
            if re.search(r'rm\s+-rf|drop\s+database|deleteMany\s*\(|truncate\s+table', cmd, re.I): return deny('CodeFlow 阻止了危险命令。需要用户明确确认。')
            if re.search(r'/graphify\s+\.|graphify\s+.*--update', cmd, re.I): return deny('CodeFlow 阻止了 Graphify 写入/更新操作。执行前必须获得用户确认。')
        if tool in {'Edit','Write','MultiEdit'} and re.search(r'\.env$|\.pem$|\.key$|id_rsa|secrets|credentials', fp, re.I): return deny('CodeFlow 阻止了敏感文件修改。需要用户明确确认。')
        print('## CodeFlow PreToolUse：确认是否已有计划、是否会产生无关修改、Git 操作是否已获确认。')
        return
    if event=='Stop': print('## CodeFlow Stop：结束前确认 TDD、/review、state 更新、Git 写操作确认。')
    elif event=='SubagentStop': print('## CodeFlow SubagentStop：检查子 agent 是否输出修改文件、测试结果、阻塞点。')
    else: print('## CodeFlow Hook：始终使用简体中文回复。')
if __name__=='__main__': main()
