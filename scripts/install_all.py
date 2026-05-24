#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import argparse, json, shutil
from pathlib import Path
from datetime import datetime
VERSION='1.2.0'; CF_START='<!-- CodeFlow start -->'; CF_END='<!-- CodeFlow end -->'
def now(): return datetime.utcnow().replace(microsecond=0).isoformat()+'Z'
def backup(p,dry=False):
    if not p.exists() or dry: return None
    b=p.with_name(p.name+'.codeflow.bak.'+datetime.now().strftime('%Y%m%d%H%M%S'))
    shutil.copytree(p,b) if p.is_dir() else shutil.copy2(p,b); return str(b)
def wr(p,s,dry=False):
    if dry: print('[dry-run] write',p); return
    p.parent.mkdir(parents=True,exist_ok=True); p.write_text(s,encoding='utf-8')
def cpdir(src,dst,args,baks):
    if not src.exists(): return
    if dst.exists():
        x=backup(dst,args.dry_run); baks.append(x) if x else None
        if not args.dry_run: shutil.rmtree(dst)
    if args.dry_run: print('[dry-run] copy dir',src,'->',dst)
    else: shutil.copytree(src,dst)
def cpfile(src,dst,args,baks):
    if not src.exists(): return
    if dst.exists():
        x=backup(dst,args.dry_run); baks.append(x) if x else None
    if args.dry_run: print('[dry-run] copy file',src,'->',dst)
    else: dst.parent.mkdir(parents=True,exist_ok=True); shutil.copy2(src,dst)
def append(p,block,args,baks):
    if p.exists():
        txt=p.read_text(encoding='utf-8')
        if CF_START in txt and CF_END in txt: print('[skip] CodeFlow block exists in',p); return
        x=backup(p,args.dry_run); baks.append(x) if x else None
        out=txt.rstrip()+'\n\n'+block.strip()+'\n'
    else: out=block.strip()+'\n'
    wr(p,out,args.dry_run)
def merge_settings(root,target,args,baks):
    src=root/'templates/claude-code/.claude/settings.codeflow.example.json'; dst=target/'.claude/settings.json'
    if not src.exists(): return
    inc=json.loads(src.read_text(encoding='utf-8')); cur={}
    if dst.exists():
        x=backup(dst,args.dry_run); baks.append(x) if x else None
        try: cur=json.loads(dst.read_text(encoding='utf-8'))
        except Exception: cur={}
    hooks=cur.setdefault('hooks',{})
    for ev,items in inc.get('hooks',{}).items():
        arr=hooks.setdefault(ev,[])
        for it in items:
            s=json.dumps(it,sort_keys=True,ensure_ascii=False)
            if not any(json.dumps(x,sort_keys=True,ensure_ascii=False)==s for x in arr): arr.append(it)
    wr(dst,json.dumps(cur,ensure_ascii=False,indent=2)+'\n',args.dry_run)
def install_core(root,target,tools,args,baks):
    cpdir(root/'templates/core/workflows',target/'.codeflow/workflows',args,baks)
    for n in ['state.md','active-change.md']: cpfile(root/f'templates/core/state/{n}',target/f'.codeflow/{n}',args,baks)
    manifest={'version':VERSION,'tools':tools,'installedAt':now(),'coreState':'.codeflow/state.md','adapters':{t:t in tools for t in ['claude','codex','cursor']}}
    wr(target/'.codeflow/manifest.json',json.dumps(manifest,ensure_ascii=False,indent=2)+'\n',args.dry_run)
def install_claude(root,target,args,baks):
    for d in ['codeflow','rules','commands','agents','skills']: cpdir(root/f'templates/claude-code/.claude/{d}',target/f'.claude/{d}',args,baks)
    append(target/'CLAUDE.md',(root/'templates/claude-code/CLAUDE_IMPORT_SNIPPET.md').read_text(encoding='utf-8'),args,baks); merge_settings(root,target,args,baks)
    hook=target/'.claude/codeflow/hooks/codeflow_hooks.py'
    if hook.exists() and not args.dry_run: hook.chmod(hook.stat().st_mode|0o111)
def install_codex(root,target,args,baks): cpdir(root/'templates/codex/.agents/skills',target/'.agents/skills',args,baks); append(target/'AGENTS.md',(root/'templates/codex/AGENTS_SNIPPET.md').read_text(encoding='utf-8'),args,baks)
def install_cursor(root,target,args,baks): cpdir(root/'templates/cursor/.cursor/rules',target/'.cursor/rules',args,baks); append(target/'AGENTS.md',(root/'templates/cursor/AGENTS_SNIPPET.md').read_text(encoding='utf-8'),args,baks)
def graph_notice(target):
    g=target/'graphify-out/graph.json'
    if not g.exists(): return '未发现 graphify-out/graph.json，可选执行：/graphify .'
    newest=0
    for ext in ('*.go','*.js','*.ts','*.tsx','*.vue','*.py','*.java','*.md'):
        for p in target.rglob(ext):
            if 'graphify-out' in str(p) or '.git' in str(p): continue
            try: newest=max(newest,p.stat().st_mtime)
            except Exception: pass
    try:
        if newest>g.stat().st_mtime: return '图谱可能已落后，可选执行：/graphify . --update'
    except Exception: pass
    return '已检测到 graphify-out/graph.json'
def main():
    ap=argparse.ArgumentParser(); ap.add_argument('--target',default='.'); ap.add_argument('--tools',default='claude'); ap.add_argument('--dry-run',action='store_true'); ap.add_argument('--upgrade',action='store_true'); ap.add_argument('--force',action='store_true'); ap.add_argument('--template-root',default=None); args=ap.parse_args()
    target=Path(args.target).resolve(); root=Path(args.template_root).resolve() if args.template_root else Path(__file__).resolve().parents[1]
    tools=['claude','codex','cursor'] if args.tools=='all' else [x.strip() for x in args.tools.split(',') if x.strip()]
    baks=[]; print(f'CodeFlow {VERSION} install target={target} tools={tools} dry_run={args.dry_run}')
    install_core(root,target,tools,args,baks)
    if 'claude' in tools: install_claude(root,target,args,baks)
    if 'codex' in tools: install_codex(root,target,args,baks)
    if 'cursor' in tools: install_cursor(root,target,args,baks)
    cpfile(root/'templates/claude-code/.graphifyignore.example',target/'.graphifyignore.example',args,baks)
    print('\n安装报告'); print('- backups:', baks if baks else '无'); print('- graphify:', graph_notice(target)); print('- 未自动修改 .gitignore，未执行任何 Git 写操作。')
    if args.dry_run: print('- dry-run 模式：未实际写入。')
if __name__=='__main__': main()
