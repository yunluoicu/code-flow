#!/usr/bin/env python3
import argparse, json, shutil
from pathlib import Path
from datetime import datetime
START="<!-- CodeFlow start -->"
END="<!-- CodeFlow end -->"
SNIPPET=f"""{START}\n@.claude/codeflow/CLAUDE.md\n{END}\n"""
def backup(path: Path):
    if path.exists():
        ts=datetime.now().strftime('%Y%m%d%H%M%S')
        bak=path.with_suffix(path.suffix+f'.codeflow.bak.{ts}')
        if path.is_dir(): shutil.copytree(path,bak)
        else: shutil.copy2(path,bak)
        return bak
def merge_hooks(existing, codeflow):
    hooks=existing.setdefault('hooks',{})
    for event, entries in codeflow.get('hooks',{}).items():
        cur=hooks.setdefault(event,[])
        for entry in entries:
            s=json.dumps(entry,sort_keys=True,ensure_ascii=False)
            if not any(json.dumps(x,sort_keys=True,ensure_ascii=False)==s for x in cur): cur.append(entry)
    return existing
def main():
    ap=argparse.ArgumentParser()
    ap.add_argument('--target',default='.')
    ap.add_argument('--template',default=None)
    args=ap.parse_args()
    target=Path(args.target).resolve()
    template=Path(args.template).resolve() if args.template else Path(__file__).resolve().parents[1]/'templates'/'claude-code'
    if not target.exists(): raise SystemExit(f'目标目录不存在: {target}')
    src=template/'.claude'/'codeflow'
    dst=target/'.claude'/'codeflow'
    if not src.exists(): raise SystemExit(f'模板目录不存在: {src}')
    target.joinpath('.claude').mkdir(exist_ok=True)
    if dst.exists(): backup(dst); shutil.rmtree(dst)
    shutil.copytree(src,dst)
    for sh in dst.glob('hooks/*.sh'): sh.chmod(sh.stat().st_mode|0o111)
    claude=target/'CLAUDE.md'
    if not claude.exists(): claude.write_text(SNIPPET+'\n',encoding='utf-8')
    else:
        text=claude.read_text(encoding='utf-8')
        if START not in text: claude.write_text(text.rstrip()+'\n\n'+SNIPPET+'\n',encoding='utf-8')
    settings=target/'.claude'/'settings.json'
    cf=json.loads((template/'.claude'/'settings.codeflow.example.json').read_text(encoding='utf-8'))
    if settings.exists():
        backup(settings)
        try: data=json.loads(settings.read_text(encoding='utf-8'))
        except Exception: data={}
    else: data={}
    data=merge_hooks(data,cf)
    settings.write_text(json.dumps(data,ensure_ascii=False,indent=2)+'\n',encoding='utf-8')
    print('CodeFlow Claude Code 安装完成')
    print(f'- 目标项目: {target}')
    print(f'- 已安装目录: {dst}')
    print(f'- 已更新 CLAUDE.md: {claude}')
    print(f'- 已更新 settings.json: {settings}')
    print('请确认项目已安装 OpenSpec / OPSX 和 Superpowers。')
if __name__=='__main__': main()
