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
.top-right { display:flex; align-items:center; gap:8px; flex-wrap:wrap; justify-content:flex-end; }
.version { font:12px var(--mono); color:var(--muted); border:1px solid var(--line); border-radius:999px; padding:4px 9px; white-space:nowrap; }
.lang-switch { display:inline-flex; align-items:center; gap:2px; padding:2px; border:1px solid var(--line); border-radius:999px; background:var(--soft); }
.lang-switch button { border:0; background:transparent; border-radius:999px; padding:3px 9px; font:600 11px var(--sans); color:var(--muted); }
.lang-switch button:hover { border:0; color:var(--ink); }
.lang-switch button.active { background:var(--ink); color:#fff; }
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
    <div><h1>Mistral</h1><p id="tagline"></p></div>
  </div>
  <div class="top-right">
    <div class="lang-switch" role="group" aria-label="Language">
      <button type="button" data-lang="zh" aria-pressed="true">中文</button>
      <button type="button" data-lang="en" aria-pressed="false">EN</button>
    </div>
    <span id="state" class="version">v0.0.0</span>
  </div>
</header>
<nav class="tabs" aria-label="Mistral">
  <button class="tab active" data-tab="overview" data-i18n="navOverview">概览</button>
  <button class="tab" data-tab="credentials" data-i18n="navCredentials">凭据</button>
  <button class="tab" data-tab="models" data-i18n="navModels">模型</button>
  <button class="tab" data-tab="logs" data-i18n="navLogs">日志</button>
  <button class="tab" data-tab="config" data-i18n="navConfig">配置</button>
</nav>
<main>
  <section id="overview" class="view"></section>
  <section id="credentials" class="view" hidden></section>
  <section id="models" class="view" hidden></section>
  <section id="logs" class="view" hidden>
    <div class="card">
      <div class="view-head">
        <div><h2 data-i18n="logsTitle">请求日志</h2><p data-i18n="logsHint">日志在插件数据目录的 requests.jsonl。测试请求标为 probe，不计入统计。</p></div>
        <div class="log-toolbar"><select id="log-filter"><option value="all" data-i18n="filterAll">全部</option><option value="success" data-i18n="filterSuccess">成功</option><option value="error" data-i18n="filterError">失败</option></select><select id="log-limit"><option value="50" data-i18n="limit50">50 条</option><option value="100" selected data-i18n="limit100">100 条</option><option value="200" data-i18n="limit200">200 条</option></select><button id="log-refresh" data-i18n="refresh">刷新</button></div>
      </div>
      <div id="log-summary" class="log-summary"></div>
      <div id="log-table" class="table-wrap"></div>
    </div>
  </section>
  <section id="config" class="view" hidden></section>
</main>
<p id="msg" role="status" aria-live="polite"></p>
<script>
const STRINGS = {
  zh: {
    tagline: 'CPA 上游插件 · 凭据、模型与请求账本',
    navOverview: '概览',
    navCredentials: '凭据',
    navModels: '模型',
    navLogs: '日志',
    navConfig: '配置',
    upstreamRoute: '上游路由',
    ready: '可调度',
    notReady: '还不能调度',
    activeCredential: '当前凭据 {label} · 上游失败按优先级切换',
    failoverNote: '上游失败按优先级切换凭据',
    setupNote: '还缺凭据或模型映射。先在凭据和模型页补齐。',
    requestPulse: '请求脉冲',
    pulseSummary: '最近 {n} 次 · 高度为耗时',
    pulseEmpty: '还没有真实请求。客户端接上后，这里会显示最近 30 次的耗时分布。',
    metricRequests: '请求',
    metricSuccess: '成功率',
    metricTtfb: '平均首字',
    metricDuration: '平均耗时',
    metricSpeed: '输出速度',
    metricCache: '输入缓存',
    tokenLedger: 'Token 账本',
    inputCacheHit: '输入缓存命中',
    tokenInput: '输入',
    tokenCached: '缓存输入',
    tokenOutput: '输出',
    noUsage: '还没有用量数据。',
    recentRequests: '最近请求',
    recentHint: '按时间倒序，测试请求标为 probe。',
    allLogs: '全部日志',
    colTime: '时间',
    colModel: '模型',
    colStatus: '状态',
    colDuration: '耗时',
    colInOut: '输入 → 输出',
    noRequests: '还没有请求',
    noRequestsHint: '客户端把模型名指到这里后，流量会出现在这张表里。',
    credentialPool: '凭据池',
    credentialsAvailable: '{healthy} / {total} 可用',
    manage: '管理',
    noCredentials: '还没有凭据',
    noCredentialsHint: '先添加一把 Mistral API key，调度才会开始。',
    runtime: '运行参数',
    settings: '配置',
    modelMappings: '模型映射',
    modelCount: '{n} 条',
    failover: '失败切换',
    failoverOn: '开启 · 冷却 {n}s',
    off: '关闭',
    idleCutoff: '流空闲截断',
    logRetention: '日志保留',
    entriesCount: '{n} 条',
    healthy: '正常',
    disabled: '已停用',
    missingKey: '缺少 key',
    noKey: '无 key',
    cooling: '冷却 {n}s',
    active: '当前',
    check: '检测',
    edit: '编辑',
    resetCooldown: '解除冷却',
    enable: '启用',
    disable: '停用',
    delete: '删除',
    checkedAt: '检测于 {time} · {latency}',
    priority: '优先级 {n}',
    credentialsTitle: '凭据',
    credentialsHint: '优先级高的先用。401、403、429、5xx 和传输错误会切到下一把，并按冷却基数退避。',
    credentialsEmptyHint: '添加一把 Mistral API key 后，调度和模型发现才能用。',
    editCredential: '编辑 {label}',
    addCredential: '添加凭据',
    keyStoredNote: 'API key 只保存在 CPA 配置里，页面只显示掩码。',
    namePlaceholder: '名称',
    keyKeepPlaceholder: 'API key（留空不改）',
    keyPlaceholder: 'API key',
    priorityLabel: '优先级',
    saveChanges: '保存修改',
    add: '添加',
    cancel: '取消',
    priorityHint: '优先级越大越先用。添加后可点「检测」调用 GET /models 验证 key。',
    createVibeKey: '创建 Vibe API key',
    vibeKeyHint: '在 Mistral 的 Code extensions 页面生成 key，再粘贴到上面的 API key 输入框。',
    confirmDeleteCredential: '删除凭据 {id}？',
    cooldownCleared: '已解除冷却',
    enabled: '已启用',
    disabledNote: '已停用',
    deleted: '已删除',
    newKeyRequired: '新凭据需要 API key。',
    saved: '已保存',
    added: '已添加',
    checkingCredential: '正在检测凭据...',
    credentialWorks: '凭据可用：{n} 个模型，{time}',
    checkFailed: '检测失败：{error}',
    modelsTitle: '模型映射',
    modelsHint: '客户端模型名映射到 Mistral 上游模型 ID。reasoning 留空则透传客户端的 reasoning_effort。测试请求会计入上游账单。',
    fetchUpstream: '拉取上游模型',
    fetchHint: '用选中的凭据调用 GET /models',
    upstreamModels: '上游模型',
    upstreamCount: '{n} 个 · {label}',
    selectAll: '全选',
    addSelected: '添加选中',
    close: '关闭',
    filterModels: '过滤模型 ID',
    alreadyAdded: '已添加',
    colClient: '客户端',
    colUpstream: '上游',
    colReasoning: '思考',
    colLastTest: '最近测试',
    notTested: '未测试',
    testing: '测试中...',
    failed: '失败',
    test: '测试',
    clientPlaceholder: '客户端模型名',
    upstreamPlaceholder: '上游模型 ID',
    addManually: '手动添加',
    confirmDeleteModel: '删除模型映射 {id}？',
    reasoningSaved: '思考档位已保存',
    modelAdded: '模型已添加',
    noModels: '还没有模型映射。',
    modelFieldsRequired: '模型名和上游 ID 都要填。',
    modelsAlreadyAdded: '选中的模型都已经添加。',
    selectModelsFirst: '先勾选要添加的模型。',
    fetchedModels: '拉取到 {n} 个模型。勾选后点「添加选中」。',
    fetchingModels: '正在拉取上游模型...',
    fetchFailed: '拉取失败',
    modelsAdded: '已添加 {n} 个模型',
    modelWorks: '模型可用：{time}',
    modelTestFailed: '模型测试失败：{error}',
    passthrough: '透传',
    noAvailableCredential: '没有可用凭据',
    configTitle: '配置',
    configHint: '保存后写入 CPA 配置文件的 plugins.configs.mistral 块。',
    baseUrl: 'Base URL',
    baseUrlHint: 'Mistral API 地址，带 /v1。',
    logRetentionHint: '10 到 5000。',
    idleHint: '0 为关闭，单位秒。',
    cooldown: '冷却基数',
    cooldownHint: '0 表示不冷却，单位秒。',
    failoverHint: '401、403、429、5xx、传输错误时换下一把。',
    on: '开启',
    saveSettings: '保存配置',
    failoverOffHint: 'failover 关闭时只使用当前排在最前的一把凭据。',
    settingsSaved: '配置已保存',
    logsTitle: '请求日志',
    logsHint: '日志在插件数据目录的 requests.jsonl。测试请求标为 probe，不计入统计。',
    filterAll: '全部',
    filterSuccess: '成功',
    filterError: '失败',
    limit50: '50 条',
    limit100: '100 条',
    limit200: '200 条',
    refresh: '刷新',
    sumEntries: '记录',
    sumRequests: '请求',
    sumSuccess: '成功率',
    sumTtfb: '平均首字',
    sumDuration: '平均耗时',
    sumInput: '输入',
    sumOutput: '输出',
    sumCached: '缓存输入',
    sumCacheRate: '缓存率',
    colCredential: '凭据',
    colAttempts: '尝试',
    colTtfb: '首字',
    colInput: '输入',
    colOutput: '输出',
    colCached: '缓存输入',
    colCacheRate: '缓存率',
    colError: '错误',
    noLogs: '还没有请求记录。',
    probeTag: '测试',
    missingSession: '没有读到管理页的登录状态。请从 CPA 管理页的 Mistral 菜单打开。'
  },
  en: {
    tagline: 'CPA upstream plugin · credentials, models, and request ledger',
    navOverview: 'Overview',
    navCredentials: 'Credentials',
    navModels: 'Models',
    navLogs: 'Logs',
    navConfig: 'Settings',
    upstreamRoute: 'Upstream route',
    ready: 'Ready',
    notReady: 'Not ready',
    activeCredential: 'Active: {label} · failover follows priority',
    failoverNote: 'Failover follows credential priority',
    setupNote: 'A credential or model mapping is missing. Add them under Credentials and Models.',
    requestPulse: 'Request pulse',
    pulseSummary: 'Last {n} · bar height is duration',
    pulseEmpty: 'No real requests yet. Once a client connects, the last 30 durations appear here.',
    metricRequests: 'Requests',
    metricSuccess: 'Success',
    metricTtfb: 'Avg first byte',
    metricDuration: 'Avg duration',
    metricSpeed: 'Output speed',
    metricCache: 'Input cache',
    tokenLedger: 'Token ledger',
    inputCacheHit: 'input cache hit',
    tokenInput: 'Input',
    tokenCached: 'Cached input',
    tokenOutput: 'Output',
    noUsage: 'No usage data yet.',
    recentRequests: 'Recent requests',
    recentHint: 'Newest first. Test requests are marked probe.',
    allLogs: 'All logs',
    colTime: 'Time',
    colModel: 'Model',
    colStatus: 'Status',
    colDuration: 'Duration',
    colInOut: 'Input → output',
    noRequests: 'No requests yet',
    noRequestsHint: 'Point a client at this provider and traffic will appear here.',
    credentialPool: 'Credential pool',
    credentialsAvailable: '{healthy} of {total} available',
    manage: 'Manage',
    noCredentials: 'No credentials yet',
    noCredentialsHint: 'Add a Mistral API key to start scheduling.',
    runtime: 'Runtime',
    settings: 'Settings',
    modelMappings: 'Model mappings',
    modelCount: '{n} mappings',
    failover: 'Failover',
    failoverOn: 'On · {n}s cooldown',
    off: 'Off',
    idleCutoff: 'Idle cutoff',
    logRetention: 'Log retention',
    entriesCount: '{n} entries',
    healthy: 'Healthy',
    disabled: 'Disabled',
    missingKey: 'Missing key',
    noKey: 'No key',
    cooling: 'Cooling {n}s',
    active: 'Active',
    check: 'Check',
    edit: 'Edit',
    resetCooldown: 'Reset cooldown',
    enable: 'Enable',
    disable: 'Disable',
    delete: 'Delete',
    checkedAt: 'Checked {time} · {latency}',
    priority: 'priority {n}',
    credentialsTitle: 'Credentials',
    credentialsHint: 'Higher priority goes first. 401, 403, 429, 5xx, and transport errors move to the next credential and back off by the cooldown base.',
    credentialsEmptyHint: 'Add a Mistral API key to enable scheduling and model discovery.',
    editCredential: 'Edit {label}',
    addCredential: 'Add credential',
    keyStoredNote: 'The API key stays in the CPA config. This page only shows a mask.',
    namePlaceholder: 'Name',
    keyKeepPlaceholder: 'API key (leave blank to keep)',
    keyPlaceholder: 'API key',
    priorityLabel: 'Priority',
    saveChanges: 'Save changes',
    add: 'Add',
    cancel: 'Cancel',
    priorityHint: 'Higher priority is tried first. Use Check to call GET /models and verify the key.',
    createVibeKey: 'Create Vibe API key',
    vibeKeyHint: 'Generate a key on the Mistral Code extensions page, then paste it into the API key field above.',
    confirmDeleteCredential: 'Delete credential {id}?',
    cooldownCleared: 'Cooldown cleared',
    enabled: 'Enabled',
    disabledNote: 'Disabled',
    deleted: 'Deleted',
    newKeyRequired: 'A new credential needs an API key.',
    saved: 'Saved',
    added: 'Added',
    checkingCredential: 'Checking credential...',
    credentialWorks: 'Credential works: {n} models, {time}',
    checkFailed: 'Check failed: {error}',
    modelsTitle: 'Model mappings',
    modelsHint: 'Client model names map to Mistral upstream IDs. Empty reasoning passes through the client reasoning_effort. Test requests count toward the upstream bill.',
    fetchUpstream: 'Fetch upstream models',
    fetchHint: 'Calls GET /models with the selected credential',
    upstreamModels: 'Upstream models',
    upstreamCount: '{n} models · {label}',
    selectAll: 'Select all',
    addSelected: 'Add selected',
    close: 'Close',
    filterModels: 'Filter model IDs',
    alreadyAdded: 'Added',
    colClient: 'Client',
    colUpstream: 'Upstream',
    colReasoning: 'Reasoning',
    colLastTest: 'Last test',
    notTested: 'Not tested',
    testing: 'Testing...',
    failed: 'Failed',
    test: 'Test',
    clientPlaceholder: 'Client model name',
    upstreamPlaceholder: 'Upstream model ID',
    addManually: 'Add manually',
    confirmDeleteModel: 'Delete model mapping {id}?',
    reasoningSaved: 'Reasoning saved',
    modelAdded: 'Model added',
    noModels: 'No model mappings yet.',
    modelFieldsRequired: 'Model name and upstream ID are both required.',
    modelsAlreadyAdded: 'The selected models are already added.',
    selectModelsFirst: 'Select at least one model.',
    fetchedModels: 'Fetched {n} models. Select the ones to add.',
    fetchingModels: 'Fetching upstream models...',
    fetchFailed: 'Fetch failed',
    modelsAdded: 'Added {n} models',
    modelWorks: 'Model works: {time}',
    modelTestFailed: 'Model test failed: {error}',
    passthrough: 'Passthrough',
    noAvailableCredential: 'No available credential',
    configTitle: 'Settings',
    configHint: 'Saving writes the plugins.configs.mistral block in the CPA config file.',
    baseUrl: 'Base URL',
    baseUrlHint: 'Mistral API address, including /v1.',
    logRetentionHint: '10 to 5000.',
    idleHint: 'Seconds. 0 disables it.',
    cooldown: 'Cooldown base',
    cooldownHint: 'Seconds. 0 disables cooldown.',
    failoverHint: 'Move to the next credential on 401, 403, 429, 5xx, or transport errors.',
    on: 'On',
    saveSettings: 'Save settings',
    failoverOffHint: 'When failover is off, only the first credential in the current order is used.',
    settingsSaved: 'Settings saved',
    logsTitle: 'Request logs',
    logsHint: 'Logs live in requests.jsonl under the plugin data directory. Test requests are marked probe and excluded from statistics.',
    filterAll: 'All',
    filterSuccess: 'Success',
    filterError: 'Error',
    limit50: '50 entries',
    limit100: '100 entries',
    limit200: '200 entries',
    refresh: 'Refresh',
    sumEntries: 'Entries',
    sumRequests: 'Requests',
    sumSuccess: 'Success',
    sumTtfb: 'Avg first byte',
    sumDuration: 'Avg duration',
    sumInput: 'Input',
    sumOutput: 'Output',
    sumCached: 'Cached input',
    sumCacheRate: 'Cache rate',
    colCredential: 'Credential',
    colAttempts: 'Attempts',
    colTtfb: 'First byte',
    colInput: 'Input',
    colOutput: 'Output',
    colCached: 'Cached input',
    colCacheRate: 'Cache rate',
    colError: 'Error',
    noLogs: 'No requests logged yet.',
    probeTag: 'Test',
    missingSession: 'Could not read the management session. Open this page from the Mistral menu in CPA.'
  }
};
const MANUAL_LANGUAGE_KEY = 'mistral-console-language';
const msg = document.querySelector('#msg');
let data = null;
let managementKey = readPanelKey() || readParentKey();
let lang = detectLanguage();
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
function storedLanguage(raw) {
  if (!raw) return '';
  let value = raw;
  try {
    const parsed = JSON.parse(raw);
    value = parsed && parsed.state && parsed.state.language ? parsed.state.language : (parsed && parsed.language ? parsed.language : parsed);
  } catch (error) {}
  return typeof value === 'string' ? value : '';
}
function cpaLanguage() {
  let value = '';
  try { value = storedLanguage(localStorage.getItem('cli-proxy-language')); } catch (error) {}
  const lowered = String(value || '').toLowerCase();
  if (lowered.indexOf('zh') === 0) return 'zh';
  if (lowered.indexOf('en') === 0) return 'en';
  try {
    if (window.parent && window.parent !== window) {
      const parentLang = String(window.parent.document.documentElement.lang || '').toLowerCase();
      if (parentLang.indexOf('zh') === 0) return 'zh';
      if (parentLang.indexOf('en') === 0) return 'en';
    }
  } catch (error) {}
  return 'zh';
}
function manualLanguage() {
  try {
    const value = localStorage.getItem(MANUAL_LANGUAGE_KEY);
    return value === 'zh' || value === 'en' ? value : '';
  } catch (error) { return ''; }
}
function detectLanguage() {
  return manualLanguage() || cpaLanguage();
}
function t(key, vars) {
  const table = STRINGS[lang] || STRINGS.zh;
  let text = table[key];
  if (text === undefined) text = STRINGS.zh[key];
  if (text === undefined) return key;
  if (vars) Object.keys(vars).forEach(name => { text = text.split('{' + name + '}').join(vars[name]); });
  return text;
}
function applyStaticI18n() {
  document.documentElement.lang = lang === 'zh' ? 'zh-CN' : 'en';
  document.querySelector('#tagline').textContent = t('tagline');
  document.querySelectorAll('[data-i18n]').forEach(node => { node.textContent = t(node.dataset.i18n); });
  document.querySelectorAll('[data-lang]').forEach(node => {
    const active = node.dataset.lang === lang;
    node.classList.toggle('active', active);
    node.setAttribute('aria-pressed', active ? 'true' : 'false');
  });
}
function setLanguage(next) {
  if ((next !== 'zh' && next !== 'en') || next === lang) return;
  lang = next;
  try { localStorage.setItem(MANUAL_LANGUAGE_KEY, next); } catch (error) {}
  applyStaticI18n();
  if (data) render();
  if (!document.querySelector('#logs').hidden) renderLogs();
}
function switchTab(name) {
  document.querySelectorAll('.tab').forEach(item => item.classList.toggle('active', item.dataset.tab === name));
  document.querySelectorAll('main section').forEach(section => section.hidden = section.id !== name);
  msg.textContent = '';
  if (name === 'logs') loadLogs();
}
document.querySelectorAll('.tab').forEach(tab => tab.onclick = () => switchTab(tab.dataset.tab));
document.querySelectorAll('[data-lang]').forEach(button => button.onclick = () => setLanguage(button.dataset.lang));
window.addEventListener('storage', event => {
  if (event.key !== 'cli-proxy-language' || manualLanguage()) return;
  const next = cpaLanguage();
  if (next === lang) return;
  lang = next;
  applyStaticI18n();
  if (data) render();
  if (!document.querySelector('#logs').hidden) renderLogs();
});
async function api(body) {
  if (!managementKey) {
    msg.textContent = t('missingSession');
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
function num(value) { return Number(value || 0).toLocaleString(lang === 'en' ? 'en-US' : 'zh-CN'); }
function pct(value) { return Math.round(Number(value || 0) * 100) + '%'; }
function pct1(value) { return (Number(value || 0) * 100).toFixed(1) + '%'; }
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
  if (item.disabled) return {cls:'off', text:t('disabled')};
  if (!item.has_key) return {cls:'warn', text:t('missingKey')};
  if (health.cooldown_left_seconds > 0) return {cls:'warn', text:t('cooling', {n:health.cooldown_left_seconds})};
  return {cls:'ok', text:t('healthy')};
}
function statusBadge(item) {
  const status = credentialStatus(item);
  return '<span class="badge ' + status.cls + '">' + status.text + '</span>';
}
function credentialOptions() {
  const items = (data.credentials || []).filter(item => item.has_key);
  if (!items.length) return '<option value="">' + t('noAvailableCredential') + '</option>';
  if (!items.some(item => item.id === fetchCredential)) fetchCredential = items[0].id;
  return items.map(item => '<option value="' + esc(item.id) + '"' + (item.id === fetchCredential ? ' selected' : '') + '>' + esc(item.label) + (item.disabled ? ' (' + t('disabled') + ')' : '') + '</option>').join('');
}
function reasonOptions(current) {
  return ['','low','medium','high','xhigh','max'].map(value => '<option value="' + value + '"' + (value === current ? ' selected' : '') + '>' + (value || t('passthrough')) + '</option>').join('');
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
    : '<div class="pulse-empty">' + t('pulseEmpty') + '</div>';
  const cacheRate = has && stats.cache_hit_rate !== undefined ? Math.max(0, Math.min(100, Number(stats.cache_hit_rate) * 100)) : 0;
  const tokenHtml = has
    ? '<div class="token-box"><div class="token-top"><span class="eyebrow">' + t('tokenLedger') + '</span><div class="token-rate"><b>' + (stats.cache_hit_rate !== undefined ? pct1(stats.cache_hit_rate) : '—') + '</b><span>' + t('inputCacheHit') + '</span></div></div>'
      + '<div class="token-bar"><i style="width:' + cacheRate.toFixed(1) + '%"></i></div>'
      + '<div class="token-legend"><div><span>' + t('tokenInput') + '</span><b>' + compact(stats.prompt_tokens) + '</b></div><div><span>' + t('tokenCached') + '</span><b>' + compact(stats.cached_tokens) + '</b></div><div><span>' + t('tokenOutput') + '</span><b>' + compact(stats.completion_tokens) + '</b></div></div></div>'
    : '<div class="token-box"><span class="eyebrow">' + t('tokenLedger') + '</span><div class="empty" style="padding:10px 0 0">' + t('noUsage') + '</div></div>';
  const credItems = creds.map(item => {
    const health = item.health || {};
    const status = credentialStatus(item);
    return '<div class="cred"><span class="dot ' + status.cls + '"></span><div class="cred-main">'
      + '<div class="cred-title"><b>' + esc(item.label) + '</b>' + (item.id === activeId ? '<span class="chip">' + t('active') + '</span>' : '') + '</div>'
      + '<div class="cred-meta">' + esc(item.mask || t('noKey')) + ' · ' + t('priority', {n:item.priority}) + '</div>'
      + (health.last_checked ? '<div class="cred-meta">' + t('checkedAt', {time:esc(timeText(health.last_checked)), latency:sec(health.last_latency_ms)}) + '</div>' : '')
      + (health.last_error ? '<div class="err" title="' + esc(health.last_error) + '">' + esc(clip(health.last_error, 70)) + '</div>' : '')
      + '</div>' + statusBadge(item) + '</div>';
  }).join('');
  const credentialCard = '<div class="card"><div class="rail-head"><div><div class="eyebrow">' + t('credentialPool') + '</div><span>' + t('credentialsAvailable', {healthy:healthy, total:creds.length}) + '</span></div><button class="linklike" id="ov-manage-creds">' + t('manage') + '</button></div>'
    + (credItems || '<div class="empty"><b>' + t('noCredentials') + '</b>' + t('noCredentialsHint') + '</div>') + '</div>';
  const policyCard = '<div class="card"><div class="rail-head"><div><div class="eyebrow">' + t('runtime') + '</div></div><button class="linklike" id="ov-config">' + t('settings') + '</button></div><div class="facts">'
    + fact(t('modelMappings'), t('modelCount', {n:models.length}))
    + fact(t('failover'), data.failover ? t('failoverOn', {n:data.cooldown_seconds}) : t('off'))
    + fact(t('idleCutoff'), (data.stream_idle_seconds || 0) > 0 ? data.stream_idle_seconds + 's' : t('off'))
    + fact(t('logRetention'), t('entriesCount', {n:data.log_retention}))
    + '</div></div>';
  const recentRows = recent.slice(0, 8).map(item => {
    const ok = item.status >= 200 && item.status < 400;
    return '<tr><td class="mono nowrap">' + esc(shortTime(item.time)) + '</td>'
      + '<td class="model-cell" title="' + esc(item.model || '') + '">' + esc(item.model) + '</td>'
      + '<td class="nowrap"><span class="badge ' + (ok ? 'ok' : 'bad') + '">' + item.status + '</span>' + (item.probe ? ' <span class="badge probe">' + t('probeTag') + '</span>' : '') + '</td>'
      + '<td class="num">' + sec(item.duration_ms) + '</td>'
      + '<td class="num">' + compactTokens(item.prompt_tokens) + ' → ' + compactTokens(item.completion_tokens) + '</td></tr>';
  }).join('');
  const recentCard = '<div class="card"><div class="view-head" style="margin-bottom:8px"><div><h2>' + t('recentRequests') + '</h2><p>' + t('recentHint') + '</p></div><button id="ov-logs">' + t('allLogs') + '</button></div>'
    + (recentRows
      ? '<div class="table-wrap"><table><thead><tr><th>' + t('colTime') + '</th><th>' + t('colModel') + '</th><th>' + t('colStatus') + '</th><th class="num">' + t('colDuration') + '</th><th class="num">' + t('colInOut') + '</th></tr></thead><tbody>' + recentRows + '</tbody></table></div>'
      : '<div class="empty"><b>' + t('noRequests') + '</b>' + t('noRequestsHint') + '</div>')
    + '</div>';
  const routeNote = data.configured
    ? (data.active_label ? t('activeCredential', {label:esc(data.active_label)}) : t('failoverNote'))
    : t('setupNote');
  document.querySelector('#overview').innerHTML = '<div class="overview"><div class="stack">'
    + '<div class="card"><div class="route-head"><div><span class="eyebrow">' + t('upstreamRoute') + '</span><div class="route-url">' + esc(data.base_url) + '</div><div class="route-note">' + routeNote + '</div></div>'
    + '<span class="state ' + (data.configured ? 'ok' : 'warn') + '"><i></i>' + (data.configured ? t('ready') : t('notReady')) + '</span></div>'
    + '<div class="pulse-box"><div class="pulse-head"><b>' + t('requestPulse') + '</b><span>' + t('pulseSummary', {n:pulseItems.length}) + '</span></div>' + pulse + '</div>'
    + '<div class="ledger">'
    + metric(t('metricRequests'), num(stats.requests))
    + metric(t('metricSuccess'), has ? pct(stats.success_rate) : '—')
    + metric(t('metricTtfb'), sec(stats.avg_ttfb_ms))
    + metric(t('metricDuration'), sec(stats.avg_duration_ms))
    + metric(t('metricSpeed'), has ? Number(stats.tokens_per_second || 0).toFixed(1) : '—', 'tok/s')
    + metric(t('metricCache'), has && stats.cache_hit_rate !== undefined ? pct1(stats.cache_hit_rate) : '—')
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
    const checked = health.last_checked ? '<div class="cred-meta">' + t('checkedAt', {time:esc(timeText(health.last_checked)), latency:sec(health.last_latency_ms)}) + '</div>' : '';
    const error = health.last_error ? '<div class="err" title="' + esc(health.last_error) + '">' + esc(clip(health.last_error, 90)) + '</div>' : '';
    return '<div class="cred-card' + (item.id === editingCredential ? ' editing' : '') + '"><div class="grow"><div class="cred-title"><b>' + esc(item.label) + '</b> <span class="mono muted">' + esc(item.id) + '</span>'
      + statusBadge(item) + (item.id === data.active ? '<span class="chip">' + t('active') + '</span>' : '') + '</div>'
      + '<div class="cred-meta">' + esc(item.mask || t('noKey')) + ' · ' + t('priority', {n:item.priority}) + '</div>' + checked + error + '</div>'
      + '<div class="actions"><button data-probe="' + esc(item.id) + '">' + t('check') + '</button>'
      + '<button data-edit="' + esc(item.id) + '">' + t('edit') + '</button>'
      + (health.cooldown_left_seconds > 0 ? '<button data-reset="' + esc(item.id) + '">' + t('resetCooldown') + '</button>' : '')
      + '<button data-toggle="' + esc(item.id) + '">' + (item.disabled ? t('enable') : t('disable')) + '</button>'
      + '<button class="danger" data-del="' + esc(item.id) + '">' + t('delete') + '</button></div></div>';
  }).join('');
  const editor = '<div class="panel"><div class="row-head"><b>' + (editing ? t('editCredential', {label:esc(editing.label)}) : t('addCredential')) + '</b><span class="muted">' + t('keyStoredNote') + '</span></div>'
    + '<div class="form" style="margin-top:12px;padding-top:0;border-top:0">'
    + '<input id="c-label" placeholder="' + t('namePlaceholder') + '" value="' + esc(editing ? editing.label : '') + '">'
    + '<input id="c-key" placeholder="' + (editing ? t('keyKeepPlaceholder') : t('keyPlaceholder')) + '" type="password">'
    + '<input id="c-pri" class="pri" type="number" min="0" max="1000" aria-label="' + t('priorityLabel') + '" title="' + t('priorityLabel') + '" value="' + (editing ? editing.priority : 10) + '">'
    + '<button class="primary" id="c-save">' + (editing ? t('saveChanges') : t('add')) + '</button>'
    + (editing ? '<button id="c-cancel">' + t('cancel') + '</button>' : '') + '</div>'
    + '<div class="muted" style="margin-top:8px">' + t('priorityHint') + '</div></div>';
  document.querySelector('#credentials').innerHTML =
    '<div class="card"><div class="view-head"><div><h2>' + t('credentialsTitle') + '</h2><p>' + t('credentialsHint') + '</p></div><span class="muted">' + t('credentialsAvailable', {healthy:healthyCount(), total:(data.credentials || []).length}) + '</span></div>'
    + (rows || '<div class="empty"><b>' + t('noCredentials') + '</b>' + t('credentialsEmptyHint') + '</div>')
    + editor
    + '<div class="form"><a class="button" href="https://chat.mistral.ai/code/extensions" target="_blank" rel="noopener noreferrer">' + t('createVibeKey') + '</a><span class="muted">' + t('vibeKeyHint') + '</span></div></div>';
  document.querySelector('#c-save').onclick = saveCredential;
  const cancel = document.querySelector('#c-cancel');
  if (cancel) cancel.onclick = () => { editingCredential = ''; render(); };
  document.querySelectorAll('[data-probe]').forEach(btn => btn.onclick = () => probeCredential(btn.dataset.probe));
  document.querySelectorAll('[data-edit]').forEach(btn => btn.onclick = () => { editingCredential = btn.dataset.edit; render(); });
  document.querySelectorAll('[data-reset]').forEach(btn => btn.onclick = () => applyAction({action:'reset-cooldown', id:btn.dataset.reset}, t('cooldownCleared')));
  document.querySelectorAll('[data-toggle]').forEach(btn => btn.onclick = () => {
    const item = (data.credentials || []).find(entry => entry.id === btn.dataset.toggle);
    if (!item) return;
    applyAction({action:'put-credential', id:item.id, disabled:!item.disabled}, item.disabled ? t('enabled') : t('disabledNote'));
  });
  document.querySelectorAll('[data-del]').forEach(btn => btn.onclick = () => {
    if (!confirm(t('confirmDeleteCredential', {id:btn.dataset.del}))) return;
    applyAction({action:'delete-credential', id:btn.dataset.del}, t('deleted'));
  });
}
async function saveCredential() {
  const label = val('#c-label');
  const key = val('#c-key');
  const priority = Number(val('#c-pri') || 10);
  if (!editingCredential && !key) { msg.textContent = t('newKeyRequired'); return; }
  const action = {action:'put-credential', label:label, priority:priority};
  if (editingCredential) action.id = editingCredential;
  if (key) action.api_key = key;
  const ok = await applyAction(action, editingCredential ? t('saved') : t('added'));
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
      + (added ? ' <span class="badge off">' + t('alreadyAdded') + '</span>' : '') + '</label>';
  }).join('');
  return '<div class="panel"><div class="row-head"><div><b>' + t('upstreamModels') + '</b> <span class="muted">' + t('upstreamCount', {n:upstreamModels.length, label:esc(upstreamCredential)}) + '</span></div>'
    + '<div class="actions"><button id="up-all">' + t('selectAll') + '</button><button class="primary" id="up-add">' + t('addSelected') + '</button><button id="up-close">' + t('close') + '</button></div></div>'
    + '<input id="up-filter" class="wide" placeholder="' + t('filterModels') + '">'
    + '<div id="upstream-list" class="checks">' + rows + '</div></div>';
}
function renderModels() {
  const rows = (data.models || []).map(item => {
    const test = modelTests[item.id];
    let result = '<span class="muted">' + t('notTested') + '</span>';
    if (test && test.pending) result = '<span class="muted">' + t('testing') + '</span>';
    else if (test && test.ok) result = '<span class="badge ok">' + test.status + '</span> <span class="muted">' + sec(test.duration_ms) + '</span>' + (test.reply ? '<div class="muted">' + esc(clip(test.reply, 60)) + '</div>' : '');
    else if (test) result = '<span class="badge bad">' + t('failed') + '</span><div class="err" title="' + esc(test.error || '') + '">' + esc(clip(test.error || '', 90)) + '</div>';
    return '<tr><td><b>' + esc(item.id) + '</b></td><td class="mono">' + esc(item.upstream_id) + '</td>'
      + '<td><select data-reason="' + esc(item.id) + '">' + reasonOptions(item.reasoning) + '</select></td>'
      + '<td>' + result + '</td>'
      + '<td><div class="actions"><button data-mtest="' + esc(item.id) + '">' + t('test') + '</button><button class="danger" data-mdel="' + esc(item.id) + '">' + t('delete') + '</button></div></td></tr>';
  }).join('');
  document.querySelector('#models').innerHTML =
    '<div class="card"><div class="view-head"><div><h2>' + t('modelsTitle') + '</h2><p>' + t('modelsHint') + '</p></div><span class="muted">' + t('modelCount', {n:(data.models || []).length}) + '</span></div>'
    + '<div class="form" style="margin-top:0;padding-top:0;border-top:0"><select id="fetch-cred">' + credentialOptions() + '</select><button id="fetch-models">' + t('fetchUpstream') + '</button><span class="muted">' + t('fetchHint') + '</span></div>'
    + upstreamPanel()
    + '<div class="table-wrap"><table><thead><tr><th>' + t('colClient') + '</th><th>' + t('colUpstream') + '</th><th>' + t('colReasoning') + '</th><th>' + t('colLastTest') + '</th><th></th></tr></thead><tbody>'
    + (rows || '<tr><td colspan="5" class="muted">' + t('noModels') + '</td></tr>') + '</tbody></table></div>'
    + '<div class="form"><input id="m-id" placeholder="' + t('clientPlaceholder') + '"><input id="m-up" placeholder="' + t('upstreamPlaceholder') + '"><select id="m-reason">' + reasonOptions('') + '</select><button class="primary" id="m-add">' + t('addManually') + '</button></div></div>';
  const selector = document.querySelector('#fetch-cred');
  selector.onchange = () => { fetchCredential = selector.value; };
  document.querySelector('#fetch-models').onclick = fetchUpstream;
  document.querySelectorAll('[data-mtest]').forEach(btn => btn.onclick = () => probeModel(btn.dataset.mtest));
  document.querySelectorAll('[data-mdel]').forEach(btn => btn.onclick = () => {
    if (!confirm(t('confirmDeleteModel', {id:btn.dataset.mdel}))) return;
    delete modelTests[btn.dataset.mdel];
    applyAction({action:'delete-model', id:btn.dataset.mdel}, t('deleted'));
  });
  document.querySelectorAll('[data-reason]').forEach(select => select.onchange = () => {
    const item = (data.models || []).find(entry => entry.id === select.dataset.reason);
    if (!item) return;
    applyAction({action:'put-model', id:item.id, upstream_id:item.upstream_id, reasoning:select.value}, t('reasoningSaved'));
  });
  document.querySelector('#m-add').onclick = () => {
    const id = val('#m-id');
    const upstream = val('#m-up');
    if (!id || !upstream) { msg.textContent = t('modelFieldsRequired'); return; }
    applyAction({action:'put-model', id:id, upstream_id:upstream, reasoning:val('#m-reason')}, t('modelAdded'));
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
  msg.textContent = t('fetchingModels');
  try {
    const result = await api({action:'fetch-models', credential_id:fetchCredential});
    if (!result.ok) {
      upstreamModels = [];
      upstreamError = result.error || t('fetchFailed');
      render();
      msg.textContent = upstreamError;
      return;
    }
    upstreamModels = result.models || [];
    upstreamCredential = result.credential_label || '';
    render();
    msg.textContent = t('fetchedModels', {n:upstreamModels.length});
  } catch (error) { msg.textContent = error.message; }
}
async function addSelectedModels() {
  const checked = Array.prototype.slice.call(document.querySelectorAll('#upstream-list input:checked'));
  if (!checked.length) { msg.textContent = t('selectModelsFirst'); return; }
  const existing = new Set((data.models || []).map(item => item.upstream_id));
  const items = checked.filter(input => !existing.has(input.value)).map(input => ({id:input.value, upstream_id:input.value}));
  if (!items.length) { msg.textContent = t('modelsAlreadyAdded'); return; }
  await applyAction({action:'put-models', items:items}, t('modelsAdded', {n:items.length}));
}
async function probeCredential(id) {
  msg.textContent = t('checkingCredential');
  try {
    const result = await api({action:'test-credential', id:id});
    data = await api({action:'state'});
    render();
    msg.textContent = result.ok ? t('credentialWorks', {n:result.model_count, time:sec(result.duration_ms)}) : t('checkFailed', {error:result.error || 'unknown'});
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
    msg.textContent = result.ok ? t('modelWorks', {time:sec(result.duration_ms)}) : t('modelTestFailed', {error:result.error || ''});
  } catch (error) { delete modelTests[id]; render(); msg.textContent = error.message; }
}
function renderConfig() {
  document.querySelector('#config').innerHTML =
    '<div class="card"><div class="view-head"><div><h2>' + t('configTitle') + '</h2><p>' + t('configHint') + '</p></div></div>'
    + '<div class="config-grid">'
    + '<div class="field full"><label for="base">' + t('baseUrl') + '<span class="hint">' + t('baseUrlHint') + '</span></label><input id="base" value="' + esc(data.base_url) + '"></div>'
    + '<div class="field"><label for="keep">' + t('logRetention') + '<span class="hint">' + t('logRetentionHint') + '</span></label><input id="keep" type="number" min="10" max="5000" value="' + data.log_retention + '"></div>'
    + '<div class="field"><label for="idle">' + t('idleCutoff') + '<span class="hint">' + t('idleHint') + '</span></label><input id="idle" type="number" min="0" max="3600" value="' + (data.stream_idle_seconds || 0) + '"></div>'
    + '<div class="field"><label for="cooldown">' + t('cooldown') + '<span class="hint">' + t('cooldownHint') + '</span></label><input id="cooldown" type="number" min="0" max="3600" value="' + data.cooldown_seconds + '"></div>'
    + '<div class="field"><label>' + t('failover') + '<span class="hint">' + t('failoverHint') + '</span></label><span class="switch"><input id="failover" type="checkbox"' + (data.failover ? ' checked' : '') + '> ' + t('on') + '</span></div>'
    + '</div><div class="form"><button class="primary" id="save-settings">' + t('saveSettings') + '</button><span class="muted">' + t('failoverOffHint') + '</span></div></div>';
  document.querySelector('#save-settings').onclick = () => applyAction({
    action:'put-settings', base_url:val('#base'), log_retention:Number(val('#keep')),
    stream_idle_seconds:Number(val('#idle')), failover:document.querySelector('#failover').checked,
    cooldown_seconds:Number(val('#cooldown'))
  }, t('settingsSaved'));
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
    '<span>' + t('sumEntries') + '<b>' + lastLogs.length + '</b></span>'
    + '<span>' + t('sumRequests') + '<b>' + num(summary.requests) + '</b></span>'
    + '<span>' + t('sumSuccess') + '<b>' + (summary.success_rate !== undefined ? pct(summary.success_rate) : '—') + '</b></span>'
    + '<span>' + t('sumTtfb') + '<b>' + sec(summary.avg_ttfb_ms) + '</b></span>'
    + '<span>' + t('sumDuration') + '<b>' + sec(summary.avg_duration_ms) + '</b></span>'
    + '<span>' + t('sumInput') + '<b>' + compact(summary.prompt_tokens) + '</b></span>'
    + '<span>' + t('sumOutput') + '<b>' + compact(summary.completion_tokens) + '</b></span>'
    + '<span>' + t('sumCached') + '<b>' + compact(summary.cached_tokens) + '</b></span>'
    + '<span>' + t('sumCacheRate') + '<b>' + (summary.cache_hit_rate !== undefined ? pct1(summary.cache_hit_rate) : '—') + '</b></span>';
  const rows = lastLogs.filter(item => {
    if (logFilter === 'success') return item.status >= 200 && item.status < 400;
    if (logFilter === 'error') return item.status < 200 || item.status >= 400;
    return true;
  }).map(item => {
    const ok = item.status >= 200 && item.status < 400;
    const input = Number(item.prompt_tokens || 0);
    const cached = Number(item.cached_tokens || 0);
    const rate = input > 0 ? cached / input : null;
    return '<tr><td class="mono nowrap">' + esc(timeText(item.time)) + '</td>'
      + '<td class="model-cell" title="' + esc(item.model || '') + '">' + esc(item.model) + '</td>'
      + '<td class="nowrap" title="' + esc(item.credential || '') + '">' + esc(item.credential) + '</td>'
      + '<td class="nowrap"><span class="badge ' + (ok ? 'ok' : 'bad') + '">' + item.status + '</span>' + (item.probe ? ' <span class="badge probe">' + t('probeTag') + '</span>' : '') + '</td>'
      + '<td class="num">' + (item.attempts || 1) + '</td>'
      + '<td class="num">' + sec(item.ttfb_ms) + '</td>'
      + '<td class="num">' + sec(item.duration_ms) + '</td>'
      + '<td class="num">' + tokens(item.prompt_tokens) + '</td>'
      + '<td class="num">' + tokens(item.completion_tokens) + '</td>'
      + '<td class="num">' + tokens(item.cached_tokens) + '</td>'
      + '<td class="num">' + (rate === null ? '—' : pct1(rate)) + '</td>'
      + '<td class="num">' + (Number(item.tokens_per_second || 0) > 0 ? Number(item.tokens_per_second).toFixed(1) : '—') + '</td>'
      + '<td class="error-cell err" title="' + esc(item.error || '') + '">' + esc(clip(item.error || '', 80)) + '</td></tr>';
  }).join('');
  const headers = '<tr><th>' + t('colTime') + '</th><th>' + t('colModel') + '</th><th>' + t('colCredential') + '</th><th>' + t('colStatus') + '</th><th class="num">' + t('colAttempts') + '</th><th class="num">' + t('colTtfb') + '</th><th class="num">' + t('colDuration') + '</th><th class="num">' + t('colInput') + '</th><th class="num">' + t('colOutput') + '</th><th class="num">' + t('colCached') + '</th><th class="num">' + t('colCacheRate') + '</th><th class="num">tok/s</th><th>' + t('colError') + '</th></tr>';
  document.querySelector('#log-table').innerHTML = '<table><thead>' + headers + '</thead><tbody>'
    + (rows || '<tr><td colspan="13" class="muted">' + t('noLogs') + '</td></tr>') + '</tbody></table>';
}
document.querySelector('#log-refresh').onclick = loadLogs;
document.querySelector('#log-filter').onchange = event => { logFilter = event.target.value; renderLogs(); };
document.querySelector('#log-limit').onchange = loadLogs;
function val(sel) { const node = document.querySelector(sel); return node ? node.value.trim() : ''; }
async function refresh() {
  data = await api({action:'state'});
  render();
}
applyStaticI18n();
refresh().catch(error => { msg.textContent = error.message; });
</script>`
}
