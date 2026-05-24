import React, { useEffect, useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { Activity, Boxes, GitPullRequest, Search, ShieldCheck, Sparkles, UsersRound } from 'lucide-react';
import './style.css';

type Project = {
  id: string;
  name: string;
  path: string;
  type: string;
  language: string[];
  framework: string[];
  aiTools: string[];
  requirementCount: number;
  changeCount: number;
  mustFixCount: number;
};

async function api<T>(path: string): Promise<T> {
  const res = await fetch(path);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

function App() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [summary, setSummary] = useState<any>({});
  const [route, setRoute] = useState('dashboard');

  useEffect(() => {
    Promise.all([api<any>('/api/summary'), api<{items: Project[]}>('/api/projects')]).then(([s, p]) => {
      setSummary(s);
      setProjects(p.items || []);
    });
  }, []);

  return (
    <div className="shell">
      <aside className="sidebar">
        <div className="brand"><div className="mark">CF</div><div><b>CodeFlow</b><span>AI Development Workspace</span></div></div>
        {[['dashboard','项目总览'],['projects','项目列表'],['requirements','需求管理'],['openspec','OpenSpec'],['reviews','Review / Check'],['agents','Collaborative Agents'],['graphify','Graphify']].map(([key,label]) => (
          <button key={key} className={route===key?'active':''} onClick={() => setRoute(key)}>{label}</button>
        ))}
      </aside>
      <main className="main">
        <header className="topbar"><div><h1>CodeFlow Workspace</h1><p>统一管理 AI Coding 项目、OpenSpec、Superpowers、Review 与 Graphify 状态。</p></div><div className="search"><Search size={18}/><input placeholder="搜索项目、需求、OpenSpec…" /></div></header>
        <section className="kpis">
          <Kpi icon={<Boxes/>} label="项目数" value={summary.projectCount}/>
          <Kpi icon={<Activity/>} label="活跃需求" value={summary.activeReqCount}/>
          <Kpi icon={<GitPullRequest/>} label="Changes" value={summary.activeChangeCount}/>
          <Kpi icon={<ShieldCheck/>} label="Must Fix" value={summary.mustFixCount}/>
          <Kpi icon={<UsersRound/>} label="Agents" value={3}/>
        </section>
        <section className="grid">
          {projects.map(p => <ProjectCard key={p.id} project={p}/>) }
          {!projects.length && <div className="empty"><Sparkles/>执行 <code>codeflow init --tools claude,codex,cursor</code> 后刷新 Dashboard。</div>}
        </section>
      </main>
    </div>
  );
}

function Kpi({icon,label,value}: any) { return <div className="kpi">{icon}<span>{label}</span><b>{value || 0}</b></div>; }
function ProjectCard({project}: {project: Project}) {
  return <article className="card"><h3>{project.name}</h3><p>{project.path}</p><div className="tags">{[...project.language, ...project.framework, ...project.aiTools].map(x => <span key={x}>{x}</span>)}</div><div className="stats"><b>{project.requirementCount}</b><b>{project.changeCount}</b><b>{project.mustFixCount}</b></div></article>;
}

createRoot(document.getElementById('root')!).render(<App/>);

// CodeFlow 2.1: Collaborative Agents 页面由 /api/agents 提供数据。
