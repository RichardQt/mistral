package main

func consolePage() string {
	return `<!doctype html>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Mistral</title>
<style>
:root {
  color-scheme: light;
  --paper:#f2f3ef;
  --surface:#fff;
  --soft:#f7f8f4;
  --ink:#171a17;
  --muted:#6b716a;
  --line:#e0e4dc;
  --accent:#ef4d22;
  --accent-soft:#fff0e9;
  --ok:#11744a;
  --ok-soft:#eef9f3;
  --warn:#8f6200;
  --warn-soft:#fdf6e3;
  --bad:#b42318;
  --bad-soft:#fef1f0;
  --mono:ui-monospace,SFMono-Regular,"SF Mono",Menlo,Consolas,monospace;
  --sans:"Avenir Next",Avenir,ui-sans-serif,system-ui,-apple-system,"Segoe UI",sans-serif;
}
* { box-sizing:border-box; }
[hidden] { display:none !important; }
html, body { margin:0; }
body { font:14px/1.55 var(--sans); background:var(--paper); color:var(--ink); -webkit-font-smoothing:antialiased; }
button, select, input { font:inherit; }
button { border:1px solid var(--line); background:var(--surface); color:var(--ink); border-radius:9px; padding:7px 12px; cursor:pointer; transition:border-color .15s, background .15s, color .15s; }
button:hover { border-color:#c5cbc0; }
button:disabled { opacity:.45; cursor:not-allowed; }
button.primary { background:var(--ink); border-color:var(--ink); color:#fff; font-weight:600; }
button.primary:hover { background:var(--accent); border-color:var(--accent); }
button.danger { color:var(--bad); }
button.danger:hover { border-color:#e7b4af; background:var(--bad-soft); }
a { color:var(--accent); text-underline-offset:2px; }
a.button { display:inline-flex; align-items:center; gap:6px; border:1px solid var(--line); background:var(--surface); border-radius:9px; padding:7px 12px; color:var(--ink); text-decoration:none; }
a.button:hover { border-color:#c5cbc0; }
:focus-visible { outline:2px solid var(--accent); outline-offset:2px; }
.top { display:flex; justify-content:space-between; align-items:center; gap:16px; padding:16px 24px; background:var(--surface); border-bottom:1px solid var(--line); }
.brand { display:flex; align-items:center; gap:12px; min-width:0; }
.mark { width:38px; height:38px; border:1px solid #f3c6b4; border-radius:11px; background:var(--accent-soft); color:var(--accent); display:grid; place-items:center; flex:none; }
.mark svg { width:21px; height:21px; display:block; }
h1 { margin:0; font-size:19px; line-height:1.2; letter-spacing:-.02em; }
.brand p { margin:2px 0 0; color:var(--muted); font-size:12px; }
.version { font:12px var(--mono); color:var(--muted); border:1px solid var(--line); border-radius:999px; padding:4px 9px; white-space:nowrap; }
.tabs { display:flex; gap:2px; padding:0 24px; background:var(--surface); border-bottom:1px solid var(--line); overflow-x:auto; }
.tab { border:0; border-bottom:2px solid transparent; border-radius:0; background:transparent; color:var(--muted); padding:11px 14px; white-space:nowrap; font-weight:550; }
.tab:hover { color:var(--ink); border-color:var(--line); }
.tab.active { color:var(--ink); border-bottom-color:var(--accent); }
main { max-width:1280px; margin:0 auto; padding:22px 24px 64px; }
h2 { margin:0; font-size:16px; letter-spacing:-.015em; }
.view-head { display:flex; justify-content:space-between; align-items:flex-start; gap:16px; margin-bottom:14px; }
.view-head p { margin:4px 0 0; color:var(--muted); font-size:12.5px; max-width:760px; }
.card { background:var(--surface); border:1px solid var(--line); border-radius:14px; padding:20px; }
.card + .card { margin-top:16px; }
.overview { display:grid; grid-template-columns:minmax(0,1.55fr) minmax(300px,.8fr); gap:16px; align-items:start; }
.stack { min-width:0; }
.stack .card + .card { margin-top:16px; }
.eyebrow { font:11px var(--mono); letter-spacing:.12em; text-transform:uppercase; color:var(--muted); }
.route-head { display:flex; justify-content:space-between; gap:16px; align-items:flex-start; }
.route-head > div:first-child { min-width:0; }
.route-url { font:600 15px/1.45 var(--mono); word-break:break-all; margin:6px 0 4px; }
.route-note { color:var(--muted); font-size:12.5px; }
.state { display:inline-flex; align-items:center; gap:7px; border-radius:999px; border:1px solid; padding:5px 10px; font-size:12px; white-space:nowrap; }
.state i { width:6px; height:6px; border-radius:50%; background:currentColor; }
.state.ok { color:var(--ok); border-color:#c6e7d6; background:var(--ok-soft); }
.state.warn { color:var(--warn); border-color:#edd9a3; background:var(--warn-soft); }
.pulse-box { margin-top:16px; padding-top:12px; border-top:1px solid var(--line); }
.pulse-head { display:flex; justify-content:space-between; gap:12px; align-items:baseline; color:var(--muted); font-size:11.5px; }
.pulse-head b { color:var(--ink); font-weight:600; font-size:12px; }
.pulse { display:flex; align-items:flex-end; gap:3px; height:42px; margin-top:8px; }
.tick { flex:1 1 0; min-width:3px; border-radius:2px 2px 0 0; background:var(--ok); opacity:.8; }
.tick.bad { background:var(--bad); }
.pulse-empty { display:flex; align-items:center; height:42px; margin-top:8px; color:var(--muted); font-size:12px; border-bottom:1px dashed var(--line); }
.ledger { display:grid; grid-template-columns:repeat(3,minmax(0,1fr)); gap:1px; background:var(--line); border:1px solid var(--line); border-radius:11px; overflow:hidden; margin-top:16px; }
.metric { background:var(--surface); padding:12px 14px; min-width:0; }
.metric span { display:block; color:var(--muted); font-size:11px; letter-spacing:.04em; }
.metric b { display:block; font:600 21px/1.25 var(--mono); letter-spacing:-.035em; margin-top:3px; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.metric small { font:500 11px var(--sans); color:var(--muted); letter-spacing:0; margin-left:5px; }
.token-box { margin-top:16px; padding-top:14px; border-top:1px solid var(--line); }
.token-top { display:flex; justify-content:space-between; align-items:baseline; gap:12px; }
.token-rate b { font:600 20px/1 var(--mono); letter-spacing:-.03em; }
.token-rate span { color:var(--muted); font-size:11.5px; margin-left:8px; }
.token-bar { height:7px; border-radius:999px; background:var(--soft); border:1px solid var(--line); overflow:hidden; margin-top:9px; }
.token-bar i { display:block; height:100%; background:var(--accent); border-radius:999px 0 0 999px; }
.token-legend { display:grid; grid-template-columns:repeat(3,minmax(0,1fr)); gap:10px; margin-top:9px; }
.token-legend span { display:block; color:var(--muted); font-size:11px; }
.token-legend b { display:block; font:600 13px var(--mono); margin-top:1px; }
.rail-head { display:flex; justify-content:space-between; align-items:center; gap:12px; padding-bottom:12px; border-bottom:1px solid var(--line); }
.rail-head span { color:var(--muted); font-size:12px; }
.linklike { border:0; background:transparent; padding:2px; color:var(--accent); font-weight:600; font-size:12px; }
.linklike:hover { color:var(--ink); border:0; }
.cred { display:flex; gap:10px; padding:12px 0; border-bottom:1px solid var(--line); align-items:flex-start; }
.cred:last-of-type { border-bottom:0; }
.dot { width:8px; height:8px; border-radius:50%; margin-top:5px; flex:none; background:var(--ok); box-shadow:0 0 0 3px var(--ok-soft); }
.dot.warn { background:var(--warn); box-shadow:0 0 0 3px var(--warn-soft); }
.dot.off { background:#a8aea4; box-shadow:0 0 0 3px #f1f2ef; }
.cred-main { min-width:0; flex:1; }
.cred-title { display:flex; gap:7px; align-items:center; flex-wrap:wrap; }
.cred-title b { font-weight:600; }
.cred-meta { color:var(--muted); font-size:11.5px; margin-top:2px; overflow-wrap:anywhere; }
.chip { display:inline-flex; align-items:center; border:1px solid #f3c6b4; background:var(--accent-soft); color:var(--accent); border-radius:999px; padding:1px 7px; font-size:10.5px; font-weight:600; }
.badge { display:inline-flex; align-items:center; border-radius:999px; padding:2px 8px; font-size:11px; border:1px solid; white-space:nowrap; }
.badge.ok { color:var(--ok); border-color:#c6e7d6; background:var(--ok-soft); }
.badge.warn { color:var(--warn); border-color:#edd9a3; background:var(--warn-soft); }
.badge.bad { color:var(--bad); border-color:#f0c8c4; background:var(--bad-soft); }
.badge.off { color:var(--muted); border-color:var(--line); background:var(--soft); }
.badge.probe { color:var(--muted); border-color:var(--line); background:transparent; }
.err { color:var(--bad); font-size:11.5px; overflow-wrap:anywhere; }
.facts { margin-top:2px; }
.fact { display:flex; justify-content:space-between; gap:14px; padding:9px 0; border-bottom:1px solid var(--line); font-size:13px; }
.fact:last-child { border-bottom:0; padding-bottom:0; }
.fact span { color:var(--muted); }
.fact b { font-weight:600; text-align:right; }
.empty { padding:22px 0; color:var(--muted); font-size:12.5px; }
.empty b { display:block; color:var(--ink); font-size:13px; margin-bottom:2px; }
.table-wrap { overflow-x:auto; }
table { width:100%; border-collapse:collapse; }
th { font:500 10.5px var(--mono); letter-spacing:.05em; color:var(--muted); text-align:left; padding:9px 10px; border-bottom:1px solid var(--line); white-space:nowrap; }
td { padding:10px; border-bottom:1px solid var(--line); font-size:12.5px; vertical-align:middle; }
tbody tr:last-child td { border-bottom:0; }
tbody tr:hover td { background:var(--soft); }
.num { text-align:right; font-family:var(--mono); font-variant-numeric:tabular-nums; white-space:nowrap; }
.mono { font-family:var(--mono); font-size:11.5px; }
.muted { color:var(--muted); font-size:12px; }
.nowrap { white-space:nowrap; }
.grow { min-width:0; flex:1; }
.model-cell { max-width:190px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.error-cell { max-width:200px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.actions { display:flex; gap:6px; align-items:center; flex-wrap:wrap; justify-content:flex-end; }
.cred-card { display:flex; justify-content:space-between; gap:16px; align-items:flex-start; padding:15px 0; border-top:1px solid var(--line); }
.cred-card:first-of-type { border-top:0; }
.cred-card.editing { background:var(--accent-soft); margin:0 -20px; padding:15px 20px; }
.cred-card .cred-meta { font-size:12px; }
.form { display:flex; gap:8px; align-items:center; flex-wrap:wrap; margin-top:16px; padding-top:16px; border-top:1px solid var(--line); }
input, select { padding:8px 10px; border:1px solid var(--line); border-radius:9px; background:var(--surface); color:var(--ink); min-width:0; }
input.wide { width:100%; margin-top:10px; }
input.pri { width:92px; }
select { max-width:100%; }
.panel { border:1px solid var(--line); border-radius:11px; padding:14px; margin:16px 0; background:var(--soft); }
.panel .row-head { display:flex; justify-content:space-between; gap:12px; align-items:center; flex-wrap:wrap; }
.checks { display:grid; grid-template-columns:repeat(3,minmax(0,1fr)); gap:3px 12px; max-height:320px; overflow:auto; margin-top:10px; }
.check { display:flex; gap:6px; align-items:center; padding:4px 2px; font-size:12.5px; min-width:0; }
.check.disabled { opacity:.55; }
.check .mono { overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.log-summary { display:flex; flex-wrap:wrap; gap:8px 18px; padding:12px 0 4px; color:var(--muted); font-size:12px; }
.log-summary b { color:var(--ink); font-family:var(--mono); font-weight:600; margin-left:3px; }
.log-toolbar { display:flex; gap:8px; align-items:center; flex-wrap:wrap; justify-content:flex-end; }
.config-grid { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:0 28px; }
.field { display:flex; justify-content:space-between; gap:16px; align-items:center; padding:13px 0; border-top:1px solid var(--line); }
.field.full { grid-column:1/-1; flex-direction:column; align-items:stretch; }
.field.full input { width:100%; }
.field label { font-size:13px; }
.field .hint { color:var(--muted); font-size:11.5px; display:block; margin-top:2px; }
.field input[type=number] { width:110px; }
.field input[type=text], .field input:not([type]) { width:min(420px,100%); }
.switch { display:flex; align-items:center; gap:8px; }
#msg { position:fixed; left:50%; bottom:22px; transform:translateX(-50%); margin:0; padding:9px 15px; border-radius:999px; background:var(--ink); color:#fff; font-size:12.5px; box-shadow:0 8px 24px rgba(0,0,0,.18); max-width:min(80vw,640px); z-index:40; }
#msg:empty { display:none; }
@media (max-width:1000px) { .overview { grid-template-columns:1fr; } .checks { grid-template-columns:repeat(2,minmax(0,1fr)); } }
@media (max-width:720px) { main { padding:16px 14px 56px; } .top { padding:14px; } .tabs { padding:0 10px; } .ledger { grid-template-columns:repeat(2,minmax(0,1fr)); } .config-grid { grid-template-columns:1fr; } .field.full { grid-column:auto; } .checks { grid-template-columns:1fr; } .route-head { flex-direction:column; } .route-url { font-size:13px; } .cred-card { flex-direction:column; } .actions { justify-content:flex-start; } .card { padding:16px; } }
@media (prefers-reduced-motion:reduce) { * { transition:none !important; } }
</style>
<header class="top">
  <div class="brand">
    <span class="mark" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M3 8h10.5a2.5 2.5 0 1 0-2.4-3.2"/><path d="M3 12h15a2.5 2.5 0 1 1-2.4 3.2"/><path d="M3 16h7"/></svg></span>
    <div><h1>Mistral</h1><p>CPA 上游插件 · 凭据、模型与请求账本</p></div>
  </div>
  <span id="state" class="version">v0.0.0</span>
</header>
<nav class="tabs" aria-label="插件页面">
  <button class="tab active" data-tab="overview">概览</button>
  <button class="tab" data-tab="credentials">凭据</button>
  <button class="tab" data-tab="models">模型</button>
  <button class="tab" data-tab="logs">日志</button>
  <button class="tab" data-tab="config">配置</button>
</nav>
<main>
  <section id="overview" class="view"></section>
  <section id="credentials" class="view" hidden></section>
  <section id="models" class="view" hidden></section>
  <section id="logs" class="view" hidden>
    <div class="card">
      <div class="view-head">
        <div><h2>请求日志</h2><p>日志在插件数据目录的 requests.jsonl。测试请求标为 probe，不计入统计。</p></div>
        <div class="log-toolbar"><select id="log-filter"><option value="all">全部</option><option value="success">成功</option><option value="error">失败</option></select><select id="log-limit"><option value="50">50 条</option><option value="100" selected>100 条</option><option value="200">200 条</option></select><button id="log-refresh">刷新</button></div>
      </div>
      <div id="log-summary" class="log-summary"></div>
      <div id="log-table" class="table-wrap"></div>
    </div>
  </section>
  <section id="config" class="view" hidden></section>
</main>
<p id="msg" role="status" aria-live="polite"></p>
<script>
const msg = document.querySelector('#msg');
let data = null;
let managementKey = readPanelKey() || readParentKey();
let editingCredential = '';
let fetchCredential = '';
let upstreamModels = [];
let upstreamCredential = '';
let upstreamError = '';
let modelTests = {};
let lastLogs = [];
let lastLogSummary = {};
let logFilter = 'all';

function decodeAuth(raw) {
  if (!raw) return null;
  if (raw.indexOf('enc::v2::') !== 0 && raw.indexOf('enc::v1::') !== 0) {
    try { return JSON.parse(raw); } catch (error) { return null; }
  }
  const v1 = raw.indexOf('enc::v1::') === 0;
  const masks = v1
    ? ['cli-proxy-api-webui::secure-storage|' + location.host + '|' + navigator.userAgent]
    : ['cli-proxy-api-webui::secure-storage|' + location.host, 'cli-proxy-api-webui::secure-storage|v2|' + location.host];
  let binary = '';
  try { binary = atob(raw.slice(9)); } catch (error) { return null; }
  for (let n = 0; n < masks.length; n++) {
    try {
      const mask = new TextEncoder().encode(masks[n]);
      const out = new Uint8Array(binary.length);
      for (let i = 0; i < binary.length; i++) out[i] = binary.charCodeAt(i) ^ mask[i % mask.length];
      return JSON.parse(new TextDecoder().decode(out));
    } catch (error) {}
  }
  return null;
}
function readPanelKey() {
  const parsed = decodeAuth(localStorage.getItem('cli-proxy-auth'));
  const state = parsed && parsed.state ? parsed.state : parsed;
  return state && state.managementKey ? state.managementKey : '';
}
function readParentKey() {
  try {
    if (!window.parent || window.parent === window) return '';
    const root = window.parent.document.getElementById('root');
    if (!root) return '';
    const prop = Object.getOwnPropertyNames(root).find(name => name.indexOf('__reactContainer') === 0 || name.indexOf('__reactFiber') === 0);
    if (!prop) return '';
    const seen = new Set();
    let found = '';
    let nodes = 0;
    function consider(value) {
      if (!found && value && typeof value === 'object' && typeof value.managementKey === 'string' && value.managementKey.trim()) found = value.managementKey.trim();
    }
    function walk(node) {
      if (!node || found || typeof node !== 'object' || seen.has(node) || nodes > 8000) return;
      seen.add(node);
      nodes++;
      consider(node.memoizedState);
      let hook = node.memoizedState;
      for (let i = 0; hook && i < 60; i++) {
        consider(hook.memoizedState);
        consider(hook.baseState);
        if (hook.memoizedState && typeof hook.memoizedState === 'object') {
          consider(hook.memoizedState.value);
          if (typeof hook.memoizedState.getSnapshot === 'function') {
            try { consider(hook.memoizedState.getSnapshot()); } catch (error) {}
          }
        }
        if (hook.queue) consider(hook.queue.value);
        if (hook.queue && typeof hook.queue.getSnapshot === 'function') {
          try { consider(hook.queue.getSnapshot()); } catch (error) {}
        }
        hook = hook.next;
      }
      walk(node.child);
      walk(node.sibling);
    }
    const container = root[prop];
    walk(container && container.current ? container.current : container);
    return found;
  } catch (error) { return ''; }
}
function switchTab(name) {
  document.querySelectorAll('.tab').forEach(item => item.classList.toggle('active', item.dataset.tab === name));
  document.querySelectorAll('main section').forEach(section => section.hidden = section.id !== name);
  msg.textContent = '';
  if (name === 'logs') loadLogs();
}
document.querySelectorAll('.tab').forEach(tab => tab.onclick = () => switchTab(tab.dataset.tab));
async function api(body) {
  if (!managementKey) {
    msg.textContent = '没有读到管理页的登录状态。请从 CPA 管理页的 Mistral 菜单打开。';
    throw new Error('missing management session');
  }
  const response = await fetch('/v0/management/mistral/api', {
    method:'POST', headers:{'content-type':'application/json','authorization':'Bearer '+managementKey},
    body: JSON.stringify(body)
  });
  const text = await response.text();
  const parsed = (() => { try { return JSON.parse(text); } catch (error) { return null; } })();
  if (!response.ok) {
    const message = parsed && (parsed.error || parsed.message) ? (parsed.error || parsed.message) : text;
    throw new Error(response.status + ' ' + message);
  }
  if (parsed && parsed.error) throw new Error(parsed.error);
  return parsed;
}
function esc(value) {
  return String(value === null || value === undefined ? '' : value).replace(/[&<>"']/g, ch => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[ch]));
}
function clip(text, limit) {
  text = String(text === null || text === undefined ? '' : text);
  return text.length > limit ? text.slice(0, limit) + '...' : text;
}
function num(value) { return Number(value || 0).toLocaleString('zh-CN'); }
function pct(value) { return Math.round(Number(value || 0) * 100) + '%'; }
function sec(ms) {
  const value = Number(ms || 0);
  return value > 0 ? (value / 1000).toFixed(1) + 's' : '—';
}
function compact(value) {
  const n = Number(value || 0);
  if (n >= 1000000000) return (n / 1000000000).toFixed(1) + 'B';
  if (n >= 1000000) return (n / 1000000).toFixed(1) + 'M';
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k';
  return String(n);
}
function tokens(value) {
  const n = Number(value || 0);
  return n > 0 ? num(n) : '—';
}
function compactTokens(value) {
  const n = Number(value || 0);
  return n > 0 ? compact(n) : '—';
}
function timeText(value) {
  if (!value) return '';
  return String(value).slice(5, 19).replace('T', ' ');
}
function shortTime(value) {
  if (!value) return '';
  return String(value).slice(5, 16).replace('T', ' ');
}
function metric(label, value, unit) {
  return '<div class="metric"><span>' + label + '</span><b>' + value + (unit ? '<small>' + unit + '</small>' : '') + '</b></div>';
}
function fact(label, value) {
  return '<div class="fact"><span>' + label + '</span><b>' + value + '</b></div>';
}
function credentialStatus(item) {
  const health = item.health || {};
  if (item.disabled) return {cls:'off', text:'已停用'};
  if (!item.has_key) return {cls:'warn', text:'缺少 key'};
  if (health.cooldown_left_seconds > 0) return {cls:'warn', text:'冷却 ' + health.cooldown_left_seconds + 's'};
  return {cls:'ok', text:'正常'};
}
function statusBadge(item) {
  const status = credentialStatus(item);
  return '<span class="badge ' + status.cls + '">' + status.text + '</span>';
}
function credentialOptions() {
  const items = (data.credentials || []).filter(item => item.has_key);
  if (!items.length) return '<option value="">没有可用凭据</option>';
  if (!items.some(item => item.id === fetchCredential)) fetchCredential = items[0].id;
  return items.map(item => '<option value="' + esc(item.id) + '"' + (item.id === fetchCredential ? ' selected' : '') + '>' + esc(item.label) + (item.disabled ? '（已停用）' : '') + '</option>').join('');
}
function reasonOptions(current) {
  return ['','low','medium','high','xhigh','max'].map(value => '<option value="' + value + '"' + (value === current ? ' selected' : '') + '>' + (value || '透传') + '</option>').join('');
}
function healthyCount() {
  return (data.credentials || []).filter(item => !item.disabled && item.has_key && !(item.health && item.health.cooldown_left_seconds > 0)).length;
}
function render() {
  if (!data) return;
  document.querySelector('#state').textContent = 'v' + data.version;
  renderOverview();
  renderCredentials();
  renderModels();
  renderConfig();
}
function renderOverview() {
  const stats = data.stats || {};
  const creds = data.credentials || [];
  const models = data.models || [];
  const activeId = data.active || '';
  const healthy = healthyCount();
  const requests = Number(stats.requests || 0);
  const has = requests > 0;
  const recent = Array.isArray(data.recent) ? data.recent : [];
  const pulseItems = recent.filter(item => !item.probe).slice(0, 30).reverse();
  const maxDuration = pulseItems.reduce((max, item) => Math.max(max, Number(item.duration_ms || 0)), 1);
  const ticks = pulseItems.map(item => {
    const ok = item.status >= 200 && item.status < 400;
    const height = Math.max(4, Math.round(Math.log1p(Number(item.duration_ms || 0)) / Math.log1p(maxDuration) * 34));
    const label = timeText(item.time) + ' · ' + (item.model || '') + ' · ' + sec(item.duration_ms);
    return '<span class="tick' + (ok ? '' : ' bad') + '" style="height:' + height + 'px" title="' + esc(label) + '"></span>';
  }).join('');
  const pulse = pulseItems.length
    ? '<div class="pulse">' + ticks + '</div>'
    : '<div class="pulse-empty">还没有真实请求。客户端接上后，这里会显示最近 30 次的耗时分布。</div>';
  const cacheRate = has && stats.cache_hit_rate !== undefined ? Math.max(0, Math.min(100, Math.round(Number(stats.cache_hit_rate) * 100))) : 0;
  const tokenHtml = has
    ? '<div class="token-box"><div class="token-top"><span class="eyebrow">Token 账本</span><div class="token-rate"><b>' + (stats.cache_hit_rate !== undefined ? pct(stats.cache_hit_rate) : '—') + '</b><span>输入缓存命中</span></div></div>'
      + '<div class="token-bar"><i style="width:' + cacheRate + '%"></i></div>'
      + '<div class="token-legend"><div><span>输入</span><b>' + compact(stats.prompt_tokens) + '</b></div><div><span>缓存输入</span><b>' + compact(stats.cached_tokens) + '</b></div><div><span>输出</span><b>' + compact(stats.completion_tokens) + '</b></div></div></div>'
    : '<div class="token-box"><span class="eyebrow">Token 账本</span><div class="empty" style="padding:10px 0 0">还没有用量数据。</div></div>';
  const credItems = creds.map(item => {
    const health = item.health || {};
    const status = credentialStatus(item);
    return '<div class="cred"><span class="dot ' + status.cls + '"></span><div class="cred-main">'
      + '<div class="cred-title"><b>' + esc(item.label) + '</b>' + (item.id === activeId ? '<span class="chip">当前</span>' : '') + '</div>'
      + '<div class="cred-meta">' + esc(item.mask || '无 key') + ' · 优先级 ' + item.priority + '</div>'
      + (health.last_checked ? '<div class="cred-meta">检测 ' + esc(timeText(health.last_checked)) + ' · ' + sec(health.last_latency_ms) + '</div>' : '')
      + (health.last_error ? '<div class="err" title="' + esc(health.last_error) + '">' + esc(clip(health.last_error, 70)) + '</div>' : '')
      + '</div>' + statusBadge(item) + '</div>';
  }).join('');
  const credentialCard = '<div class="card"><div class="rail-head"><div><div class="eyebrow">凭据池</div><span>' + healthy + ' / ' + creds.length + ' 可用</span></div><button class="linklike" id="ov-manage-creds">管理</button></div>'
    + (credItems || '<div class="empty"><b>还没有凭据</b>先添加一把 Mistral API key，调度才会开始。</div>') + '</div>';
  const policyCard = '<div class="card"><div class="rail-head"><div><div class="eyebrow">运行参数</div></div><button class="linklike" id="ov-config">配置</button></div><div class="facts">'
    + fact('模型映射', models.length + ' 条')
    + fact('失败切换', data.failover ? '开启 · 冷却 ' + data.cooldown_seconds + 's' : '关闭')
    + fact('流空闲截断', (data.stream_idle_seconds || 0) > 0 ? data.stream_idle_seconds + 's' : '关闭')
    + fact('日志保留', data.log_retention + ' 条')
    + '</div></div>';
  const recentRows = recent.slice(0, 8).map(item => {
    const ok = item.status >= 200 && item.status < 400;
    return '<tr><td class="mono nowrap">' + esc(shortTime(item.time)) + '</td>'
      + '<td class="model-cell" title="' + esc(item.model || '') + '">' + esc(item.model) + '</td>'
      + '<td class="nowrap"><span class="badge ' + (ok ? 'ok' : 'bad') + '">' + item.status + '</span>' + (item.probe ? ' <span class="badge probe">测试</span>' : '') + '</td>'
      + '<td class="num">' + sec(item.duration_ms) + '</td>'
      + '<td class="num">' + compactTokens(item.prompt_tokens) + ' → ' + compactTokens(item.completion_tokens) + '</td></tr>';
  }).join('');
  const recentCard = '<div class="card"><div class="view-head" style="margin-bottom:8px"><div><h2>最近请求</h2><p>按时间倒序，测试请求标为 probe。</p></div><button id="ov-logs">全部日志</button></div>'
    + (recentRows
      ? '<div class="table-wrap"><table><thead><tr><th>时间</th><th>模型</th><th>状态</th><th class="num">耗时</th><th class="num">输入 → 输出</th></tr></thead><tbody>' + recentRows + '</tbody></table></div>'
      : '<div class="empty"><b>还没有请求</b>客户端把模型名指到这里后，流量会出现在这张表里。</div>')
    + '</div>';
  const routeNote = data.configured
    ? (data.active_label ? '当前凭据 ' + esc(data.active_label) + ' · 上游失败按优先级切换' : '上游失败按优先级切换凭据')
    : '还缺凭据或模型映射。先在凭据和模型页补齐。';
  document.querySelector('#overview').innerHTML = '<div class="overview"><div class="stack">'
    + '<div class="card"><div class="route-head"><div><span class="eyebrow">上游路由</span><div class="route-url">' + esc(data.base_url) + '</div><div class="route-note">' + routeNote + '</div></div>'
    + '<span class="state ' + (data.configured ? 'ok' : 'warn') + '"><i></i>' + (data.configured ? '可调度' : '还不能调度') + '</span></div>'
    + '<div class="pulse-box"><div class="pulse-head"><b>请求脉冲</b><span>最近 ' + pulseItems.length + ' 次 · 高度为耗时</span></div>' + pulse + '</div>'
    + '<div class="ledger">'
    + metric('请求', num(stats.requests))
    + metric('成功率', has ? pct(stats.success_rate) : '—')
    + metric('平均首字', sec(stats.avg_ttfb_ms))
    + metric('平均耗时', sec(stats.avg_duration_ms))
    + metric('输出速度', has ? Number(stats.tokens_per_second || 0).toFixed(1) : '—', 'tok/s')
    + metric('输入缓存', has && stats.cache_hit_rate !== undefined ? pct(stats.cache_hit_rate) : '—')
    + '</div>' + tokenHtml + '</div>'
    + recentCard + '</div><div class="stack">' + credentialCard + policyCard + '</div></div>';
  const logsButton = document.querySelector('#ov-logs');
  if (logsButton) logsButton.onclick = () => switchTab('logs');
  const credsButton = document.querySelector('#ov-manage-creds');
  if (credsButton) credsButton.onclick = () => switchTab('credentials');
  const configButton = document.querySelector('#ov-config');
  if (configButton) configButton.onclick = () => switchTab('config');
}
function renderCredentials() {
  const editing = editingCredential ? (data.credentials || []).find(item => item.id === editingCredential) : null;
  if (editingCredential && !editing) editingCredential = '';
  const rows = (data.credentials || []).map(item => {
    const health = item.health || {};
    const checked = health.last_checked ? '<div class="cred-meta">检测于 ' + esc(timeText(health.last_checked)) + ' · ' + sec(health.last_latency_ms) + '</div>' : '';
    const error = health.last_error ? '<div class="err" title="' + esc(health.last_error) + '">' + esc(clip(health.last_error, 90)) + '</div>' : '';
    return '<div class="cred-card' + (item.id === editingCredential ? ' editing' : '') + '"><div class="grow"><div class="cred-title"><b>' + esc(item.label) + '</b> <span class="mono muted">' + esc(item.id) + '</span>'
      + statusBadge(item) + (item.id === data.active ? '<span class="chip">当前</span>' : '') + '</div>'
      + '<div class="cred-meta">' + esc(item.mask || '无 key') + ' · 优先级 ' + item.priority + '</div>' + checked + error + '</div>'
      + '<div class="actions"><button data-probe="' + esc(item.id) + '">检测</button>'
      + '<button data-edit="' + esc(item.id) + '">编辑</button>'
      + (health.cooldown_left_seconds > 0 ? '<button data-reset="' + esc(item.id) + '">解除冷却</button>' : '')
      + '<button data-toggle="' + esc(item.id) + '">' + (item.disabled ? '启用' : '停用') + '</button>'
      + '<button class="danger" data-del="' + esc(item.id) + '">删除</button></div></div>';
  }).join('');
  const editor = '<div class="panel"><div class="row-head"><b>' + (editing ? '编辑 ' + esc(editing.label) : '添加凭据') + '</b><span class="muted">API key 只保存在 CPA 配置里，页面只显示掩码。</span></div>'
    + '<div class="form" style="margin-top:12px;padding-top:0;border-top:0">'
    + '<input id="c-label" placeholder="名称" value="' + esc(editing ? editing.label : '') + '">'
    + '<input id="c-key" placeholder="' + (editing ? 'API key（留空不改）' : 'API key') + '" type="password">'
    + '<input id="c-pri" class="pri" type="number" min="0" max="1000" value="' + (editing ? editing.priority : 10) + '">'
    + '<button class="primary" id="c-save">' + (editing ? '保存修改' : '添加') + '</button>'
    + (editing ? '<button id="c-cancel">取消</button>' : '') + '</div>'
    + '<div class="muted" style="margin-top:8px">优先级越大越先用。添加后可点「检测」调用 GET /models 验证 key。</div></div>';
  document.querySelector('#credentials').innerHTML =
    '<div class="card"><div class="view-head"><div><h2>凭据</h2><p>优先级高的先用。401、403、429、5xx 和传输错误会切到下一把，并按冷却基数退避。</p></div><span class="muted">' + healthyCount() + ' / ' + (data.credentials || []).length + ' 可用</span></div>'
    + (rows || '<div class="empty"><b>还没有凭据</b>添加一把 Mistral API key 后，调度和模型发现才能用。</div>')
    + editor
    + '<div class="form"><a class="button" href="https://chat.mistral.ai/code/extensions" target="_blank" rel="noopener noreferrer">创建 Vibe API key</a><span class="muted">在 Mistral 的 Code extensions 页面生成 key，再粘贴到上面的 API key 输入框。</span></div></div>';
  document.querySelector('#c-save').onclick = saveCredential;
  const cancel = document.querySelector('#c-cancel');
  if (cancel) cancel.onclick = () => { editingCredential = ''; render(); };
  document.querySelectorAll('[data-probe]').forEach(btn => btn.onclick = () => probeCredential(btn.dataset.probe));
  document.querySelectorAll('[data-edit]').forEach(btn => btn.onclick = () => { editingCredential = btn.dataset.edit; render(); });
  document.querySelectorAll('[data-reset]').forEach(btn => btn.onclick = () => applyAction({action:'reset-cooldown', id:btn.dataset.reset}, '已解除冷却'));
  document.querySelectorAll('[data-toggle]').forEach(btn => btn.onclick = () => {
    const item = (data.credentials || []).find(entry => entry.id === btn.dataset.toggle);
    if (!item) return;
    applyAction({action:'put-credential', id:item.id, disabled:!item.disabled}, item.disabled ? '已启用' : '已停用');
  });
  document.querySelectorAll('[data-del]').forEach(btn => btn.onclick = () => {
    if (!confirm('删除凭据 ' + btn.dataset.del + '？')) return;
    applyAction({action:'delete-credential', id:btn.dataset.del}, '已删除');
  });
}
async function saveCredential() {
  const label = val('#c-label');
  const key = val('#c-key');
  const priority = Number(val('#c-pri') || 10);
  if (!editingCredential && !key) { msg.textContent = '新凭据需要 API key。'; return; }
  const action = {action:'put-credential', label:label, priority:priority};
  if (editingCredential) action.id = editingCredential;
  if (key) action.api_key = key;
  const ok = await applyAction(action, editingCredential ? '已保存' : '已添加');
  if (ok) { editingCredential = ''; render(); }
}
function upstreamPanel() {
  if (upstreamError) return '<div class="panel err">' + esc(upstreamError) + '</div>';
  if (!upstreamModels.length) return '';
  const existing = new Set((data.models || []).map(item => item.upstream_id));
  const rows = upstreamModels.map(item => {
    const added = existing.has(item.id);
    return '<label class="check' + (added ? ' disabled' : '') + '" data-model="' + esc(item.id.toLowerCase()) + '"><input type="checkbox" value="' + esc(item.id) + '"' + (added ? ' disabled' : '') + '> <span class="mono">' + esc(item.id) + '</span>'
      + (item.owned_by ? ' <span class="muted">' + esc(item.owned_by) + '</span>' : '')
      + (added ? ' <span class="badge off">已添加</span>' : '') + '</label>';
  }).join('');
  return '<div class="panel"><div class="row-head"><div><b>上游模型</b> <span class="muted">' + upstreamModels.length + ' 个 · ' + esc(upstreamCredential) + '</span></div>'
    + '<div class="actions"><button id="up-all">全选</button><button class="primary" id="up-add">添加选中</button><button id="up-close">关闭</button></div></div>'
    + '<input id="up-filter" class="wide" placeholder="过滤模型 ID">'
    + '<div id="upstream-list" class="checks">' + rows + '</div></div>';
}
function renderModels() {
  const rows = (data.models || []).map(item => {
    const test = modelTests[item.id];
    let result = '<span class="muted">未测试</span>';
    if (test && test.pending) result = '<span class="muted">测试中...</span>';
    else if (test && test.ok) result = '<span class="badge ok">' + test.status + '</span> <span class="muted">' + sec(test.duration_ms) + '</span>' + (test.reply ? '<div class="muted">' + esc(clip(test.reply, 60)) + '</div>' : '');
    else if (test) result = '<span class="badge bad">失败</span><div class="err" title="' + esc(test.error || '') + '">' + esc(clip(test.error || '', 90)) + '</div>';
    return '<tr><td><b>' + esc(item.id) + '</b></td><td class="mono">' + esc(item.upstream_id) + '</td>'
      + '<td><select data-reason="' + esc(item.id) + '">' + reasonOptions(item.reasoning) + '</select></td>'
      + '<td>' + result + '</td>'
      + '<td><div class="actions"><button data-mtest="' + esc(item.id) + '">测试</button><button class="danger" data-mdel="' + esc(item.id) + '">删除</button></div></td></tr>';
  }).join('');
  document.querySelector('#models').innerHTML =
    '<div class="card"><div class="view-head"><div><h2>模型映射</h2><p>客户端模型名映射到 Mistral 上游模型 ID。reasoning 留空则透传客户端的 reasoning_effort。测试请求会计入上游账单。</p></div><span class="muted">' + (data.models || []).length + ' 条</span></div>'
    + '<div class="form" style="margin-top:0;padding-top:0;border-top:0"><select id="fetch-cred">' + credentialOptions() + '</select><button id="fetch-models">拉取上游模型</button><span class="muted">用选中的凭据调用 GET /models</span></div>'
    + upstreamPanel()
    + '<div class="table-wrap"><table><thead><tr><th>客户端</th><th>上游</th><th>思考</th><th>最近测试</th><th></th></tr></thead><tbody>'
    + (rows || '<tr><td colspan="5" class="muted">还没有模型映射。</td></tr>') + '</tbody></table></div>'
    + '<div class="form"><input id="m-id" placeholder="客户端模型名"><input id="m-up" placeholder="上游模型 ID"><select id="m-reason">' + reasonOptions('') + '</select><button class="primary" id="m-add">手动添加</button></div></div>';
  const selector = document.querySelector('#fetch-cred');
  selector.onchange = () => { fetchCredential = selector.value; };
  document.querySelector('#fetch-models').onclick = fetchUpstream;
  document.querySelectorAll('[data-mtest]').forEach(btn => btn.onclick = () => probeModel(btn.dataset.mtest));
  document.querySelectorAll('[data-mdel]').forEach(btn => btn.onclick = () => {
    if (!confirm('删除模型映射 ' + btn.dataset.mdel + '？')) return;
    delete modelTests[btn.dataset.mdel];
    applyAction({action:'delete-model', id:btn.dataset.mdel}, '已删除');
  });
  document.querySelectorAll('[data-reason]').forEach(select => select.onchange = () => {
    const item = (data.models || []).find(entry => entry.id === select.dataset.reason);
    if (!item) return;
    applyAction({action:'put-model', id:item.id, upstream_id:item.upstream_id, reasoning:select.value}, '思考档位已保存');
  });
  document.querySelector('#m-add').onclick = () => {
    const id = val('#m-id');
    const upstream = val('#m-up');
    if (!id || !upstream) { msg.textContent = '模型名和上游 ID 都要填。'; return; }
    applyAction({action:'put-model', id:id, upstream_id:upstream, reasoning:val('#m-reason')}, '模型已添加');
  };
  const upstreamList = document.querySelector('#upstream-list');
  if (upstreamList) {
    document.querySelector('#up-all').onclick = () => {
      upstreamList.querySelectorAll('input:not([disabled])').forEach(input => { input.checked = true; });
    };
    document.querySelector('#up-add').onclick = addSelectedModels;
    document.querySelector('#up-close').onclick = () => { upstreamModels = []; upstreamError = ''; render(); };
    document.querySelector('#up-filter').oninput = event => {
      const needle = event.target.value.trim().toLowerCase();
      upstreamList.querySelectorAll('[data-model]').forEach(label => { label.hidden = needle !== '' && label.dataset.model.indexOf(needle) === -1; });
    };
  }
}
async function fetchUpstream() {
  upstreamError = '';
  msg.textContent = '正在拉取上游模型...';
  try {
    const result = await api({action:'fetch-models', credential_id:fetchCredential});
    if (!result.ok) {
      upstreamModels = [];
      upstreamError = result.error || '拉取失败';
      render();
      msg.textContent = upstreamError;
      return;
    }
    upstreamModels = result.models || [];
    upstreamCredential = result.credential_label || '';
    render();
    msg.textContent = '拉取到 ' + upstreamModels.length + ' 个模型。勾选后点「添加选中」。';
  } catch (error) { msg.textContent = error.message; }
}
async function addSelectedModels() {
  const checked = Array.prototype.slice.call(document.querySelectorAll('#upstream-list input:checked'));
  if (!checked.length) { msg.textContent = '先勾选要添加的模型。'; return; }
  const existing = new Set((data.models || []).map(item => item.upstream_id));
  const items = checked.filter(input => !existing.has(input.value)).map(input => ({id:input.value, upstream_id:input.value}));
  if (!items.length) { msg.textContent = '选中的模型都已经添加。'; return; }
  await applyAction({action:'put-models', items:items}, '已添加 ' + items.length + ' 个模型');
}
async function probeCredential(id) {
  msg.textContent = '正在检测凭据...';
  try {
    const result = await api({action:'test-credential', id:id});
    data = await api({action:'state'});
    render();
    msg.textContent = result.ok ? ('凭据可用：' + result.model_count + ' 个模型，' + sec(result.duration_ms)) : ('检测失败：' + (result.error || 'unknown'));
  } catch (error) { msg.textContent = error.message; }
}
async function probeModel(id) {
  modelTests[id] = {pending:true};
  render();
  try {
    const result = await api({action:'test-model', id:id, credential_id:fetchCredential});
    modelTests[id] = result;
    data = await api({action:'state'});
    render();
    msg.textContent = result.ok ? ('模型可用：' + sec(result.duration_ms)) : ('模型测试失败：' + (result.error || ''));
  } catch (error) { delete modelTests[id]; render(); msg.textContent = error.message; }
}
function renderConfig() {
  document.querySelector('#config').innerHTML =
    '<div class="card"><div class="view-head"><div><h2>配置</h2><p>保存后写入 CPA 配置文件的 plugins.configs.mistral 块。</p></div></div>'
    + '<div class="config-grid">'
    + '<div class="field full"><label for="base">Base URL<span class="hint">Mistral API 地址，带 /v1。</span></label><input id="base" value="' + esc(data.base_url) + '"></div>'
    + '<div class="field"><label for="keep">日志保留条数<span class="hint">10 到 5000。</span></label><input id="keep" type="number" min="10" max="5000" value="' + data.log_retention + '"></div>'
    + '<div class="field"><label for="idle">流空闲截断<span class="hint">0 为关闭，单位秒。</span></label><input id="idle" type="number" min="0" max="3600" value="' + (data.stream_idle_seconds || 0) + '"></div>'
    + '<div class="field"><label for="cooldown">冷却基数<span class="hint">0 表示不冷却，单位秒。</span></label><input id="cooldown" type="number" min="0" max="3600" value="' + data.cooldown_seconds + '"></div>'
    + '<div class="field"><label>失败切换<span class="hint">401、403、429、5xx、传输错误时换下一把。</span></label><span class="switch"><input id="failover" type="checkbox"' + (data.failover ? ' checked' : '') + '> 开启</span></div>'
    + '</div><div class="form"><button class="primary" id="save-settings">保存配置</button><span class="muted">failover 关闭时只使用当前排在最前的一把凭据。</span></div></div>';
  document.querySelector('#save-settings').onclick = () => applyAction({
    action:'put-settings', base_url:val('#base'), log_retention:Number(val('#keep')),
    stream_idle_seconds:Number(val('#idle')), failover:document.querySelector('#failover').checked,
    cooldown_seconds:Number(val('#cooldown'))
  }, '配置已保存');
}
async function applyAction(action, note) {
  msg.textContent = '';
  try {
    data = await api(action);
    render();
    if (note) msg.textContent = note;
    return true;
  } catch (error) { msg.textContent = error.message; return false; }
}
async function loadLogs() {
  try {
    const result = await api({action:'logs', limit:Number(val('#log-limit') || 100)});
    lastLogs = result.logs || [];
    lastLogSummary = result.summary || {};
    renderLogs();
  } catch (error) { msg.textContent = error.message; }
}
function renderLogs() {
  const summary = lastLogSummary || {};
  document.querySelector('#log-summary').innerHTML =
    '<span>记录<b>' + lastLogs.length + '</b></span>'
    + '<span>请求<b>' + num(summary.requests) + '</b></span>'
    + '<span>成功率<b>' + (summary.success_rate !== undefined ? pct(summary.success_rate) : '—') + '</b></span>'
    + '<span>平均首字<b>' + sec(summary.avg_ttfb_ms) + '</b></span>'
    + '<span>平均耗时<b>' + sec(summary.avg_duration_ms) + '</b></span>'
    + '<span>输入<b>' + compact(summary.prompt_tokens) + '</b></span>'
    + '<span>输出<b>' + compact(summary.completion_tokens) + '</b></span>'
    + '<span>缓存输入<b>' + compact(summary.cached_tokens) + '</b></span>'
    + '<span>缓存率<b>' + (summary.cache_hit_rate !== undefined ? pct(summary.cache_hit_rate) : '—') + '</b></span>';
  const rows = lastLogs.filter(item => {
    if (logFilter === 'success') return item.status >= 200 && item.status < 400;
    if (logFilter === 'error') return item.status < 200 || item.status >= 400;
    return true;
  }).map(item => {
    const ok = item.status >= 200 && item.status < 400;
    const input = Number(item.prompt_tokens || 0);
    const cached = Number(item.cached_tokens || 0);
    const rate = input > 0 ? (item.cache_hit_rate !== undefined ? item.cache_hit_rate : cached / input) : null;
    return '<tr><td class="mono nowrap">' + esc(timeText(item.time)) + '</td>'
      + '<td class="model-cell" title="' + esc(item.model || '') + '">' + esc(item.model) + '</td>'
      + '<td class="nowrap" title="' + esc(item.credential || '') + '">' + esc(item.credential) + '</td>'
      + '<td class="nowrap"><span class="badge ' + (ok ? 'ok' : 'bad') + '">' + item.status + '</span>' + (item.probe ? ' <span class="badge probe">测试</span>' : '') + '</td>'
      + '<td class="num">' + (item.attempts || 1) + '</td>'
      + '<td class="num">' + sec(item.ttfb_ms) + '</td>'
      + '<td class="num">' + sec(item.duration_ms) + '</td>'
      + '<td class="num">' + tokens(item.prompt_tokens) + '</td>'
      + '<td class="num">' + tokens(item.completion_tokens) + '</td>'
      + '<td class="num">' + tokens(item.cached_tokens) + '</td>'
      + '<td class="num">' + (rate === null ? '—' : pct(rate)) + '</td>'
      + '<td class="num">' + (Number(item.tokens_per_second || 0) > 0 ? Number(item.tokens_per_second).toFixed(1) : '—') + '</td>'
      + '<td class="error-cell err" title="' + esc(item.error || '') + '">' + esc(clip(item.error || '', 80)) + '</td></tr>';
  }).join('');
  document.querySelector('#log-table').innerHTML = '<table><thead><tr><th>时间</th><th>模型</th><th>凭据</th><th>状态</th><th class="num">尝试</th><th class="num">首字</th><th class="num">耗时</th><th class="num">输入</th><th class="num">输出</th><th class="num">缓存输入</th><th class="num">缓存率</th><th class="num">tok/s</th><th>错误</th></tr></thead><tbody>'
    + (rows || '<tr><td colspan="13" class="muted">还没有请求记录。</td></tr>') + '</tbody></table>';
}
document.querySelector('#log-refresh').onclick = loadLogs;
document.querySelector('#log-filter').onchange = event => { logFilter = event.target.value; renderLogs(); };
document.querySelector('#log-limit').onchange = loadLogs;
function val(sel) { const node = document.querySelector(sel); return node ? node.value.trim() : ''; }
async function refresh() {
  data = await api({action:'state'});
  render();
}
refresh().catch(error => { msg.textContent = error.message; });
</script>`
}
