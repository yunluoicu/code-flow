const state = {
  route: location.hash.replace('#', '') || 'dashboard',
  projects: [],
  currentProject: null,
  summary: {},
  loading: false,
  searchQuery: '',
};

const navItems = [
  ['dashboard', '⌘', '项目总览'],
  ['projects', '▣', '项目列表'],
  ['requirements', '✓', '需求管理'],
  ['iterations', '◎', '迭代管理'],
  ['openspec', '◇', 'OpenSpec'],
  ['superpowers', '✦', 'Superpowers'],
  ['reviews', '◆', 'Review / Check'],
  ['agents', '♟', 'Collaborative Agents'],
  ['graphify', '⌁', 'Graphify'],
  ['search', '⌕', '全文搜索'],
  ['settings', '⚙', '设置'],
];

function qs(sel) { return document.querySelector(sel); }
function fmt(n) { return Number(n || 0).toLocaleString(); }
function safe(v) { return (v === undefined || v === null || v === '') ? '-' : String(v); }
function arr(v) { return Array.isArray(v) ? v : (String(v || '').split(',').filter(Boolean)); }
function statusClass(v) {
  const s = String(v || '').toLowerCase();
  if (['done','released','ready','archive-ready','success','passed'].includes(s)) return 'success';
  if (['warning','planned','planning','executing','review','active','stale'].includes(s)) return 'warn';
  if (['danger','failed','must-fix','blocked','missing'].includes(s)) return 'danger';
  return 'brand';
}
async function api(path) {
  const res = await fetch(path);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}
function setRoute(route) {
  state.route = route;
  location.hash = route;
  render();
  loadRoute();
}
window.addEventListener('hashchange', () => {
  state.route = location.hash.replace('#', '') || 'dashboard';
  render();
  loadRoute();
});

async function bootstrap() {
  renderShell();
  await loadBase();
  await loadRoute();
}

async function loadBase() {
  try {
    const [summary, projects] = await Promise.all([api('/api/summary'), api('/api/projects')]);
    state.summary = summary;
    state.projects = projects.items || [];
    state.currentProject = state.currentProject || state.projects[0] || null;
  } catch (err) {
    console.error(err);
  }
}

async function loadRoute() {
  await loadBase();
  render();
}

function renderShell() {
  document.getElementById('app').innerHTML = `
    <div class="shell">
      <aside class="sidebar">
        <div class="logo">
          <div class="logo-mark">CF</div>
          <div>
            <div class="logo-title">CodeFlow</div>
            <div class="logo-subtitle">AI Development Workspace</div>
          </div>
        </div>
        <div class="nav-section">
          <div class="nav-label">Workspace</div>
          ${navItems.map(([key, icon, label]) => `<div class="nav-item" data-route="${key}"><span class="nav-icon">${icon}</span><span>${label}</span></div>`).join('')}
        </div>
        <div class="sidebar-card">
          <strong>工作流守护</strong>
          <p>TDD 与 Review 是代码改动的默认质量门禁；Git 写操作必须用户确认。</p>
        </div>
      </aside>
      <main class="main">
        <header class="topbar">
          <div>
            <h1 class="page-title" id="pageTitle">项目总览</h1>
            <p class="page-subtitle" id="pageSubtitle">统一查看项目、需求、OpenSpec、Review 和 Graphify 状态。</p>
          </div>
          <div class="searchbar">
            <span>⌕</span>
            <input id="globalSearch" placeholder="搜索项目、需求、OpenSpec、Review…" />
          </div>
          <div class="header-actions">
            <span class="pill"><span class="dot"></span>Local 4399</span>
            <span class="pill">SQLite + FTS5</span>
          </div>
        </header>
        <section id="view"></section>
      </main>
    </div>`;
  document.querySelectorAll('.nav-item').forEach(el => el.addEventListener('click', () => setRoute(el.dataset.route)));
  qs('#globalSearch').addEventListener('keydown', e => {
    if (e.key === 'Enter') {
      state.searchQuery = e.target.value.trim();
      setRoute('search');
    }
  });
}

function setHeader(title, subtitle) {
  qs('#pageTitle').textContent = title;
  qs('#pageSubtitle').textContent = subtitle;
  document.querySelectorAll('.nav-item').forEach(el => el.classList.toggle('active', el.dataset.route === state.route));
}

