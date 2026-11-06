package main

func consolePage() string {
	return `<!doctype html>
<meta charset="utf-8">
<title>Mistral</title>
<style>
:root { color-scheme: light; --bg:#f6f7f9; --card:#fff; --ink:#1c1917; --muted:#6b6560; --line:#e7e5e4; --accent:#9a3412; --ok:#166534; --warn:#a16207; --bad:#9f1239; }
* { box-sizing: border-box; }
body { margin:0; font:14px/1.5 ui-sans-serif,system-ui,sans-serif; background:var(--bg); color:var(--ink); }
header { display:flex; justify-content:space-between; align-items:center; gap:16px; padding:18px 28px; background:var(--card); border-bottom:1px solid var(--line); }
h1 { font-size:18px; margin:0; }
nav { display:flex; gap:6px; padding:12px 28px; flex-wrap:wrap; }
button, select, input { font:inherit; }
button { border:1px solid var(--line); background:var(--card); border-radius:8px; padding:6px 11px; cursor:pointer; color:var(--ink); }
button:hover { border-color:#c8c3bd; }
a.button { display:inline-flex; align-items:center; gap:6px; border:1px solid var(--line); background:var(--card); border-radius:8px; padding:6px 11px; color:var(--ink); text-decoration:none; }
a.button:hover { border-color:#c8c3bd; }
button.danger { color:var(--bad); }
button:disabled { opacity:.45; cursor:not-allowed; }
.tab { border:1px solid var(--line); background:var(--card); border-radius:8px; padding:7px 12px; cursor:pointer; }
.tab.active { background:var(--ink); color:#fff; border-color:var(--ink); }
main { padding:0 28px 48px; max-width:1080px; }
.card { background:var(--card); border:1px solid var(--line); border-radius:12px; padding:18px; margin-top:12px; }
.head { display:flex; flex-direction:column; gap:4px; margin-bottom:12px; }
.row { display:flex; gap:14px; align-items:center; justify-content:space-between; padding:11px 0; border-top:1px solid var(--line); }
.row:first-child { border-top:0; }
.grow { min-width:0; flex:1; }
.actions { display:flex; gap:6px; align-items:center; flex-wrap:wrap; justify-content:flex-end; }
.hero { display:flex; justify-content:space-between; gap:16px; align-items:flex-start; }
.base { font-size:15px; font-weight:600; margin:8px 0 4px; word-break:break-all; }
.muted { color:var(--muted); font-size:12px; }
.small { font-size:12px; }
.err { color:var(--bad); }
.mono { font-family:ui-monospace,SFMono-Regular,Menlo,monospace; font-size:12px; }
.pill { display:inline-block; padding:2px 8px; border-radius:999px; font-size:12px; border:1px solid var(--line); white-space:nowrap; }
.pill.ok { color:var(--ok); border-color:#bbf7d0; background:#f0fdf4; }
.pill.warn { color:var(--warn); border-color:#fde68a; background:#fffbeb; }
.pill.off { color:var(--muted); background:#f5f5f4; }
.grid { display:grid; grid-template-columns:repeat(6,minmax(0,1fr)); gap:10px; margin:16px 0; }
.stat { border:1px solid var(--line); border-radius:10px; padding:10px; }
.stat b { display:block; font-size:18px; margin-top:2px; }
.rows { border-top:1px solid var(--line); }
.form { display:flex; gap:8px; align-items:center; flex-wrap:wrap; margin-top:12px; }
input, select { padding:7px 8px; border:1px solid var(--line); border-radius:8px; background:#fff; color:var(--ink); }
input.wide { width:100%; margin-top:10px; }
table { width:100%; border-collapse:collapse; margin-top:10px; }
td, th { text-align:left; padding:9px 6px; border-bottom:1px solid var(--line); font-size:13px; vertical-align:top; }
th { color:var(--muted); font-weight:600; font-size:12px; }
.panel { border:1px solid var(--line); border-radius:10px; padding:12px; margin-top:12px; background:#fafaf9; }
.checks { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:4px 12px; max-height:320px; overflow:auto; margin-top:8px; }
.check { display:flex; gap:6px; align-items:center; padding:4px 2px; font-size:13px; }
.check.disabled { opacity:.55; }
#msg { min-height:20px; color:var(--muted); font-size:13px; }
@media (max-width: 820px) { .grid { grid-template-columns:repeat(3,minmax(0,1fr)); } .checks { grid-template-columns:1fr; } .hero { flex-direction:column; } }
</style>
<header><div><h1>Mistral</h1><div class="muted">CPA 插件。凭据、模型映射、请求日志。客户端走 CPA 自己的地址和 key。</div></div><div id="state" class="muted"></div></header>
<nav>
<button class="tab active" data-tab="overview">概览</button>
<button class="tab" data-tab="credentials">凭据</button>
<button class="tab" data-tab="models">模型</button>
<button class="tab" data-tab="logs">日志</button>
<button class="tab" data-tab="config">配置</button>
</nav>
<main>
<section id="overview" class="card"></section>
<section id="credentials" class="card" hidden></section>
<section id="models" class="card" hidden></section>
<section id="logs" class="card" hidden><div class="head"><b>请求日志</b><span class="muted">日志在插件数据目录的 requests.jsonl。测试请求标为 probe，不计入统计。</span></div><div class="form"><select id="log-filter"><option value="all">全部</option><option value="success">成功</option><option value="error">失败</option></select><select id="log-limit"><option value="50">50 条</option><option value="100" selected>100 条</option><option value="200">200 条</option></select><button id="log-refresh">刷新</button></div><div id="log-summary" class="muted" style="margin-top:10px"></div><div id="log-table"></div></section>
<section id="config" class="card" hidden></section>
<p id="msg" class="muted"></p>
</main>
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
document.querySelectorAll('.tab').forEach(tab => tab.onclick = () => {
  document.querySelectorAll('.tab').forEach(item => item.classList.toggle('active', item === tab));
  document.querySelectorAll('main section').forEach(section => section.hidden = section.id !== tab.dataset.tab);
  msg.textContent = '';
  if (tab.dataset.tab === 'logs') loadLogs();
});
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
function stat(label, value) { return '<div class="stat"><span class="muted">' + label + '</span><b>' + value + '</b></div>'; }
function row(label, value) { return '<div class="row"><span>' + label + '</span><span>' + value + '</span></div>'; }
function timeText(value) {
  if (!value) return '';
  return String(value).slice(5, 19).replace('T', ' ');
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
  const rate = stats.requests ? pct(stats.success_rate) : '0%';
  document.querySelector('#overview').innerHTML =
    '<div class="hero"><div><span class="pill ' + (data.configured ? 'ok' : 'warn') + '">' + (data.configured ? '可调度' : '还不能调度') + '</span>'
    + '<div class="base">' + esc(data.base_url) + '</div>'
    + '<div class="muted">上游失败时按优先级换下一把凭据。思考块会转成各协议自己的字段。</div></div>'
    + '<div class="muted">v' + esc(data.version) + '</div></div>'
    + '<div class="grid">'
    + stat('请求', num(stats.requests)) + stat('成功率', rate) + stat('平均首字', num(stats.avg_ttfb_ms) + ' ms')
    + stat('平均耗时', num(stats.avg_duration_ms) + ' ms') + stat('输出速度', num(stats.tokens_per_second) + ' tok/s')
    + stat('缓存命中', pct(stats.cache_hit_rate)) + '</div>'
    + '<div class="rows">'
    + row('凭据', (data.credentials || []).length + ' 把，' + healthyCount() + ' 把可用')
    + row('模型映射', (data.models || []).length + ' 条')
    + row('失败切换', data.failover ? '开启' : '关闭')
    + row('冷却基数', data.cooldown_seconds + ' 秒')
    + row('流空闲截断', (data.stream_idle_seconds || 0) + ' 秒') + '</div>';
}
function renderCredentials() {
  const editing = editingCredential ? (data.credentials || []).find(item => item.id === editingCredential) : null;
  if (editingCredential && !editing) editingCredential = '';
  const rows = (data.credentials || []).map(item => {
    const health = item.health || {};
    let status = '<span class="pill ok">正常</span>';
    if (item.disabled) status = '<span class="pill off">已停用</span>';
    else if (!item.has_key) status = '<span class="pill warn">缺少 key</span>';
    else if (health.cooldown_left_seconds > 0) status = '<span class="pill warn">冷却 ' + health.cooldown_left_seconds + 's</span>';
    const checked = health.last_checked ? '<div class="muted">检测于 ' + esc(timeText(health.last_checked)) + ' · ' + num(health.last_latency_ms) + ' ms</div>' : '';
    const error = health.last_error ? '<div class="err small" title="' + esc(health.last_error) + '">' + esc(clip(health.last_error, 90)) + '</div>' : '';
    return '<div class="row"><div class="grow"><b>' + esc(item.label) + '</b> <span class="muted">' + esc(item.id) + '</span>'
      + '<div class="muted">' + esc(item.mask || '无 key') + ' · 优先级 ' + item.priority + '</div>' + checked + error + '</div>'
      + '<div class="actions">' + status
      + ' <button data-probe="' + esc(item.id) + '">检测</button>'
      + ' <button data-edit="' + esc(item.id) + '">编辑</button>'
      + (health.cooldown_left_seconds > 0 ? ' <button data-reset="' + esc(item.id) + '">解除冷却</button>' : '')
      + ' <button data-toggle="' + esc(item.id) + '">' + (item.disabled ? '启用' : '停用') + '</button>'
      + ' <button class="danger" data-del="' + esc(item.id) + '">删除</button></div></div>';
  }).join('');
  document.querySelector('#credentials').innerHTML =
    '<div class="head"><b>凭据</b><span class="muted">优先级高的先用。401、403、429、5xx 和传输错误会切到下一把，并按冷却基数退避。</span></div>'
    + (rows || '<div class="muted">还没有凭据。添加一把 Mistral API key。</div>')
    + '<div class="form">'
    + '<input id="c-label" placeholder="名称" value="' + esc(editing ? editing.label : '') + '">'
    + '<input id="c-key" placeholder="' + (editing ? 'API key（留空不改）' : 'API key') + '" type="password">'
    + '<input id="c-pri" type="number" min="0" max="1000" value="' + (editing ? editing.priority : 10) + '" style="width:90px">'
    + '<button id="c-save">' + (editing ? '保存修改' : '添加') + '</button>'
    + (editing ? '<button id="c-cancel">取消</button>' : '')
    + '</div><div class="muted" style="margin-top:8px">添加后可以点「检测」调用 GET /models 验证 key。停用的凭据不会参与调度。</div>'
    + '<div class="form"><a class="button" href="https://chat.mistral.ai/code/extensions" target="_blank" rel="noopener noreferrer">创建 Vibe API key</a><span class="muted">在 Mistral 的 Code extensions 页面生成 key，再粘贴到上面的 API key 输入框。</span></div>';
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
      + (added ? ' <span class="pill off">已添加</span>' : '') + '</label>';
  }).join('');
  return '<div class="panel"><div class="row"><b>上游模型</b><span class="muted">' + upstreamModels.length + ' 个 · ' + esc(upstreamCredential) + '</span>'
    + '<div class="actions"><button id="up-all">全选</button><button id="up-add">添加选中</button><button id="up-close">关闭</button></div></div>'
    + '<input id="up-filter" class="wide" placeholder="过滤模型 ID">'
    + '<div id="upstream-list" class="checks">' + rows + '</div></div>';
}
function renderModels() {
  const existing = new Set((data.models || []).map(item => item.upstream_id));
  const rows = (data.models || []).map(item => {
    const test = modelTests[item.id];
    let result = '<span class="muted">未测试</span>';
    if (test && test.pending) result = '<span class="muted">测试中...</span>';
    else if (test && test.ok) result = '<span class="pill ok">' + test.status + '</span> <span class="muted">' + num(test.duration_ms) + ' ms</span>' + (test.reply ? '<div class="muted">' + esc(clip(test.reply, 60)) + '</div>' : '');
    else if (test) result = '<span class="pill warn">失败</span><div class="err small" title="' + esc(test.error || '') + '">' + esc(clip(test.error || '', 90)) + '</div>';
    return '<tr><td><b>' + esc(item.id) + '</b></td><td class="mono">' + esc(item.upstream_id) + '</td>'
      + '<td><select data-reason="' + esc(item.id) + '">' + reasonOptions(item.reasoning) + '</select></td>'
      + '<td>' + result + '</td>'
      + '<td class="actions"><button data-mtest="' + esc(item.id) + '">测试</button><button class="danger" data-mdel="' + esc(item.id) + '">删除</button></td></tr>';
  }).join('');
  document.querySelector('#models').innerHTML =
    '<div class="head"><b>模型映射</b><span class="muted">客户端模型名映射到 Mistral 上游模型 ID。reasoning 留空则透传客户端的 reasoning_effort。测试请求会计入上游账单。</span></div>'
    + '<div class="form"><select id="fetch-cred">' + credentialOptions() + '</select><button id="fetch-models">拉取上游模型</button><span class="muted">用选中的凭据调用 GET /models</span></div>'
    + upstreamPanel()
    + '<table><tr><th>客户端</th><th>上游</th><th>思考</th><th>最近测试</th><th></th></tr>'
    + (rows || '<tr><td colspan="5" class="muted">还没有模型映射。</td></tr>') + '</table>'
    + '<div class="form"><input id="m-id" placeholder="客户端模型名"><input id="m-up" placeholder="上游模型 ID"><select id="m-reason">' + reasonOptions('') + '</select><button id="m-add">手动添加</button></div>';
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
    msg.textContent = result.ok ? ('凭据可用：' + result.model_count + ' 个模型，' + result.duration_ms + ' ms') : ('检测失败：' + (result.error || 'unknown'));
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
    msg.textContent = result.ok ? ('模型可用：' + result.duration_ms + ' ms') : ('模型测试失败：' + (result.error || ''));
  } catch (error) { delete modelTests[id]; render(); msg.textContent = error.message; }
}
function renderConfig() {
  document.querySelector('#config').innerHTML =
    '<div class="head"><b>配置</b><span class="muted">保存后写入 CPA 配置文件的 plugins.configs.mistral 块。</span></div>'
    + '<div class="row"><label for="base">Base URL</label><input id="base" value="' + esc(data.base_url) + '" style="width:420px;max-width:100%"></div>'
    + '<div class="row"><label for="keep">日志保留条数</label><input id="keep" type="number" min="10" max="5000" value="' + data.log_retention + '" style="width:100px"></div>'
    + '<div class="row"><label for="idle">流空闲截断秒数，0 为关闭</label><input id="idle" type="number" min="0" max="3600" value="' + (data.stream_idle_seconds || 0) + '" style="width:100px"></div>'
    + '<div class="row"><label><input id="failover" type="checkbox"' + (data.failover ? ' checked' : '') + '> 上游失败时换下一把凭据</label><span class="muted">401、403、429、5xx、传输错误</span></div>'
    + '<div class="row"><label for="cooldown">冷却基数秒数</label><input id="cooldown" type="number" min="0" max="3600" value="' + data.cooldown_seconds + '" style="width:100px"></div>'
    + '<div class="form"><button id="save-settings">保存配置</button><span class="muted">failover 关闭时只使用当前排在最前的一把凭据。</span></div>';
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
  document.querySelector('#log-summary').textContent = '近 ' + lastLogs.length + ' 条：请求 ' + num(summary.requests)
    + ' · 成功率 ' + pct(summary.success_rate) + ' · 平均首字 ' + num(summary.avg_ttfb_ms) + ' ms'
    + ' · 平均耗时 ' + num(summary.avg_duration_ms) + ' ms · 输出 ' + num(summary.tokens_per_second) + ' tok/s'
    + ' · 缓存命中 ' + pct(summary.cache_hit_rate);
  const rows = lastLogs.filter(item => {
    if (logFilter === 'success') return item.status >= 200 && item.status < 400;
    if (logFilter === 'error') return item.status < 200 || item.status >= 400;
    return true;
  }).map(item => {
    const ok = item.status >= 200 && item.status < 400;
    return '<tr><td class="mono">' + esc(timeText(item.time)) + '</td><td>' + esc(item.model) + '</td><td>' + esc(item.credential) + '</td>'
      + '<td><span class="pill ' + (ok ? 'ok' : 'warn') + '">' + item.status + '</span>' + (item.probe ? ' <span class="muted">测试</span>' : '') + '</td>'
      + '<td>' + (item.attempts || 1) + '</td><td>' + num(item.ttfb_ms) + '</td><td>' + num(item.duration_ms) + '</td>'
      + '<td>' + (item.tokens_per_second ? num(item.tokens_per_second) : '') + '</td>'
      + '<td class="err small" title="' + esc(item.error || '') + '">' + esc(clip(item.error || '', 120)) + '</td></tr>';
  }).join('');
  document.querySelector('#log-table').innerHTML = '<table><tr><th>时间</th><th>模型</th><th>凭据</th><th>状态</th><th>尝试</th><th>首字 ms</th><th>耗时 ms</th><th>tok/s</th><th>错误</th></tr>'
    + (rows || '<tr><td colspan="9" class="muted">还没有请求记录。</td></tr>') + '</table>';
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