function render() {
  const map = {
    dashboard: renderDashboard,
    projects: renderProjects,
    requirements: () => renderListPage('需求管理', '查看需求状态、关联 OpenSpec 和 Review 结果。', 'requirements'),
    iterations: () => renderListPage('迭代管理', '查看迭代进度、发布状态和需求关联。', 'iterations'),
    openspec: renderOpenSpec,
    superpowers: renderSuperpowers,
    reviews: renderReviews,
    graphify: renderGraphify,
    search: renderSearch,
    settings: renderSettings,
  };
  (map[state.route] || renderDashboard)();
}

function kpi(label, value, hint, tone='brand') {
  return `<div class="card kpi-card"><div class="kpi-label"><span>${label}</span><span class="tag ${tone}">${tone}</span></div><div class="kpi-value">${fmt(value)}</div><div class="kpi-hint">${hint}</div></div>`;
}
function projectCard(p) {
  const tags = [...arr(p.language).slice(0,3).map(x => `<span class="tag brand">${x}</span>`), ...arr(p.framework).slice(0,3).map(x => `<span class="tag purple">${x}</span>`), ...arr(p.aiTools).slice(0,3).map(x => `<span class="tag success">${x}</span>`)].join('');
  const must = Number(p.mustFixCount || 0);
  return `<div class="card project-card" onclick="openProject('${p.id}')">
    <div class="project-head"><div><h3 class="project-name">${safe(p.name)}</h3><div class="project-path">${safe(p.path)}</div></div><span class="tag ${must ? 'danger' : 'success'}">${must ? 'Must Fix' : 'Healthy'}</span></div>
    <div class="tags">${tags || '<span class="tag">未识别技术栈</span>'}</div>
    <div class="mini-stats"><div class="mini-stat"><b>${fmt(p.requirementCount)}</b><span>需求</span></div><div class="mini-stat"><b>${fmt(p.changeCount)}</b><span>Changes</span></div><div class="mini-stat"><b>${fmt(p.mustFixCount)}</b><span>Must Fix</span></div></div>
  </div>`;
}
window.openProject = async function(id) {
  const data = await api('/api/project?id=' + encodeURIComponent(id));
  state.currentProject = data;
  state.route = 'project-detail';
  location.hash = 'project-detail';
  renderProjectDetail(data);
};

function renderDashboard() {
  setHeader('项目总览', '统一查看项目、需求、OpenSpec、Review 和 Graphify 状态。');
  const s = state.summary || {};
  const projects = state.projects || [];
  qs('#view').innerHTML = `
    <div class="grid grid-4">
      ${kpi('项目数', s.projectCount, '已同步到本地 SQLite 的项目', 'brand')}
      ${kpi('活跃需求', s.activeReqCount, '未完成或未归档的需求', 'warn')}
      ${kpi('OpenSpec Changes', s.activeChangeCount, '当前活跃变更', 'purple')}
      ${kpi('Must Fix', s.mustFixCount, 'Review 阻塞项', Number(s.mustFixCount) ? 'danger' : 'success')}
    </div>
    <div class="section-title"><div><h2>项目目录</h2><p>来自 .codeflow、OpenSpec、Superpowers 和源码扫描。</p></div><span class="pill">最近同步：${safe(s.updatedAt)}</span></div>
    <div class="grid grid-3">${projects.map(projectCard).join('') || empty('还没有项目', '在项目根目录执行 codeflow init --tools claude,codex,cursor')}</div>
    <div class="section-title"><div><h2>流程提醒</h2><p>CodeFlow 2.0 默认不自动执行 Git 和 Graphify。</p></div></div>
    <div class="grid grid-3">
      ${noticeCard('TDD + Review', '代码逻辑改动默认必须执行 TDD 和 /review，Review 发现 must-fix 需先修复。', 'success')}
      ${noticeCard('Graphify 可选', '缺少 graphify-out 时提示 /graphify .，图谱可能落后时提示 /graphify . --update。', 'warn')}
      ${noticeCard('本地工作台', 'Dashboard 默认监听 127.0.0.1:4399，不做公网权限和团队账号。', 'brand')}
    </div>`;
}

function noticeCard(title, text, tone) {
  return `<div class="card card-pad"><span class="tag ${tone}">${tone}</span><h3>${title}</h3><p class="markdown-snippet">${text}</p></div>`;
}
function empty(title, cmd) {
  return `<div class="card empty"><div class="empty-icon">⌘</div><h3>${title}</h3><p>${cmd}</p><span class="command">${cmd}</span></div>`;
}

function renderProjects() {
  setHeader('项目列表', '按项目查看技术栈、AI 工具、需求、OpenSpec 和 Review 风险。');
  qs('#view').innerHTML = `<div class="grid grid-3">${(state.projects || []).map(projectCard).join('') || empty('没有项目', 'codeflow init --tools claude,codex,cursor')}</div>`;
}

async function renderProjectDetail(p) {
  setHeader(p.name || '项目详情', '项目画像、模块、需求、OpenSpec 和风险概览。');
  const [req, it, changes, reviews, checks, graph] = await Promise.all([
    api('/api/requirements?projectId=' + p.id), api('/api/iterations?projectId=' + p.id), api('/api/changes?projectId=' + p.id), api('/api/reviews?projectId=' + p.id), api('/api/checks?projectId=' + p.id), api('/api/graphify?path=' + encodeURIComponent(p.path || ''))
  ]);
  qs('#view').innerHTML = `
    <div class="detail-hero"><h2>${safe(p.name)}</h2><p>${safe(p.description || 'CodeFlow 已同步该项目。可通过 profile/index/sync 进一步完善项目画像。')}</p><div class="tags">${arr(p.language).map(x=>`<span class="tag brand">${x}</span>`).join('')}${arr(p.framework).map(x=>`<span class="tag purple">${x}</span>`).join('')}${arr(p.aiTools).map(x=>`<span class="tag success">${x}</span>`).join('')}</div></div>
    <div class="grid grid-4" style="margin-top:18px">${kpi('需求', p.counts?.requirements, '项目需求数量')}${kpi('迭代', p.counts?.iterations, '项目迭代数量')}${kpi('Changes', p.counts?.changes, 'OpenSpec 变更')}${kpi('Review', p.counts?.reviews, '审查记录')}</div>
    <div class="grid grid-2" style="margin-top:18px">
      ${tableCard('模块地图', ['模块','路径','职责'], (p.modules || []).map(m => [m.name, m.path, m.responsibilities]))}
      ${graphCard(graph)}
    </div>
    <div class="grid grid-2" style="margin-top:18px">
      ${tableCard('最近需求', ['ID','标题','状态'], (req.items || []).slice(0,8).map(x => [x.id, x.title, badge(x.status)]))}
      ${tableCard('OpenSpec Changes', ['Change','状态','任务'], (changes.items || []).slice(0,8).map(x => [x.changeId, badge(x.status), progress(x.taskDone, x.taskTotal)]))}
    </div>
    <div style="margin-top:18px">${tableCard('Review / Check 风险', ['类型','结论/状态','详情'], [...(reviews.items || []).slice(0,5).map(x => [x.reviewType, badge(x.conclusion || x.riskLevel), `Must ${x.mustFixCount || 0} · Gap ${x.testGapCount || 0}`]), ...(checks.items || []).slice(0,5).map(x => [x.checkType, badge(x.status), x.message])])}</div>`;
}

async function renderListPage(title, subtitle, type) {
  setHeader(title, subtitle);
  const current = state.currentProject || state.projects[0];
  let data = {items: []};
  if (type === 'requirements') data = await api('/api/requirements?projectId=' + (current?.id || ''));
  if (type === 'iterations') data = await api('/api/iterations?projectId=' + (current?.id || ''));
  const rows = (data.items || []).map(x => type === 'requirements' ? [x.id, x.title, badge(x.status), x.changeId || '-', x.filePath || '-'] : [x.id, x.name, badge(x.status), x.releaseDate || '-', x.filePath || '-']);
  qs('#view').innerHTML = `${projectSelector()}${tableCard(title, type === 'requirements' ? ['ID','标题','状态','Change','文件'] : ['ID','名称','状态','发布日期','文件'], rows)}`;
}

async function renderOpenSpec() {
  setHeader('OpenSpec', '查看 specs、active changes、tasks 完成度和归档状态。');
  const current = state.currentProject || state.projects[0];
  const data = await api('/api/changes?projectId=' + (current?.id || ''));
  qs('#view').innerHTML = `${projectSelector()}${tableCard('OpenSpec Changes', ['Change','状态','proposal','design','tasks','完成度'], (data.items || []).map(x => [x.changeId, badge(x.status), x.proposalPath, x.designPath, x.tasksPath, progress(x.taskDone, x.taskTotal)]))}`;
}

function renderSuperpowers() {
  setHeader('Superpowers', '以过程证据视角查看 brainstorming、plans、TDD、review、verification 和 finishing。');
  qs('#view').innerHTML = `<div class="card card-pad"><div class="timeline">
    ${['brainstorming 需求澄清','writing-plans 实施计划','executing-plans / subagent-driven-development 执行','test-driven-development 测试驱动','requesting-code-review 阶段审查','verification-before-completion 完成前验证','finishing-a-development-branch 收尾报告'].map((x,i)=>`<div class="timeline-item"><div class="timeline-dot"></div><div class="timeline-body"><div class="timeline-title">${x}</div><div class="timeline-meta">CodeFlow sync 会扫描相关过程文档，并在需求详情中展示关联证据。</div></div></div>`).join('')}
  </div></div>`;
}

async function renderReviews() {
  setHeader('Review / Check', '查看 must-fix、test-gap、风险和建议动作。');
  const current = state.currentProject || state.projects[0];
  const [reviews, checks] = await Promise.all([api('/api/reviews?projectId=' + (current?.id || '')), api('/api/checks?projectId=' + (current?.id || ''))]);
  qs('#view').innerHTML = `${projectSelector()}<div class="grid grid-2">${tableCard('Review 结果', ['类型','结论','Must','Should','Gap','风险'], (reviews.items || []).map(x => [x.reviewType, badge(x.conclusion), x.mustFixCount, x.shouldFixCount, x.testGapCount, badge(x.riskLevel)]))}${tableCard('Check 风险', ['类型','状态','消息','时间'], (checks.items || []).map(x => [x.checkType, badge(x.status), x.message, x.createdAt]))}</div>`;
}

async function renderGraphify() {
  setHeader('Graphify', '项目知识图谱状态与建议命令。');
  const current = state.currentProject || state.projects[0];
  const graph = await api('/api/graphify?path=' + encodeURIComponent(current?.path || ''));
  qs('#view').innerHTML = `${projectSelector()}${graphCard(graph, true)}`;
}

async function renderSearch() {
  setHeader('全文搜索', '搜索项目、需求、迭代、OpenSpec、Superpowers、Review 和 docs。');
  const q = state.searchQuery || qs('#globalSearch')?.value || '';
  const data = q ? await api('/api/search?q=' + encodeURIComponent(q)) : {items: []};
  qs('#view').innerHTML = `<div class="card card-pad"><h3>搜索</h3><div class="searchbar" style="box-shadow:none;max-width:none"><span>⌕</span><input id="pageSearch" value="${q}" placeholder="输入关键词，例如 human handoff、review、FAQ…" /></div><div class="tabs">${['All','Projects','Requirements','OpenSpec','Superpowers','Reviews','Docs'].map((x,i)=>`<span class="tab ${i===0?'active':''}">${x}</span>`).join('')}</div></div><div style="margin-top:18px">${searchResults(data.items || [])}</div>`;
  qs('#pageSearch').addEventListener('keydown', e => { if (e.key === 'Enter') { state.searchQuery = e.target.value.trim(); renderSearch(); } });
}

function renderSettings() {
  setHeader('设置', '查看本地工作台配置、端口、安全边界和推荐命令。');
  qs('#view').innerHTML = `<div class="grid grid-2">
    ${noticeCard('Web 端口', '默认 127.0.0.1:4399，可通过 codeflow web --port 18080 自定义。', 'brand')}
    ${noticeCard('安全边界', '本地访问，无登录权限系统；不会自动执行 Git、Graphify 或远程操作。', 'success')}
    ${noticeCard('Graphify 命令', '首次生成：/graphify .；更新图谱：/graphify . --update。', 'warn')}
    ${noticeCard('SQLite + FTS5', '数据库默认位于 ~/.codeflow/codeflow.db，用于索引和全文搜索。', 'purple')}
  </div>`;
}

function projectSelector() {
  return `<div class="card card-pad" style="margin-bottom:18px"><strong>当前项目</strong><div class="tags">${(state.projects || []).slice(0,8).map(p => `<span class="tag ${state.currentProject?.id === p.id ? 'brand' : ''}" onclick="selectProject('${p.id}')">${p.name}</span>`).join('') || '<span class="tag">暂无项目</span>'}</div></div>`;
}
window.selectProject = async function(id) { state.currentProject = (await api('/api/project?id=' + id)); render(); };
function badge(v) { const c = statusClass(v); return `<span class="tag ${c}">${safe(v)}</span>`; }
function progress(done, total) { const d = Number(done || 0), t = Number(total || 0); const pct = t ? Math.round(d*100/t) : 0; return `<div style="min-width:120px"><div class="progress"><i style="width:${pct}%"></i></div><div class="timeline-meta" style="margin-top:6px">${d}/${t || 0}</div></div>`; }
function tableCard(title, headers, rows) {
  return `<div class="card table-card"><div class="table-header"><div><strong>${title}</strong><div class="timeline-meta">${rows.length} 条记录</div></div></div>${rows.length ? `<table><thead><tr>${headers.map(h=>`<th>${h}</th>`).join('')}</tr></thead><tbody>${rows.map(r=>`<tr>${r.map(c=>`<td>${c === undefined ? '-' : c}</td>`).join('')}</tr>`).join('')}</tbody></table>` : `<div class="empty"><div class="empty-icon">∅</div><h3>暂无数据</h3><p>执行 codeflow sync 或 codeflow init 后刷新。</p></div>`}</div>`;
}
function graphCard(graph, large=false) {
  const tone = graph.status === 'ready' ? 'success' : graph.status === 'stale' ? 'warn' : 'danger';
  return `<div class="card card-pad ${large ? '' : ''}"><span class="tag ${tone}">${safe(graph.status)}</span><h3>Graphify 状态</h3><p class="markdown-snippet">${safe(graph.message)}</p>${graph.updatedAt ? `<p class="timeline-meta">更新时间：${graph.updatedAt}</p>` : ''}<span class="command">${graph.status === 'missing' ? '/graphify .' : graph.status === 'stale' ? '/graphify . --update' : '/graphify query "相关模块"'}</span>${graph.summary ? `<div class="markdown-snippet" style="margin-top:18px;white-space:pre-wrap">${graph.summary}</div>` : ''}</div>`;
}
function searchResults(items) {
  if (!items.length) return empty('暂无搜索结果', '试试搜索：OpenSpec、Review、需求名称、模块名');
  return `<div class="grid">${items.map(x => `<div class="card card-pad"><div style="display:flex;justify-content:space-between;gap:12px"><div><h3 style="margin:0 0 8px">${x.title}</h3><div class="timeline-meta">${x.path}</div></div><span class="tag ${statusClass(x.type)}">${x.type}</span></div><p class="markdown-snippet">${x.snippet || ''}</p></div>`).join('')}</div>`;
}

bootstrap();


async function renderCollaborativeAgents() {
  const el = document.querySelector('.content') || document.querySelector('main') || document.body;
  try {
    const data = await api('/api/agents');
    el.innerHTML = `
      <section class="page-head"><h2>Collaborative Agents</h2><p>统一管理 Claude Agent Teams、Codex Subagent Workflows、Cursor Parallel Agents 的使用规范。</p></section>
      <section class="grid cards">
        ${(data.tools || []).map(t => `<article class="card"><h3>${t.name}</h3><p>${t.native}</p><span class="badge">${t.runtime}</span><p>${t.codeflowRole}</p></article>`).join('')}
      </section>
      <section class="card"><h3>通用规则</h3><ul>${(data.rules || []).map(r => `<li>${r}</li>`).join('')}</ul></section>
    `;
  } catch (e) {
    el.innerHTML = '<section class="card danger">Collaborative Agents 加载失败</section>';
  }
}
const _oldSetRoute = typeof setRoute === 'function' ? setRoute : null;
if (_oldSetRoute) {
  window.addEventListener('hashchange', () => {
    if ((location.hash || '').replace('#','') === 'agents') setTimeout(renderCollaborativeAgents, 0);
  });
}
if ((location.hash || '').replace('#','') === 'agents') setTimeout(renderCollaborativeAgents, 500);
