package main

func consolePage() string {
	return `<!doctype html>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Mistral</title>
<style>
:root {
  color-scheme: light;
  --bg:#f6f7f9;
  --surface:#fff;
  --surface-2:#fafbfc;
  --soft:#f1f3f6;
  --ink:#0f172a;
  --ink-2:#334155;
  --muted:#64748b;
  --faint:#94a3b8;
  --line:#e5e9ef;
  --line-2:#eef1f5;
  --accent:#6366f1;
  --accent-2:#8b5cf6;
  --accent-soft:#eef2ff;
  --accent-ink:#4f46e5;
  --ok:#059669;
  --ok-soft:#ecfdf5;
  --warn:#d97706;
  --warn-soft:#fffbeb;
  --bad:#dc2626;
  --bad-soft:#fef2f2;
  --shadow-sm:0 1px 2px rgba(15,23,42,.04);
  --shadow-md:0 4px 16px rgba(15,23,42,.06);
  --shadow-lg:0 16px 48px rgba(15,23,42,.12);
  --radius:12px;
  --radius-sm:8px;
  --mono:ui-monospace,SFMono-Regular,"SF Mono",Menlo,Consolas,monospace;
  --sans:-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC","Hiragino Sans GB","Microsoft YaHei",sans-serif;
}
* { box-sizing:border-box; }
[hidden] { display:none !important; }
html, body { margin:0; }
body { font:14px/1.55 var(--sans); background:var(--bg); color:var(--ink); -webkit-font-smoothing:antialiased; }
button, select, input { font:inherit; }
button { border:1px solid var(--line); background:var(--surface); color:var(--ink); border-radius:var(--radius-sm); padding:7px 12px; cursor:pointer; transition:border-color .15s, background .15s, color .15s, box-shadow .15s; }
button:hover { border-color:#cbd5e1; background:var(--surface-2); }
button:disabled { opacity:.45; cursor:not-allowed; }
button.primary { background:var(--accent); border-color:var(--accent); color:#fff; font-weight:600; box-shadow:0 1px 2px rgba(99,102,241,.24); }
button.primary:hover { background:var(--accent-ink); border-color:var(--accent-ink); }
button.danger { color:var(--bad); }
button.danger:hover { border-color:#fecaca; background:var(--bad-soft); }
button.ghost { border-color:transparent; background:transparent; color:var(--muted); }
button.ghost:hover { background:var(--soft); color:var(--ink); }
a { color:var(--accent-ink); text-underline-offset:2px; }
a.button { display:inline-flex; align-items:center; gap:6px; border:1px solid var(--line); background:var(--surface); border-radius:var(--radius-sm); padding:7px 12px; color:var(--ink); text-decoration:none; }
a.button:hover { border-color:#cbd5e1; }
:focus-visible { outline:2px solid var(--accent); outline-offset:2px; }
::selection { background:var(--accent-soft); color:var(--accent-ink); }
* { scrollbar-width:thin; scrollbar-color:#cbd5e1 transparent; }
*::-webkit-scrollbar { width:8px; height:8px; }
*::-webkit-scrollbar-thumb { background:#cbd5e1; border-radius:999px; }
*::-webkit-scrollbar-track { background:transparent; }

/* header */
.top { display:flex; justify-content:space-between; align-items:center; gap:16px; padding:14px 24px; background:var(--surface); border-bottom:1px solid var(--line); }
.brand { display:flex; align-items:center; gap:12px; min-width:0; }
.mark { width:38px; height:38px; border:1px solid var(--line); border-radius:11px; background:linear-gradient(135deg,var(--accent-soft),#f5f3ff); color:var(--accent-ink); display:grid; place-items:center; flex:none; box-shadow:var(--shadow-sm); }
.mark svg { width:21px; height:21px; display:block; }
h1 { margin:0; font-size:18px; line-height:1.2; letter-spacing:-.02em; }
.brand p { margin:2px 0 0; color:var(--muted); font-size:12px; }
.top-right { display:flex; align-items:center; gap:8px; flex-wrap:wrap; justify-content:flex-end; }
.route-state { display:inline-flex; align-items:center; gap:7px; border-radius:999px; border:1px solid; padding:4px 10px; font-size:11.5px; font-weight:600; white-space:nowrap; }
.route-state i { width:6px; height:6px; border-radius:50%; background:currentColor; }
.route-state.ok { color:var(--ok); border-color:#bbf7d0; background:var(--ok-soft); }
.route-state.warn { color:var(--warn); border-color:#fde68a; background:var(--warn-soft); }
.version { font:12px var(--mono); color:var(--muted); border:1px solid var(--line); border-radius:999px; padding:4px 9px; white-space:nowrap; }
.lang-switch { display:inline-flex; align-items:center; gap:2px; padding:2px; border:1px solid var(--line); border-radius:999px; background:var(--soft); }
.lang-switch button { border:0; background:transparent; border-radius:999px; padding:3px 9px; font:600 11px var(--sans); color:var(--muted); }
.lang-switch button:hover { border:0; background:transparent; color:var(--ink); }
.lang-switch button.active { background:var(--ink); color:#fff; }

/* tabs */
.tabs { display:flex; gap:2px; padding:0 24px; background:var(--surface); border-bottom:1px solid var(--line); overflow-x:auto; }
.tab { border:0; border-bottom:2px solid transparent; border-radius:0; background:transparent; color:var(--muted); padding:11px 14px; white-space:nowrap; font-weight:550; }
.tab:hover { color:var(--ink); border-color:var(--line); background:transparent; }
.tab.active { color:var(--accent-ink); border-bottom-color:var(--accent); }

/* filter bar */
.filter-bar { max-width:1440px; margin:0 auto; padding:16px 24px 0; }
.filter-row { display:flex; gap:8px; align-items:center; flex-wrap:wrap; }
.filter-row + .filter-row { margin-top:8px; }
.filter-right { margin-left:auto; display:flex; gap:8px; align-items:center; flex-wrap:wrap; }
.seg { display:inline-flex; align-items:center; gap:2px; padding:2px; border:1px solid var(--line); border-radius:var(--radius-sm); background:var(--surface); box-shadow:var(--shadow-sm); }
.seg button { border:0; background:transparent; border-radius:6px; padding:5px 10px; font-size:12px; font-weight:600; color:var(--muted); }
.seg button:hover { background:var(--soft); color:var(--ink); }
.seg button.active { background:var(--ink); color:#fff; }
.seg.soft button.active { background:var(--accent-soft); color:var(--accent-ink); }
.seg-pill { display:inline-flex; align-items:center; gap:2px; padding:3px; border-radius:8px; background:#f1f3f5; }
.seg-pill button { border:0; background:transparent; border-radius:6px; padding:4px 14px; font-size:12px; font-weight:500; color:#64748b; cursor:pointer; transition:all .15s ease; }
.seg-pill button:hover { color:#0f172a; }
.seg-pill button.active { background:#fff; color:#0f172a; font-weight:500; border:1px solid rgba(0,0,0,.06); box-shadow:0 1px 2px rgba(15,23,42,.06); }
.cost-chip { display:inline-flex; align-items:baseline; gap:6px; border:1px solid #bbf7d0; background:var(--ok-soft); color:var(--ok); border-radius:999px; padding:5px 11px; font:600 12.5px var(--mono); white-space:nowrap; }
.cost-chip span { font:500 11px var(--sans); color:#047857; }
.filter-bar select, .filter-bar input { padding:7px 10px; border:1px solid var(--line); border-radius:var(--radius-sm); background:var(--surface); color:var(--ink); box-shadow:var(--shadow-sm); }
.filter-bar input#search { flex:1 1 220px; max-width:340px; }
.filter-note { margin:9px 0 0; color:var(--faint); font-size:11.5px; }

/* layout */
main { max-width:1440px; margin:0 auto; padding:16px 24px 64px; }
h2 { margin:0; font-size:16px; letter-spacing:-.015em; }
h3 { margin:0; font-size:14px; letter-spacing:-.01em; }
.view-head { display:flex; justify-content:space-between; align-items:flex-start; gap:16px; margin-bottom:14px; }
.view-head p { margin:4px 0 0; color:var(--muted); font-size:12.5px; max-width:760px; }
.card { position:relative; background:var(--surface); border:1px solid var(--line); border-radius:var(--radius); padding:18px; box-shadow:var(--shadow-sm); }
.card + .card { margin-top:16px; }
.card-head { display:flex; justify-content:space-between; align-items:flex-start; gap:12px; flex-wrap:wrap; margin-bottom:14px; }
.card-head p { margin:3px 0 0; color:var(--muted); font-size:11.5px; }

/* kpi */
.kpi-row { display:grid; grid-template-columns:repeat(5,minmax(0,1fr)); gap:12px; margin-bottom:16px; }
.kpi { background:var(--surface); border:1px solid var(--line); border-radius:var(--radius); padding:15px 16px; box-shadow:var(--shadow-sm); min-width:0; }
.kpi-label { display:flex; align-items:center; justify-content:space-between; gap:8px; color:var(--muted); font-size:11.5px; font-weight:600; letter-spacing:.02em; }
.kpi-label svg { width:15px; height:15px; color:var(--faint); flex:none; }
.kpi-value { margin-top:8px; font:650 26px/1.05 var(--mono); letter-spacing:-.04em; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.kpi-value.good { color:var(--ok); }
.kpi-value.bad { color:var(--bad); }
.kpi-value small { font:600 12px var(--sans); color:var(--muted); letter-spacing:0; margin-left:4px; }
.kpi-sub { margin-top:7px; color:var(--muted); font-size:11.5px; line-height:1.5; }
.kpi-sub b { color:var(--ink-2); font-weight:600; }
.kpi-sub .bad { color:var(--bad); }
.kpi-bar { height:6px; border-radius:999px; background:var(--soft); overflow:hidden; margin-top:9px; }
.kpi-bar i { display:block; height:100%; border-radius:999px; background:linear-gradient(90deg,var(--accent),var(--accent-2)); transition:width .4s ease; }

/* chart */
.trend-body { display:grid; grid-template-columns:minmax(0,1fr) 320px; align-items:stretch; }
.chart-col { min-width:0; padding-right:20px; }
.chart-metric-row { margin-bottom:10px; }
.chart-wrap { position:relative; width:100%; min-height:250px; }
.chart-wrap svg { display:block; width:100%; height:250px; overflow:visible; }
.chart-grid-line { stroke:var(--line-2); stroke-width:1; }
.chart-axis-label { fill:var(--faint); font:10.5px var(--mono); }
.chart-line { fill:none; stroke:var(--accent); stroke-width:2; stroke-linecap:round; stroke-linejoin:round; }
.chart-area { fill:url(#chartGradient); }
.chart-dot { fill:var(--surface); stroke:var(--accent); stroke-width:2; }
.chart-cross { stroke:var(--faint); stroke-width:1; stroke-dasharray:3 3; }
.chart-tip { position:absolute; z-index:5; pointer-events:none; background:var(--ink); color:#fff; border-radius:8px; padding:8px 10px; font-size:11.5px; line-height:1.5; box-shadow:var(--shadow-lg); min-width:120px; transform:translate(-50%,-100%); }
.chart-tip b { font-family:var(--mono); font-weight:600; }
.chart-tip .tip-time { color:#cbd5e1; font-size:10.5px; margin-bottom:3px; }
.chart-tip .tip-row { display:flex; justify-content:space-between; gap:12px; }
#hm-tip { white-space:nowrap; }
.chart-empty { display:flex; align-items:center; justify-content:center; height:250px; color:var(--muted); font-size:12.5px; border:1px dashed var(--line); border-radius:var(--radius-sm); }

/* breakdown */
.breakdown-panel { min-width:0; min-height:0; display:flex; flex-direction:column; border-left:1px solid var(--line); padding-left:20px; }
.breakdown-head { display:flex; align-items:baseline; justify-content:space-between; gap:10px; margin-bottom:8px; }
.breakdown-head h4 { margin:0; font-size:12.5px; font-weight:600; color:var(--ink-2); }
.breakdown-head p { margin:0; color:var(--muted); font-size:11px; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.breakdown-list { flex:1; min-height:0; max-height:250px; display:flex; flex-direction:column; gap:2px; overflow:auto; margin:0 -6px; padding:0 6px; overscroll-behavior:contain; }
.breakdown-item { display:flex; align-items:flex-start; gap:9px; padding:9px 6px; border-radius:var(--radius-sm); cursor:pointer; transition:background .12s; }
.breakdown-item:hover { background:var(--soft); }
.breakdown-item.active { background:var(--accent-soft); }
.breakdown-dot { width:8px; height:8px; border-radius:50%; margin-top:5px; flex:none; }
.breakdown-main { min-width:0; flex:1; }
.breakdown-name { font-weight:600; font-size:12.5px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.breakdown-meta { color:var(--muted); font-size:10.5px; margin-top:2px; line-height:1.45; }
.breakdown-meta .bad { color:var(--bad); }
.breakdown-value { font:600 12.5px var(--mono); white-space:nowrap; text-align:right; }
.breakdown-value small { display:block; font:500 10px var(--sans); color:var(--muted); }

/* heatmap */
.heatmap { display:flex; gap:8px; overflow-x:auto; padding-bottom:6px; }
.hm-weekdays { display:grid; grid-template-rows:repeat(7,var(--hm-size,12px)); gap:var(--hm-gap,3px); font-size:11px; color:#94a3b8; padding-top:20px; flex:none; text-align:right; padding-right:4px; }
.hm-weekdays span { line-height:var(--hm-size,12px); white-space:nowrap; }
.hm-main { min-width:0; flex:1; }
.hm-months { position:relative; height:18px; font-size:11px; color:#94a3b8; white-space:nowrap; margin-bottom:2px; }
.hm-months span { position:absolute; top:0; }
.hm-grid { display:grid; grid-auto-flow:column; grid-template-rows:repeat(7,var(--hm-size,12px)); grid-auto-columns:var(--hm-size,12px); gap:var(--hm-gap,3px); width:100%; justify-content:space-between; }
.hm-cell { width:var(--hm-size,12px); height:var(--hm-size,12px); border-radius:var(--hm-radius,3px); background:#f1f3f6; transition:outline .1s, filter .1s; cursor:pointer; }
.hm-cell.blank { background:transparent; cursor:default; }
.hm-cell.blank:hover { outline:0; filter:none; }
.hm-cell:hover { outline:1.5px solid #6366f1; outline-offset:1px; filter:brightness(1.05); }
.hm-cell.l1 { background:#c4caf9; }
.hm-cell.l2 { background:#939cf5; }
.hm-cell.l3 { background:#5e66eb; }
.hm-cell.l4 { background:#4338ca; }
.hm-legend { display:flex; align-items:center; gap:4px; justify-content:flex-end; margin-top:14px; color:#94a3b8; font-size:11px; }
.hm-legend i { width:12px; height:12px; border-radius:3px; display:inline-block; cursor:default; }
.hm-legend i.hm-cell { width:12px; height:12px; border-radius:3px; }

/* tables */
.table-wrap { overflow-x:auto; }
table { width:100%; border-collapse:collapse; }
th { font:600 10.5px var(--mono); letter-spacing:.04em; color:var(--muted); text-align:left; padding:9px 10px; border-bottom:1px solid var(--line); white-space:nowrap; background:var(--surface-2); }
td { padding:9px 10px; border-bottom:1px solid var(--line-2); font-size:12.5px; vertical-align:middle; }
tbody tr:last-child td { border-bottom:0; }
tbody tr:hover td { background:var(--surface-2); }
.num { text-align:right; font-family:var(--mono); font-variant-numeric:tabular-nums; white-space:nowrap; }
.mono { font-family:var(--mono); font-size:11.5px; }
.muted { color:var(--muted); font-size:12px; }
.nowrap { white-space:nowrap; }
.grow { min-width:0; flex:1; }
.model-cell { max-width:190px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.error-cell { max-width:220px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.actions { display:flex; gap:6px; align-items:center; flex-wrap:wrap; justify-content:flex-end; }

/* badges */
.badge { display:inline-flex; align-items:center; border-radius:999px; padding:2px 8px; font-size:10.5px; font-weight:600; border:1px solid; white-space:nowrap; }
.badge.ok { color:var(--ok); border-color:#bbf7d0; background:var(--ok-soft); }
.badge.warn { color:var(--warn); border-color:#fde68a; background:var(--warn-soft); }
.badge.bad { color:var(--bad); border-color:#fecaca; background:var(--bad-soft); }
.badge.off { color:var(--muted); border-color:var(--line); background:var(--soft); }
.badge.probe { color:var(--muted); border-color:var(--line); background:transparent; }
.badge.indigo { color:var(--accent-ink); border-color:#c7d2fe; background:var(--accent-soft); }
.chip { display:inline-flex; align-items:center; border:1px solid #c7d2fe; background:var(--accent-soft); color:var(--accent-ink); border-radius:999px; padding:1px 7px; font-size:10.5px; font-weight:600; }
.dot { width:8px; height:8px; border-radius:50%; margin-top:5px; flex:none; background:var(--ok); box-shadow:0 0 0 3px var(--ok-soft); }
.dot.warn { background:var(--warn); box-shadow:0 0 0 3px var(--warn-soft); }
.dot.off { background:#a8aea4; box-shadow:0 0 0 3px #f1f2ef; }

/* credentials / models */
.cred { display:flex; gap:10px; padding:12px 0; border-bottom:1px solid var(--line-2); align-items:flex-start; }
.cred:last-of-type { border-bottom:0; }
.cred-main { min-width:0; flex:1; }
.cred-title { display:flex; gap:7px; align-items:center; flex-wrap:wrap; }
.cred-title b { font-weight:600; }
.cred-meta { color:var(--muted); font-size:11.5px; margin-top:2px; overflow-wrap:anywhere; }
.cred-card { display:flex; justify-content:space-between; gap:16px; align-items:flex-start; padding:15px 0; border-top:1px solid var(--line-2); }
.cred-card:first-of-type { border-top:0; }
.cred-card.editing { background:var(--accent-soft); margin:0 -18px; padding:15px 18px; }
.cred-card .cred-meta { font-size:12px; }
.err { color:var(--bad); font-size:11.5px; overflow-wrap:anywhere; }
.empty { padding:22px 0; color:var(--muted); font-size:12.5px; }
.empty b { display:block; color:var(--ink); font-size:13px; margin-bottom:2px; }
.facts { margin-top:2px; }
.fact { display:flex; justify-content:space-between; gap:14px; padding:9px 0; border-bottom:1px solid var(--line-2); font-size:13px; }
.fact:last-child { border-bottom:0; padding-bottom:0; }
.fact span { color:var(--muted); }
.fact b { font-weight:600; text-align:right; }
.form { display:flex; gap:8px; align-items:center; flex-wrap:wrap; margin-top:16px; padding-top:16px; border-top:1px solid var(--line); }
input, select { padding:8px 10px; border:1px solid var(--line); border-radius:var(--radius-sm); background:var(--surface); color:var(--ink); min-width:0; }
input.wide { width:100%; margin-top:10px; }
input.pri { width:92px; }
select { max-width:100%; }
.panel { border:1px solid var(--line); border-radius:var(--radius-sm); padding:14px; margin:16px 0; background:var(--soft); }
.panel .row-head { display:flex; justify-content:space-between; gap:12px; align-items:center; flex-wrap:wrap; }
.checks { display:grid; grid-template-columns:repeat(3,minmax(0,1fr)); gap:3px 12px; max-height:320px; overflow:auto; margin-top:10px; }
.check { display:flex; gap:6px; align-items:center; padding:4px 2px; font-size:12.5px; min-width:0; }
.check.disabled { opacity:.55; }
.check .mono { overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.config-grid { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:0 28px; }
.field { display:flex; justify-content:space-between; gap:16px; align-items:center; padding:13px 0; border-top:1px solid var(--line-2); }
.field.full { grid-column:1/-1; flex-direction:column; align-items:stretch; }
.field.full input { width:100%; }
.field label { font-size:13px; }
.field .hint { color:var(--muted); font-size:11.5px; display:block; margin-top:2px; }
.field input[type=number] { width:110px; }
.field input[type=text], .field input:not([type]) { width:min(420px,100%); }
.switch { display:flex; align-items:center; gap:8px; }
.linklike { border:0; background:transparent; padding:2px; color:var(--accent-ink); font-weight:600; font-size:12px; }
.linklike:hover { color:var(--ink); border:0; background:transparent; }
.eyebrow { font:600 10.5px var(--mono); letter-spacing:.1em; text-transform:uppercase; color:var(--muted); }
.route-head { display:flex; justify-content:space-between; gap:16px; align-items:flex-start; }
.route-head > div:first-child { min-width:0; }
.route-url { font:600 15px/1.45 var(--mono); word-break:break-all; margin:6px 0 4px; }
.route-note { color:var(--muted); font-size:12.5px; }
.pulse-box { margin-top:16px; padding-top:12px; border-top:1px solid var(--line-2); }
.pulse-head { display:flex; justify-content:space-between; gap:12px; align-items:baseline; color:var(--muted); font-size:11.5px; }
.pulse-head b { color:var(--ink); font-weight:600; font-size:12px; }
.pulse { display:flex; align-items:flex-end; gap:3px; height:42px; margin-top:8px; }
.tick { flex:1 1 0; min-width:3px; border-radius:2px 2px 0 0; background:var(--accent); opacity:.75; }
.tick.bad { background:var(--bad); }
.pulse-empty { display:flex; align-items:center; height:42px; margin-top:8px; color:var(--muted); font-size:12px; border-bottom:1px dashed var(--line); }
.ledger { display:grid; grid-template-columns:repeat(3,minmax(0,1fr)); gap:1px; background:var(--line); border:1px solid var(--line); border-radius:11px; overflow:hidden; margin-top:16px; }
.metric { background:var(--surface); padding:12px 14px; min-width:0; }
.metric span { display:block; color:var(--muted); font-size:11px; letter-spacing:.04em; }
.metric b { display:block; font:600 21px/1.25 var(--mono); letter-spacing:-.035em; margin-top:3px; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
.metric small { font:500 11px var(--sans); color:var(--muted); letter-spacing:0; margin-left:5px; }
.rail-head { display:flex; justify-content:space-between; align-items:center; gap:12px; padding-bottom:12px; border-bottom:1px solid var(--line-2); }
.rail-head span { color:var(--muted); font-size:12px; }

/* logs */
.log-summary { display:flex; flex-wrap:wrap; gap:8px 18px; padding:12px 0 4px; color:var(--muted); font-size:12px; }
.log-summary b { color:var(--ink); font-family:var(--mono); font-weight:600; margin-left:3px; }
.log-summary .bad b { color:var(--bad); }
.log-toolbar { display:flex; gap:8px; align-items:center; flex-wrap:wrap; justify-content:flex-end; }
.log-filter-row { display:flex; gap:8px; align-items:center; flex-wrap:wrap; margin-bottom:10px; }
.log-row { cursor:pointer; }
.log-row.expanded td { background:var(--accent-soft); }
.log-detail td { padding:0; background:var(--surface-2); border-bottom:1px solid var(--line); }
.detail-grid { display:grid; grid-template-columns:repeat(4,minmax(0,1fr)); gap:12px; padding:16px; }
.detail-card { background:var(--surface); border:1px solid var(--line); border-radius:var(--radius-sm); padding:12px 14px; }
.detail-title { font:600 11px var(--mono); letter-spacing:.05em; text-transform:uppercase; color:var(--muted); margin-bottom:8px; }
.detail-facts { display:flex; flex-direction:column; gap:7px; }
.detail-fact { display:flex; justify-content:space-between; gap:10px; font-size:11.5px; align-items:baseline; }
.detail-fact span { color:var(--muted); white-space:nowrap; }
.detail-fact b { font-family:var(--mono); font-weight:600; text-align:right; overflow-wrap:anywhere; }
.detail-actions { display:flex; gap:6px; flex-wrap:wrap; padding:0 16px 16px; }
.detail-error { margin:0 16px 16px; border:1px solid #fecaca; background:var(--bad-soft); border-radius:var(--radius-sm); padding:12px 14px; }
.detail-error .detail-title { color:var(--bad); }
.detail-error pre { margin:0; white-space:pre-wrap; word-break:break-word; font:11.5px/1.55 var(--mono); color:#7f1d1d; }
.detail-error .hint { margin-top:8px; font-size:11.5px; color:#991b1b; }
.chevron { display:inline-block; transition:transform .15s; color:var(--faint); font-size:10px; }
.log-row.expanded .chevron { transform:rotate(90deg); }

/* toast */
#msg { position:fixed; left:50%; bottom:22px; transform:translateX(-50%); margin:0; padding:9px 15px; border-radius:999px; background:var(--ink); color:#fff; font-size:12.5px; box-shadow:var(--shadow-lg); max-width:min(80vw,640px); z-index:40; }
#msg:empty { display:none; }

@media (max-width:1180px) { .kpi-row { grid-template-columns:repeat(3,minmax(0,1fr)); } .detail-grid { grid-template-columns:repeat(2,minmax(0,1fr)); } }
@media (max-width:1000px) {
  .trend-body { grid-template-columns:1fr; }
  .chart-col { padding-right:0; }
  .breakdown-panel { border-left:0; border-top:1px solid var(--line); padding-left:0; padding-top:14px; margin-top:16px; }
  .checks { grid-template-columns:repeat(2,minmax(0,1fr)); }
}
@media (max-width:720px) {
  main { padding:16px 14px 56px; }
  .filter-bar { padding:14px 14px 0; }
  .top { padding:14px; }
  .tabs { padding:0 10px; }
  .kpi-row { grid-template-columns:repeat(2,minmax(0,1fr)); }
  .ledger { grid-template-columns:repeat(2,minmax(0,1fr)); }
  .config-grid { grid-template-columns:1fr; }
  .field.full { grid-column:auto; }
  .checks { grid-template-columns:1fr; }
  .route-head { flex-direction:column; }
  .route-url { font-size:13px; }
  .cred-card { flex-direction:column; }
  .actions { justify-content:flex-start; }
  .card { padding:16px; }
  .detail-grid { grid-template-columns:1fr; }
  .filter-bar input#search { max-width:none; }
}
@media (prefers-reduced-motion:reduce) { * { transition:none !important; } }
</style>
<header class="top">
  <div class="brand">
    <span class="mark" aria-hidden="true"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M3 8h10.5a2.5 2.5 0 1 0-2.4-3.2"/><path d="M3 12h15a2.5 2.5 0 1 1-2.4 3.2"/><path d="M3 16h7"/></svg></span>
    <div><h1>Mistral</h1><p id="tagline"></p></div>
  </div>
  <div class="top-right">
    <span id="route-state" class="route-state warn"><i></i><span id="route-state-text"></span></span>
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
<div id="filter-bar" class="filter-bar">
  <div class="filter-row">
    <div class="seg" id="time-seg" role="group" aria-label="Time range">
      <button data-time="today" class="active" data-i18n="timeToday">今天</button>
      <button data-time="24h" data-i18n="time24h">24 小时</button>
      <button data-time="7d" data-i18n="time7d">7 天</button>
      <button data-time="30d" data-i18n="time30d">30 天</button>
      <button data-time="all" data-i18n="timeAll">全部</button>
    </div>
    <div class="filter-right">
      <span id="cost-chip" class="cost-chip">≈$0.00 <span data-i18n="liveCost">生效费用</span></span>
      <select id="auto-refresh" aria-label="Auto refresh">
        <option value="0" data-i18n="autoOff">手动</option>
        <option value="5">5s</option>
        <option value="10">10s</option>
        <option value="30">30s</option>
      </select>
      <button id="refresh-btn" class="ghost" title="Refresh" data-i18n="refresh">刷新</button>
      <button id="export-btn" data-i18n="exportCsv">导出 CSV</button>
    </div>
  </div>
  <div class="filter-row">
    <select id="model-filter" aria-label="Model"></select>
    <select id="cred-filter" aria-label="Credential"></select>
    <div class="seg soft" id="status-seg" role="group" aria-label="Status">
      <button data-status="all" class="active" data-i18n="filterAll">全部</button>
      <button data-status="success" data-i18n="filterSuccess">成功</button>
      <button data-status="error" data-i18n="filterError">失败</button>
      <button data-status="retry" data-i18n="filterRetry">重试</button>
      <button data-status="probe" data-i18n="filterProbe">测试</button>
    </div>
    <input id="search" type="search" placeholder="搜索模型、凭据、错误、请求 ID" data-i18n-placeholder="searchPlaceholder">
  </div>
  <p class="filter-note" data-i18n="statsNote">测试请求不计入汇总；费用按 Mistral 公开价估算。</p>
</div>
<main>
  <section id="overview" class="view">
    <div id="kpi-row" class="kpi-row"></div>
    <div class="card">
      <div class="card-head">
        <div><h3 data-i18n="chartTitle">用量趋势</h3><p data-i18n="chartHint">按当前时间范围和筛选条件聚合</p></div>
        <div class="seg-pill" id="chart-split-seg">
          <button data-split="model" class="active" data-i18n="chartSplitModel">模型</button>
          <button data-split="credential" data-i18n="chartSplitCredential">凭据</button>
          <button data-split="status" data-i18n="chartSplitStatus">状态</button>
        </div>
      </div>
      <div class="trend-body">
        <div class="chart-col">
          <div class="seg-pill chart-metric-row" id="chart-metric-seg">
            <button data-metric="tokens" data-i18n="chartTokens">Token</button>
            <button data-metric="cost" class="active" data-i18n="chartCost">费用</button>
            <button data-metric="requests" data-i18n="chartRequests">请求</button>
            <button data-metric="speed" data-i18n="chartSpeed">速度</button>
          </div>
          <div id="chart-wrap" class="chart-wrap"></div>
        </div>
        <aside class="breakdown-panel">
          <div class="breakdown-head"><h4 data-i18n="breakdownTitle">分拆</h4><p id="breakdown-sub"></p></div>
          <div id="breakdown-list" class="breakdown-list"></div>
        </aside>
      </div>
    </div>
    <div class="card">
      <div class="card-head">
        <div><h3 data-i18n="heatmapTitle">用量热力图</h3><p data-i18n="heatmapHint">最近 53 周的每日用量</p></div>
        <div class="seg-pill" id="heatmap-seg">
          <button data-hm="tokens" data-i18n="heatmapTokens">Token</button>
          <button data-hm="requests" data-i18n="heatmapRequests">请求</button>
          <button data-hm="cost" class="active" data-i18n="heatmapCost">费用</button>
        </div>
      </div>
      <div id="heatmap" class="heatmap"></div>
      <div id="hm-tip" class="chart-tip" hidden></div>
      <div class="hm-legend"><span data-i18n="heatmapLess">少</span><i class="hm-cell l0"></i><i class="hm-cell l1"></i><i class="hm-cell l2"></i><i class="hm-cell l3"></i><i class="hm-cell l4"></i><span data-i18n="heatmapMore">多</span></div>
    </div>
    <div class="card">
      <div class="card-head">
        <div><h3 data-i18n="recentRequests">最近请求</h3><p data-i18n="recentHint">按时间倒序，测试请求标为 probe。</p></div>
        <button id="ov-logs" class="linklike" data-i18n="allLogs">全部日志</button>
      </div>
      <div id="recent-table" class="table-wrap"></div>
    </div>
  </section>
  <section id="credentials" class="view" hidden></section>
  <section id="models" class="view" hidden></section>
  <section id="logs" class="view" hidden>
    <div class="card">
      <div class="card-head">
        <div><h3 data-i18n="logsTitle">请求日志</h3><p data-i18n="logsHint">日志在插件数据目录的 requests.jsonl。测试请求标为 probe，不计入统计。</p></div>
        <div class="log-toolbar">
          <select id="log-limit"><option value="50">50</option><option value="100" selected>100</option><option value="200">200</option><option value="500">500</option></select>
          <button id="log-refresh" data-i18n="refresh">刷新</button>
        </div>
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
    timeToday: '今天',
    time24h: '24 小时',
    time7d: '7 天',
    time30d: '30 天',
    timeAll: '全部',
    autoOff: '手动',
    liveCost: '生效费用',
    exportCsv: '导出 CSV',
    filterAll: '全部',
    filterSuccess: '成功',
    filterError: '失败',
    filterRetry: '重试',
    filterProbe: '测试',
    searchPlaceholder: '搜索模型、凭据、错误、请求 ID',
    statsNote: '测试请求不计入汇总；费用按 Mistral 公开价估算。',
    filterModel: '全部模型',
    filterCredential: '全部凭据',
    chartTitle: '用量趋势',
    chartHint: '按当前时间范围和筛选条件聚合',
    chartTokens: 'Token',
    chartCost: '费用',
    chartRequests: '请求',
    chartSpeed: '速度',
    chartSplitModel: '模型',
    chartSplitCredential: '凭据',
    chartSplitStatus: '状态',
    chartEmpty: '当前筛选下没有数据。',
    breakdownTitle: '分拆',
    breakdownEmpty: '没有可显示的分组。',
    heatmapTitle: '用量热力图',
    heatmapHint: '最近 53 周的每日用量',
    heatmapLess: '少',
    heatmapMore: '多',
    hmTipRequests: '{n} 次',
    hmTipTokens: '{n} token',
    heatmapTokens: 'Token',
    heatmapRequests: '请求',
    heatmapCost: '费用',
    kpiTokens: 'Token',
    kpiTokensSub: '输入 {input} · 输出 {output} · 缓存 {cached}',
    kpiRequests: '请求',
    kpiRequestsSub: '{failed} 次失败 · 成功 {rate}',
    kpiRequestsSub2: '{retries} 次重试 · {failovers} 次切换',
    kpiCost: '费用',
    kpiCostSub: '按公开价估算',
    kpiCache: '缓存命中率',
    kpiCacheSub: '{cached} / {input} 输入 token',
    kpiSpeed: '输出速度',
    kpiSpeedSub: '首字平均 {ttfb} · {n} 次计时',
    kpiNoData: '暂无数据',
    ready: '可调度',
    notReady: '还不能调度',
    activeCredential: '当前凭据 {label} · 上游失败按优先级切换',
    failoverNote: '上游失败按优先级切换凭据',
    setupNote: '还缺凭据或模型映射。先在凭据和模型页补齐。',
    requestPulse: '请求脉冲',
    pulseSummary: '最近 {n} 次 · 高度为耗时',
    pulseEmpty: '还没有真实请求。客户端接上后，这里会显示最近 30 次的耗时分布。',
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
    ignoreProxy: '忽略全局代理',
    ignoreProxyHint: '勾选后本插件直接连接 Mistral，不走 CPA 的全局代理。默认关闭。',
    ignoreProxyOn: '直连 · 绕过全局代理',
    ignoreProxyOff: '跟随 CPA 全局代理',
    on: '开启',
    saveSettings: '保存配置',
    failoverOffHint: 'failover 关闭时只使用当前排在最前的一把凭据。',
    settingsSaved: '配置已保存',
    logsTitle: '请求日志',
    logsHint: '日志在插件数据目录的 requests.jsonl。测试请求标为 probe，不计入统计。',
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
    sumCost: '费用',
    sumErrors: '失败',
    sumRetries: '重试',
    sumFailovers: '切换',
    sumProbes: '测试',
    colCredential: '凭据',
    colAttempts: '尝试',
    colTtfb: '首字',
    colInput: '输入',
    colOutput: '输出',
    colCached: '缓存',
    colCacheRate: '缓存率',
    colError: '错误',
    colMode: '模式',
    colSpeed: '速度',
    colCost: '费用',
    noLogs: '还没有请求记录。',
    probeTag: '测试',
    detailRequest: '请求与路由',
    detailModel: '模型与参数',
    detailTokens: 'Token 与缓存',
    detailPerf: '性能与延迟',
    detailError: '错误诊断',
    detailRequestId: '请求 ID',
    detailTime: '时间',
    detailFormat: '接口格式',
    detailStream: '响应模式',
    detailAttempts: '尝试次数',
    detailFailover: '故障切换',
    detailInputBytes: '请求大小',
    detailOutputBytes: '响应大小',
    detailClientModel: '客户端模型',
    detailUpstreamModel: '上游模型',
    detailReasoning: '思考档位',
    detailCredential: '凭据',
    detailPromptTokens: '输入 Token',
    detailCachedTokens: '缓存 Token',
    detailUncachedTokens: '未缓存输入',
    detailCompletionTokens: '输出 Token',
    detailTotalTokens: '总计 Token',
    detailCacheRate: '缓存命中率',
    detailCacheSource: '缓存键来源',
    detailCacheKey: '缓存键',
    detailPrefixHash: '前缀哈希',
    detailCost: '预估费用',
    detailTtfb: '首字延迟',
    detailGeneration: '生成耗时',
    detailTotalTime: '总耗时',
    detailSpeed: '生成速度',
    detailStreamYes: '流式',
    detailStreamNo: '非流式',
    yes: '是',
    no: '否',
    copyJson: '复制 JSON',
    copied: '已复制',
    filterByModel: '筛选此模型',
    filterByCredential: '筛选此凭据',
    errorHint: '检查凭据是否有效、模型映射是否正确，以及上游是否临时不可用。',
    cacheKeyClient: '缓存键来自客户端 prompt_cache_key',
    cacheKeyAffinity: '缓存键来自客户端 x-affinity 头',
    cacheKeyRoot: '缓存键按会话内容根（system 指令 + 首条用户消息）生成',
    cacheKeyInstructions: '缓存键按 system 指令生成',
    cacheKeySession: '缓存键来自 CPA 会话身份',
    cacheKeyUser: '缓存键按首条用户消息生成',
    cacheKeyNone: '没有缓存键',
    missingSession: '没有读到管理页的登录状态。请从 CPA 管理页的 Mistral 菜单打开。'
  },
  en: {
    tagline: 'CPA upstream plugin · credentials, models, and request ledger',
    navOverview: 'Overview',
    navCredentials: 'Credentials',
    navModels: 'Models',
    navLogs: 'Logs',
    navConfig: 'Settings',
    timeToday: 'Today',
    time24h: '24h',
    time7d: '7d',
    time30d: '30d',
    timeAll: 'All',
    autoOff: 'Manual',
    liveCost: 'estimated cost',
    exportCsv: 'Export CSV',
    filterAll: 'All',
    filterSuccess: 'Success',
    filterError: 'Error',
    filterRetry: 'Retried',
    filterProbe: 'Probe',
    searchPlaceholder: 'Search model, credential, error, request ID',
    statsNote: 'Probe requests are excluded from totals; cost is estimated from Mistral list prices.',
    filterModel: 'All models',
    filterCredential: 'All credentials',
    chartTitle: 'Usage trend',
    chartHint: 'Aggregated over the current range and filters',
    chartTokens: 'Tokens',
    chartCost: 'Cost',
    chartRequests: 'Requests',
    chartSpeed: 'Speed',
    chartSplitModel: 'Model',
    chartSplitCredential: 'Credential',
    chartSplitStatus: 'Status',
    chartEmpty: 'No data under the current filters.',
    breakdownTitle: 'Breakdown',
    breakdownEmpty: 'No groups to show.',
    heatmapTitle: 'Usage heatmap',
    heatmapHint: 'Daily usage over the last 53 weeks',
    heatmapLess: 'Less',
    heatmapMore: 'More',
    hmTipRequests: '{n} requests',
    hmTipTokens: '{n} tokens',
    heatmapTokens: 'Tokens',
    heatmapRequests: 'Requests',
    heatmapCost: 'Cost',
    kpiTokens: 'Tokens',
    kpiTokensSub: 'Input {input} · Output {output} · Cached {cached}',
    kpiRequests: 'Requests',
    kpiRequestsSub: '{failed} failed · {rate} success',
    kpiRequestsSub2: '{retries} retries · {failovers} failovers',
    kpiCost: 'Cost',
    kpiCostSub: 'Estimated from list prices',
    kpiCache: 'Cache hit rate',
    kpiCacheSub: '{cached} / {input} input tokens',
    kpiSpeed: 'Output speed',
    kpiSpeedSub: 'Avg first byte {ttfb} · {n} timed',
    kpiNoData: 'No data yet',
    ready: 'Ready',
    notReady: 'Not ready',
    activeCredential: 'Active: {label} · failover follows priority',
    failoverNote: 'Failover follows credential priority',
    setupNote: 'A credential or model mapping is missing. Add them under Credentials and Models.',
    requestPulse: 'Request pulse',
    pulseSummary: 'Last {n} · bar height is duration',
    pulseEmpty: 'No real requests yet. Once a client connects, the last 30 durations appear here.',
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
    ignoreProxy: 'Ignore global proxy',
    ignoreProxyHint: 'Connect to Mistral directly instead of through the CPA global proxy. Off by default.',
    ignoreProxyOn: 'Direct · bypasses the global proxy',
    ignoreProxyOff: 'Follows the CPA global proxy',
    on: 'On',
    saveSettings: 'Save settings',
    failoverOffHint: 'When failover is off, only the first credential in the current order is used.',
    settingsSaved: 'Settings saved',
    logsTitle: 'Request logs',
    logsHint: 'Logs live in requests.jsonl under the plugin data directory. Test requests are marked probe and excluded from statistics.',
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
    sumCost: 'Cost',
    sumErrors: 'Failed',
    sumRetries: 'Retries',
    sumFailovers: 'Failovers',
    sumProbes: 'Probes',
    colCredential: 'Credential',
    colAttempts: 'Attempts',
    colTtfb: 'First byte',
    colInput: 'Input',
    colOutput: 'Output',
    colCached: 'Cached',
    colCacheRate: 'Cache rate',
    colError: 'Error',
    colMode: 'Mode',
    colSpeed: 'Speed',
    colCost: 'Cost',
    noLogs: 'No requests logged yet.',
    probeTag: 'Test',
    detailRequest: 'Request & routing',
    detailModel: 'Model & parameters',
    detailTokens: 'Tokens & cache',
    detailPerf: 'Performance & latency',
    detailError: 'Error diagnosis',
    detailRequestId: 'Request ID',
    detailTime: 'Time',
    detailFormat: 'Format',
    detailStream: 'Response mode',
    detailAttempts: 'Attempts',
    detailFailover: 'Failover',
    detailInputBytes: 'Request size',
    detailOutputBytes: 'Response size',
    detailClientModel: 'Client model',
    detailUpstreamModel: 'Upstream model',
    detailReasoning: 'Reasoning',
    detailCredential: 'Credential',
    detailPromptTokens: 'Input tokens',
    detailCachedTokens: 'Cached tokens',
    detailUncachedTokens: 'Uncached input',
    detailCompletionTokens: 'Output tokens',
    detailTotalTokens: 'Total tokens',
    detailCacheRate: 'Cache hit rate',
    detailCacheSource: 'Cache key source',
    detailCacheKey: 'Cache key',
    detailPrefixHash: 'Prefix hash',
    detailCost: 'Estimated cost',
    detailTtfb: 'First byte',
    detailGeneration: 'Generation time',
    detailTotalTime: 'Total time',
    detailSpeed: 'Generation speed',
    detailStreamYes: 'Stream',
    detailStreamNo: 'Non-stream',
    yes: 'Yes',
    no: 'No',
    copyJson: 'Copy JSON',
    copied: 'Copied',
    filterByModel: 'Filter by this model',
    filterByCredential: 'Filter by this credential',
    errorHint: 'Check the credential, the model mapping, and whether the upstream is temporarily unavailable.',
    cacheKeyClient: 'Cache key from client prompt_cache_key',
    cacheKeyAffinity: 'Cache key from the client x-affinity header',
    cacheKeyRoot: 'Cache key from the conversation root (instructions + first user message)',
    cacheKeyInstructions: 'Cache key derived from system instructions',
    cacheKeySession: 'Cache key from CPA session identity',
    cacheKeyUser: 'Cache key derived from first user message',
    cacheKeyNone: 'No cache key',
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
let timeFilter = 'today';
let statusFilter = 'all';
let modelFilter = '';
let credFilter = '';
let searchQuery = '';
let chartMetric = 'cost';
let chartSplit = 'model';
let heatmapMetric = 'cost';
let autoRefreshSeconds = 0;
let autoRefreshTimer = null;
let chartStep = 3600000;
const PALETTE = ['#6366f1','#8b5cf6','#ec4899','#f59e0b','#10b981','#06b6d4','#f43f5e','#84cc16'];

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
  document.querySelectorAll('[data-i18n-placeholder]').forEach(node => { node.placeholder = t(node.dataset.i18nPlaceholder); });
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
}
function switchTab(name) {
  document.querySelectorAll('.tab').forEach(item => item.classList.toggle('active', item.dataset.tab === name));
  document.querySelectorAll('main section').forEach(section => section.hidden = section.id !== name);
  const filterBar = document.querySelector('#filter-bar');
  if (filterBar) filterBar.hidden = name === 'credentials' || name === 'models' || name === 'config';
  msg.textContent = '';
  if (name === 'overview') {
    renderChart(filteredLogs());
    renderHeatmap(heatmapLogs());
  }
  if (name === 'logs') renderLogs();
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
  if (!value) return '—';
  return value >= 1000 ? (value / 1000).toFixed(1) + 's' : Math.round(value) + 'ms';
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
function formatBytes(value) {
  const n = Number(value || 0);
  if (!n) return '—';
  if (n >= 1073741824) return (n / 1073741824).toFixed(1) + ' GB';
  if (n >= 1048576) return (n / 1048576).toFixed(1) + ' MB';
  if (n >= 1024) return (n / 1024).toFixed(1) + ' KB';
  return n + ' B';
}
function timeText(value) {
  if (!value) return '';
  return String(value).slice(5, 19).replace('T', ' ');
}
function shortTime(value) {
  if (!value) return '';
  return String(value).slice(5, 16).replace('T', ' ');
}
function fullTime(value) {
  if (!value) return '—';
  return String(value).slice(0, 19).replace('T', ' ');
}
function metric(label, value, unit) {
  return '<div class="metric"><span>' + label + '</span><b>' + value + (unit ? '<small>' + unit + '</small>' : '') + '</b></div>';
}
function fact(label, value) {
  return '<div class="fact"><span>' + label + '</span><b>' + value + '</b></div>';
}
function estimateCostClient(upstream, prompt, cached, completion) {
  const up = String(upstream || '').toLowerCase();
  let promptRate = 1.0, cachedRate = 0.5, completionRate = 3.0;
  if (up.indexOf('large') !== -1) { promptRate = 2.0; cachedRate = 1.0; completionRate = 6.0; }
  else if (up.indexOf('codestral') !== -1) { promptRate = 0.3; cachedRate = 0.15; completionRate = 0.9; }
  else if (up.indexOf('small') !== -1) { promptRate = 0.1; cachedRate = 0.05; completionRate = 0.3; }
  else if (up.indexOf('ministral-8b') !== -1) { promptRate = 0.1; cachedRate = 0.1; completionRate = 0.1; }
  else if (up.indexOf('ministral-3b') !== -1) { promptRate = 0.04; cachedRate = 0.04; completionRate = 0.04; }
  else if (up.indexOf('pixtral-12b') !== -1) { promptRate = 0.15; cachedRate = 0.15; completionRate = 0.15; }
  else if (up.indexOf('pixtral') !== -1) { promptRate = 2.0; cachedRate = 1.0; completionRate = 6.0; }
  else if (up.indexOf('embed') !== -1 || up.indexOf('moderation') !== -1) { promptRate = 0.1; cachedRate = 0.1; completionRate = 0; }
  const uncached = Math.max(0, Number(prompt || 0) - Number(cached || 0));
  return (uncached * promptRate + Number(cached || 0) * cachedRate + Number(completion || 0) * completionRate) / 1000000;
}
function itemCost(item) {
  if (item && item.cost) return Number(item.cost);
  return estimateCostClient(item && item.upstream, item && item.prompt_tokens, item && item.cached_tokens, item && item.completion_tokens);
}
function money(value) {
  const n = Number(value || 0);
  if (n === 0) return '$0.00';
  if (n < 0.01) return '$' + n.toFixed(4);
  return '$' + n.toFixed(2);
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
function allLogs() {
  return Array.isArray(data && data.recent) ? data.recent : [];
}
function timeCutoff() {
  const now = new Date();
  if (timeFilter === 'today') {
    const d = new Date(now); d.setHours(0, 0, 0, 0); return d.getTime();
  }
  if (timeFilter === '24h') {
    const d = new Date(now.getTime() - 23 * 3600000); d.setMinutes(0, 0, 0); return d.getTime();
  }
  if (timeFilter === '7d') {
    const d = new Date(now); d.setDate(now.getDate() - 6); d.setHours(0, 0, 0, 0); return d.getTime();
  }
  if (timeFilter === '30d') {
    const d = new Date(now); d.setDate(now.getDate() - 29); d.setHours(0, 0, 0, 0); return d.getTime();
  }
  return 0;
}
function filterBase(cutoff) {
  const needle = searchQuery.trim().toLowerCase();
  return allLogs().filter(item => {
    if (cutoff) {
      const ts = new Date(item.time).getTime();
      if (!ts || ts < cutoff) return false;
    }
    if (modelFilter && item.model !== modelFilter) return false;
    if (credFilter && item.credential !== credFilter) return false;
    if (statusFilter === 'probe') { if (!item.probe) return false; }
    else if (statusFilter === 'success') { if (item.probe || item.status < 200 || item.status >= 400) return false; }
    else if (statusFilter === 'error') { if (item.probe || (item.status >= 200 && item.status < 400)) return false; }
    else if (statusFilter === 'retry') { if (item.probe || !(item.attempts > 1 || item.failover)) return false; }
    if (needle) {
      const hay = [item.model, item.upstream, item.credential, item.error, item.request_id, item.format].join(' ').toLowerCase();
      if (hay.indexOf(needle) === -1) return false;
    }
    return true;
  });
}
function filteredLogs() {
  return filterBase(timeCutoff());
}
function heatmapLogs() {
  return filterBase(0);
}
function computeStats(logs) {
  let requests = 0, success = 0, errors = 0, streams = 0, probes = 0, retries = 0, failovers = 0;
  let prompt = 0, completion = 0, cached = 0, cost = 0, ttfbTotal = 0, durationTotal = 0, generationMs = 0, inputBytes = 0, outputBytes = 0;
  let speedSum = 0, speedCount = 0, ttfbCount = 0;
  const byModel = {};
  logs.forEach(item => {
    if (item.probe) { probes++; return; }
    requests++;
    if (item.stream) streams++;
    if (item.attempts > 1) retries += item.attempts - 1;
    if (item.failover) failovers++;
    if (item.status >= 200 && item.status < 400) success++; else errors++;
    if (item.ttfb_ms > 0) { ttfbTotal += item.ttfb_ms; ttfbCount++; }
    durationTotal += item.duration_ms || 0;
    prompt += item.prompt_tokens || 0;
    completion += item.completion_tokens || 0;
    cached += item.cached_tokens || 0;
    inputBytes += item.input_bytes || 0;
    outputBytes += item.output_bytes || 0;
    cost += itemCost(item);
    if (item.tokens_per_second > 0) { speedSum += item.tokens_per_second; speedCount++; }
    const gen = Math.max(0, (item.duration_ms || 0) - (item.ttfb_ms || 0));
    if (item.completion_tokens > 0 && gen > 0) generationMs += gen;
    const key = item.model || item.upstream || '';
    if (key) {
      if (!byModel[key]) byModel[key] = { model: key, requests: 0, errors: 0, prompt_tokens: 0, completion_tokens: 0, cached_tokens: 0, cost: 0, duration_total: 0, ttfb_total: 0, speed_sum: 0, speed_count: 0 };
      const bm = byModel[key];
      bm.requests++;
      if (item.status < 200 || item.status >= 400) bm.errors++;
      bm.prompt_tokens += item.prompt_tokens || 0;
      bm.completion_tokens += item.completion_tokens || 0;
      bm.cached_tokens += item.cached_tokens || 0;
      bm.cost += itemCost(item);
      bm.duration_total += item.duration_ms || 0;
      bm.ttfb_total += item.ttfb_ms || 0;
      if (item.tokens_per_second > 0) { bm.speed_sum += item.tokens_per_second; bm.speed_count++; }
    }
  });
  const out = {
    requests, success, errors, streams, probes, retries, failovers,
    prompt_tokens: prompt, completion_tokens: completion, cached_tokens: cached,
    total_tokens: prompt + completion, total_cost: cost,
    input_bytes: inputBytes, output_bytes: outputBytes,
    models_breakdown: byModel
  };
  if (requests > 0) {
    out.success_rate = success / requests;
    out.avg_ttfb_ms = ttfbCount ? ttfbTotal / ttfbCount : 0;
    out.avg_duration_ms = durationTotal / requests;
  }
  if (generationMs > 0) out.tokens_per_second = completion * 1000 / generationMs;
  if (prompt > 0) out.cache_hit_rate = cached / prompt;
  if (speedCount > 0) out.avg_speed = speedSum / speedCount;
  return out;
}
function render() {
  if (!data) return;
  document.querySelector('#state').textContent = 'v' + data.version;
  const routeState = document.querySelector('#route-state');
  const routeText = document.querySelector('#route-state-text');
  routeState.className = 'route-state ' + (data.configured ? 'ok' : 'warn');
  routeText.textContent = data.configured ? t('ready') : t('notReady');
  renderFilterOptions();
  renderOverview();
  renderCredentials();
  renderModels();
  renderConfig();
  if (!document.querySelector('#logs').hidden) renderLogs();
}
function renderFilterOptions() {
  const logs = allLogs();
  const models = Array.from(new Set(logs.map(item => item.model).filter(Boolean))).sort();
  const creds = Array.from(new Set(logs.map(item => item.credential).filter(Boolean))).sort();
  const modelSel = document.querySelector('#model-filter');
  const credSel = document.querySelector('#cred-filter');
  if (modelSel) {
    const current = modelSel.value;
    modelSel.innerHTML = '<option value="">' + t('filterModel') + '</option>' + models.map(m => '<option value="' + esc(m) + '">' + esc(m) + '</option>').join('');
    modelSel.value = models.indexOf(current) !== -1 ? current : '';
  }
  if (credSel) {
    const current = credSel.value;
    credSel.innerHTML = '<option value="">' + t('filterCredential') + '</option>' + creds.map(c => '<option value="' + esc(c) + '">' + esc(c) + '</option>').join('');
    credSel.value = creds.indexOf(current) !== -1 ? current : '';
  }
}
function renderOverview() {
  const logs = filteredLogs();
  const stats = computeStats(logs);
  document.querySelector('#cost-chip').innerHTML = '≈' + money(stats.total_cost) + ' <span>' + t('liveCost') + '</span>';
  renderKpis(stats);
  renderChart(logs);
  renderBreakdown(logs);
  renderHeatmap(heatmapLogs());
  renderRecent(logs);
}
function renderKpis(stats) {
  const has = stats.requests > 0;
  const cacheRate = has && stats.cache_hit_rate !== undefined ? Math.max(0, Math.min(1, stats.cache_hit_rate)) : 0;
  const speed = stats.avg_speed || 0;
  const timed = stats.requests;
  const kpis = [
    '<div class="kpi"><div class="kpi-label">' + t('kpiTokens') + icon('database') + '</div>'
      + '<div class="kpi-value">' + (has ? compact(stats.total_tokens) : '—') + '</div>'
      + '<div class="kpi-sub">' + t('kpiTokensSub', {input:compact(stats.prompt_tokens), output:compact(stats.completion_tokens), cached:compact(stats.cached_tokens)}) + '</div></div>',
    '<div class="kpi"><div class="kpi-label">' + t('kpiRequests') + icon('activity') + '</div>'
      + '<div class="kpi-value">' + (has ? num(stats.requests) : '—') + '</div>'
      + '<div class="kpi-sub"><span class="' + (stats.errors ? 'bad' : '') + '">' + t('kpiRequestsSub', {failed:num(stats.errors), rate:has ? pct(stats.success_rate) : '—'}) + '</span><br>' + t('kpiRequestsSub2', {retries:num(stats.retries), failovers:num(stats.failovers)}) + '</div></div>',
    '<div class="kpi"><div class="kpi-label">' + t('kpiCost') + icon('coin') + '</div>'
      + '<div class="kpi-value good">' + (has ? money(stats.total_cost) : '—') + '</div>'
      + '<div class="kpi-sub">' + t('kpiCostSub') + '</div></div>',
    '<div class="kpi"><div class="kpi-label">' + t('kpiCache') + icon('zap') + '</div>'
      + '<div class="kpi-value">' + (has && stats.cache_hit_rate !== undefined ? pct1(stats.cache_hit_rate) : '—') + '</div>'
      + '<div class="kpi-bar"><i style="width:' + (cacheRate * 100).toFixed(1) + '%"></i></div>'
      + '<div class="kpi-sub">' + t('kpiCacheSub', {cached:compact(stats.cached_tokens), input:compact(stats.prompt_tokens)}) + '</div></div>',
    '<div class="kpi"><div class="kpi-label">' + t('kpiSpeed') + icon('gauge') + '</div>'
      + '<div class="kpi-value">' + (has && speed ? speed.toFixed(0) : '—') + (has && speed ? '<small>tok/s</small>' : '') + '</div>'
      + '<div class="kpi-sub">' + t('kpiSpeedSub', {ttfb:sec(stats.avg_ttfb_ms), n:num(timed)}) + '</div></div>'
  ];
  document.querySelector('#kpi-row').innerHTML = kpis.join('');
}
function icon(name) {
  const paths = {
    database: '<ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14a9 3 0 0 0 18 0V5"/><path d="M3 12a9 3 0 0 0 18 0"/>',
    activity: '<path d="M22 12h-4l-3 9L9 3l-3 9H2"/>',
    coin: '<circle cx="12" cy="12" r="9"/><path d="M14.8 9.2c-.5-.8-1.5-1.2-2.8-1.2-1.6 0-2.8.8-2.8 2s1.2 1.8 2.8 2.1 2.8.9 2.8 2.1-1.2 2-2.8 2c-1.3 0-2.3-.4-2.8-1.2"/><path d="M12 6v12"/>',
    zap: '<path d="M13 2 3 14h8l-1 8 10-12h-8l1-8z"/>',
    gauge: '<path d="M12 15l3.5-3.5"/><path d="M20.3 18a9 9 0 1 0-16.6 0"/>'
  };
  return '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">' + (paths[name] || '') + '</svg>';
}
function splitKey(item, mode) {
  if (mode === 'credential') return item.credential || '—';
  if (mode === 'status') return item.status >= 200 && item.status < 400 ? 'success' : 'error';
  return item.model || item.upstream || '—';
}
function buildBuckets(logs, metric) {
  const now = new Date();
  let start, step, count;
  if (timeFilter === 'today') {
    start = new Date(now); start.setHours(0, 0, 0, 0);
    step = 3600000;
    count = now.getHours() + 1;
  } else if (timeFilter === '24h') {
    start = new Date(now.getTime() - 23 * 3600000);
    start.setMinutes(0, 0, 0);
    step = 3600000;
    count = 24;
  } else if (timeFilter === '7d') {
    start = new Date(now);
    start.setDate(now.getDate() - 6);
    start.setHours(0, 0, 0, 0);
    step = 86400000;
    count = 7;
  } else if (timeFilter === '30d') {
    start = new Date(now);
    start.setDate(now.getDate() - 29);
    start.setHours(0, 0, 0, 0);
    step = 86400000;
    count = 30;
  } else {
    const times = logs.map(item => new Date(item.time).getTime()).filter(ts => ts > 0);
    const earliest = times.length ? Math.min.apply(null, times) : now.getTime() - 29 * 86400000;
    start = new Date(earliest);
    start.setHours(0, 0, 0, 0);
    step = 86400000;
    const days = Math.floor((now.getTime() - start.getTime()) / 86400000);
    count = Math.max(1, days + 1);
  }
  chartStep = step;
  const buckets = [];
  for (let i = 0; i < count; i++) {
    buckets.push({ time: start.getTime() + i * step, value: 0, speedSum: 0, speedCount: 0, groups: {} });
  }
  logs.forEach(item => {
    if (item.probe) return;
    const ts = new Date(item.time).getTime();
    if (!ts) return;
    let idx = Math.floor((ts - start.getTime()) / step);
    if (idx >= 0 && idx < buckets.length) {
      const b = buckets[idx];
      const v = metricValue(item, metric);
      if (metric === 'speed') {
        if (v > 0) { b.speedSum += v; b.speedCount++; }
      } else {
        b.value += v;
      }
      const key = splitKey(item, chartSplit);
      b.groups[key] = (b.groups[key] || 0) + v;
    }
  });
  if (metric === 'speed') {
    buckets.forEach(b => { b.value = b.speedCount ? b.speedSum / b.speedCount : 0; });
  }
  return buckets;
}
function metricValue(item, metric) {
  if (metric === 'tokens') return (item.prompt_tokens || 0) + (item.completion_tokens || 0);
  if (metric === 'cost') return itemCost(item);
  if (metric === 'requests') return 1;
  if (metric === 'speed') return item.tokens_per_second || 0;
  return 0;
}
function niceCeil(value) {
  if (!(value > 0)) return 1;
  const exp = Math.floor(Math.log10(value));
  const base = Math.pow(10, exp);
  const n = value / base;
  let nice;
  if (n <= 1) nice = 1;
  else if (n <= 2) nice = 2;
  else if (n <= 5) nice = 5;
  else nice = 10;
  return nice * base;
}
function formatMetric(value, metric) {
  if (metric === 'cost') return money(value);
  if (metric === 'requests') return String(Math.round(value));
  if (metric === 'speed') return String(Math.round(value));
  return compact(value);
}
function formatBucketTime(ts) {
  const d = new Date(ts);
  if (chartStep < 86400000) return String(d.getHours()).padStart(2, '0') + ':00';
  return (d.getMonth() + 1) + '-' + String(d.getDate()).padStart(2, '0');
}
function smoothPath(points, top, bottom) {
  if (!points.length) return '';
  if (points.length === 1) return 'M' + points[0].x + ',' + points[0].y;
  let d = 'M' + points[0].x.toFixed(1) + ',' + points[0].y.toFixed(1);
  for (let i = 0; i < points.length - 1; i++) {
    const p0 = points[i - 1] || points[i];
    const p1 = points[i];
    const p2 = points[i + 1];
    const p3 = points[i + 2] || p2;
    let cp1x = p1.x + (p2.x - p0.x) / 6;
    let cp1y = p1.y + (p2.y - p0.y) / 6;
    let cp2x = p2.x - (p3.x - p1.x) / 6;
    let cp2y = p2.y - (p3.y - p1.y) / 6;
    if (top !== undefined) {
      cp1y = Math.max(cp1y, top);
      cp2y = Math.max(cp2y, top);
    }
    if (bottom !== undefined) {
      cp1y = Math.min(cp1y, bottom);
      cp2y = Math.min(cp2y, bottom);
    }
    d += ' C' + cp1x.toFixed(1) + ',' + cp1y.toFixed(1) + ' ' + cp2x.toFixed(1) + ',' + cp2y.toFixed(1) + ' ' + p2.x.toFixed(1) + ',' + p2.y.toFixed(1);
  }
  return d;
}
function renderChart(logs) {
  const wrap = document.querySelector('#chart-wrap');
  if (!wrap) return;
  const buckets = buildBuckets(logs, chartMetric);
  const width = Math.max(320, wrap.clientWidth || 640);
  const height = 250;
  const padL = 52, padR = 14, padT = 14, padB = 28;
  const innerW = width - padL - padR;
  const innerH = height - padT - padB;
  const maxVal = Math.max.apply(null, buckets.map(b => b.value).concat([0]));
  if (!maxVal) {
    wrap.innerHTML = '<div class="chart-empty">' + t('chartEmpty') + '</div>';
    return;
  }
  const niceMax = niceCeil(maxVal);
  const points = buckets.map((b, i) => ({
    x: padL + (buckets.length === 1 ? innerW / 2 : (i / (buckets.length - 1)) * innerW),
    y: padT + innerH - (b.value / niceMax) * innerH,
    value: b.value,
    time: b.time,
    groups: b.groups
  }));
  const bottom = padT + innerH;
  let grid = '';
  for (let i = 0; i <= 4; i++) {
    const y = padT + (innerH / 4) * i;
    const val = niceMax * (1 - i / 4);
    grid += '<line x1="' + padL + '" y1="' + y.toFixed(1) + '" x2="' + (width - padR) + '" y2="' + y.toFixed(1) + '" class="chart-grid-line"/>';
    grid += '<text x="' + (padL - 8) + '" y="' + (y + 4).toFixed(1) + '" class="chart-axis-label" text-anchor="end">' + formatMetric(val, chartMetric) + '</text>';
  }
  let xLabels = '';
  const labelCount = Math.min(6, buckets.length);
  for (let i = 0; i < labelCount; i++) {
    const idx = Math.round((i / Math.max(1, labelCount - 1)) * (buckets.length - 1));
    xLabels += '<text x="' + points[idx].x.toFixed(1) + '" y="' + (height - 8) + '" class="chart-axis-label" text-anchor="middle">' + formatBucketTime(points[idx].time) + '</text>';
  }
  const linePath = smoothPath(points, padT, bottom);
  const areaPath = points.length > 1
    ? linePath + ' L' + points[points.length - 1].x.toFixed(1) + ',' + bottom + ' L' + points[0].x.toFixed(1) + ',' + bottom + ' Z'
    : '';
  let dots = '';
  if (points.length === 1) {
    dots = '<circle cx="' + points[0].x.toFixed(1) + '" cy="' + points[0].y.toFixed(1) + '" r="4" class="chart-dot"/>';
  }
  const svg = '<svg viewBox="0 0 ' + width + ' ' + height + '" preserveAspectRatio="none" role="img" aria-label="' + t('chartTitle') + '">'
    + '<defs><linearGradient id="chartGradient" x1="0" y1="0" x2="0" y2="1">'
    + '<stop offset="0%" stop-color="#6366f1" stop-opacity=".22"/>'
    + '<stop offset="100%" stop-color="#6366f1" stop-opacity="0"/>'
    + '</linearGradient></defs>'
    + grid + xLabels
    + (areaPath ? '<path d="' + areaPath + '" class="chart-area"/>' : '')
    + '<path d="' + linePath + '" class="chart-line"/>'
    + dots
    + '<line id="chart-cross" class="chart-cross" x1="0" y1="' + padT + '" x2="0" y2="' + bottom + '" style="display:none"/>'
    + '<circle id="chart-point" r="4" class="chart-dot" style="display:none"/>'
    + '<rect id="chart-hit" x="' + padL + '" y="' + padT + '" width="' + innerW + '" height="' + innerH + '" fill="transparent"/>'
    + '</svg>'
    + '<div id="chart-tip" class="chart-tip" hidden></div>';
  wrap.innerHTML = svg;
  const hit = document.querySelector('#chart-hit');
  const cross = document.querySelector('#chart-cross');
  const point = document.querySelector('#chart-point');
  const tip = document.querySelector('#chart-tip');
  if (!hit) return;
  hit.addEventListener('mousemove', event => {
    const rect = wrap.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const ratio = Math.max(0, Math.min(1, (x - padL) / innerW));
    const idx = Math.round(ratio * (buckets.length - 1));
    const p = points[idx];
    cross.setAttribute('x1', p.x.toFixed(1));
    cross.setAttribute('x2', p.x.toFixed(1));
    cross.style.display = '';
    point.setAttribute('cx', p.x.toFixed(1));
    point.setAttribute('cy', p.y.toFixed(1));
    point.style.display = '';
    const groups = Object.keys(p.groups || {}).sort((a, b) => p.groups[b] - p.groups[a]).slice(0, 4);
    tip.innerHTML = '<div class="tip-time">' + formatBucketTime(p.time) + '</div>'
      + '<div class="tip-row"><span>' + t('chartTitle') + '</span><b>' + formatMetric(p.value, chartMetric) + '</b></div>'
      + groups.map(g => '<div class="tip-row"><span>' + esc(g) + '</span><b>' + formatMetric(p.groups[g], chartMetric) + '</b></div>').join('');
    tip.hidden = false;
    const tipX = Math.max(70, Math.min(width - 70, p.x));
    tip.style.left = tipX + 'px';
    tip.style.top = Math.max(0, p.y - 10) + 'px';
  });
  hit.addEventListener('mouseleave', () => {
    cross.style.display = 'none';
    point.style.display = 'none';
    tip.hidden = true;
  });
}
function renderBreakdown(logs) {
  const groups = {};
  logs.forEach(item => {
    if (item.probe) return;
    const key = splitKey(item, chartSplit);
    if (!groups[key]) groups[key] = { name: key, requests: 0, tokens: 0, cost: 0, errors: 0, speedSum: 0, speedCount: 0, cached: 0, prompt: 0 };
    const g = groups[key];
    g.requests++;
    g.tokens += (item.prompt_tokens || 0) + (item.completion_tokens || 0);
    g.cost += itemCost(item);
    if (item.status < 200 || item.status >= 400) g.errors++;
    if (item.tokens_per_second > 0) { g.speedSum += item.tokens_per_second; g.speedCount++; }
    g.cached += item.cached_tokens || 0;
    g.prompt += item.prompt_tokens || 0;
  });
  const list = Object.values(groups).sort((a, b) => b.cost - a.cost);
  const sub = document.querySelector('#breakdown-sub');
  if (sub) sub.textContent = t('chartSplit' + chartSplit.charAt(0).toUpperCase() + chartSplit.slice(1)) + ' · ' + list.length;
  const target = document.querySelector('#breakdown-list');
  if (!list.length) {
    target.innerHTML = '<div class="empty">' + t('breakdownEmpty') + '</div>';
    return;
  }
  target.innerHTML = list.map((g, i) => {
    const rate = g.prompt > 0 ? g.cached / g.prompt : 0;
    const speed = g.speedCount ? g.speedSum / g.speedCount : 0;
    const active = chartSplit === 'model' && modelFilter === g.name;
    return '<div class="breakdown-item' + (active ? ' active' : '') + '" data-bd="' + esc(g.name) + '">'
      + '<span class="breakdown-dot" style="background:' + PALETTE[i % PALETTE.length] + '"></span>'
      + '<div class="breakdown-main"><div class="breakdown-name" title="' + esc(g.name) + '">' + esc(g.name) + '</div>'
      + '<div class="breakdown-meta">' + num(g.requests) + ' ' + t('kpiRequests').toLowerCase() + ' · ' + compact(g.tokens) + ' token'
      + (speed ? ' · ' + speed.toFixed(0) + ' tok/s' : '')
      + ' · ' + t('colCacheRate') + ' ' + pct(rate)
      + (g.errors ? ' · <span class="bad">' + g.errors + ' ' + t('sumErrors').toLowerCase() + '</span>' : '') + '</div></div>'
      + '<div class="breakdown-value">' + money(g.cost) + '<small>' + compact(g.tokens) + ' tok</small></div></div>';
  }).join('');
  target.querySelectorAll('[data-bd]').forEach(node => {
    node.onclick = () => {
      const name = node.dataset.bd;
      if (chartSplit === 'model') {
        modelFilter = modelFilter === name ? '' : name;
        const sel = document.querySelector('#model-filter');
        if (sel) sel.value = modelFilter;
      } else if (chartSplit === 'credential') {
        credFilter = credFilter === name ? '' : name;
        const sel = document.querySelector('#cred-filter');
        if (sel) sel.value = credFilter;
      }
      renderOverview();
    };
  });
}
function renderHeatmap(logs) {
  const target = document.querySelector('#heatmap');
  if (!target) return;
  const byDay = {};
  logs.forEach(item => {
    if (item.probe) return;
    const d = new Date(item.time);
    if (isNaN(d.getTime())) return;
    const key = d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate();
    if (!byDay[key]) byDay[key] = { cost: 0, tokens: 0, requests: 0 };
    byDay[key].cost += itemCost(item);
    byDay[key].tokens += (item.prompt_tokens || 0) + (item.completion_tokens || 0);
    byDay[key].requests += 1;
  });
  const end = new Date(); end.setHours(0, 0, 0, 0);
  const start = new Date(end); start.setDate(start.getDate() - 364);
  const startDay = (start.getDay() + 6) % 7;
  const gridStart = new Date(start); gridStart.setDate(gridStart.getDate() - startDay);
  const cells = [];
  for (let d = new Date(gridStart); d <= end; d.setDate(d.getDate() + 1)) {
    const key = d.getFullYear() + '-' + (d.getMonth() + 1) + '-' + d.getDate();
    const rec = byDay[key];
    const value = rec ? rec[heatmapMetric] : 0;
    cells.push({ date: new Date(d), value: value, rec: rec, key: key, inRange: d >= start });
  }
  // Month labels: strictly at the column containing the 1st of the month.
  const months = [];
  const numCols = Math.ceil(cells.length / 7);
  for (let col = 0; col < numCols; col++) {
    const colCells = cells.slice(col * 7, col * 7 + 7).filter(c => c.inRange);
    if (!colCells.length) continue;
    const first = colCells.find(c => c.date.getDate() === 1);
    if (first) months.push({ label: (first.date.getMonth() + 1) + (lang === 'zh' ? '月' : ''), col: col });
  }
  // Hide labels closer than 2 columns to their predecessor.
  const shownMonths = [];
  months.forEach(m => {
    if (shownMonths.length && m.col - shownMonths[shownMonths.length - 1].col < 2) return;
    shownMonths.push(m);
  });
  const values = cells.map(c => c.value).filter(v => v > 0).sort((a, b) => a - b);
  const q = p => values.length ? values[Math.min(values.length - 1, Math.floor(values.length * p))] : 0;
  const q1 = q(0.25), q2 = q(0.5), q3 = q(0.75);
  function level(v) {
    if (!(v > 0)) return 0;
    if (values.length === 1) return 2;
    if (v <= q1) return 1;
    if (v <= q2) return 2;
    if (v <= q3) return 3;
    return 4;
  }
  const zhWeekdays = ['周日','周一','周二','周三','周四','周五','周六'];
  const enWeekdays = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'];
  const cellHtml = cells.map(c => {
    if (!c.inRange) return '<span class="hm-cell blank"></span>';
    const lv = level(c.value);
    const rec = c.rec || { cost: 0, tokens: 0, requests: 0 };
    const wd = (lang === 'zh' ? zhWeekdays : enWeekdays)[c.date.getDay()];
    const dateLabel = c.date.getFullYear() + '-' + String(c.date.getMonth() + 1).padStart(2, '0') + '-' + String(c.date.getDate()).padStart(2, '0') + ' ' + wd;
    return '<span class="hm-cell l' + lv + '" data-hm-date="' + esc(dateLabel) + '" data-hm-cost="' + money(rec.cost) + '" data-hm-requests="' + num(rec.requests) + '" data-hm-tokens="' + num(rec.tokens) + '"></span>';
  }).join('');
  const weekdays = lang === 'zh' ? ['周一','','周三','','周五','','周日'] : ['Mon','','Wed','','Fri','','Sun'];
  heatmapMonths = shownMonths;
  heatmapFitWidth = 0;
  target.innerHTML = '<div class="hm-weekdays">' + weekdays.map(w => '<span>' + w + '</span>').join('') + '</div>'
    + '<div class="hm-main"><div class="hm-months"></div>'
    + '<div class="hm-grid">' + cellHtml + '</div></div>';
  bindHeatmapTip(target);
  fitHeatmap();
}
function bindHeatmapTip(target) {
  const tip = document.querySelector('#hm-tip');
  if (!tip) return;
  const card = target.closest('.card') || target.parentElement;
  const show = cell => {
    tip.innerHTML = '<div class="tip-time">' + esc(cell.dataset.hmDate) + '</div>'
      + '<div class="tip-row"><span>' + t('kpiCost') + '</span><b>' + esc(cell.dataset.hmCost) + '</b></div>'
      + '<div class="tip-row"><span>' + t('kpiRequests') + '</span><b>' + t('hmTipRequests', {n:esc(cell.dataset.hmRequests)}) + '</b></div>'
      + '<div class="tip-row"><span>' + t('kpiTokens') + '</span><b>' + t('hmTipTokens', {n:esc(cell.dataset.hmTokens)}) + '</b></div>';
    tip.hidden = false;
    const cardRect = card.getBoundingClientRect();
    const rect = cell.getBoundingClientRect();
    const left = rect.left - cardRect.left + rect.width / 2;
    const top = rect.top - cardRect.top;
    tip.style.left = Math.max(70, Math.min(cardRect.width - 70, left)) + 'px';
    tip.style.top = Math.max(0, top - 6) + 'px';
  };
  target.querySelectorAll('.hm-cell:not(.blank)').forEach(cell => {
    cell.addEventListener('mouseenter', () => show(cell));
    cell.addEventListener('mouseleave', () => { tip.hidden = true; });
  });
}
let heatmapMonths = [];
let heatmapFitWidth = 0;
function fitHeatmap() {
  const target = document.querySelector('#heatmap');
  if (!target) return;
  const grid = target.querySelector('.hm-grid');
  const monthsEl = target.querySelector('.hm-months');
  if (!grid || !monthsEl) return;
  const numCols = Math.max(1, Math.ceil(grid.children.length / 7));
  const gap = 3;
  const gridW = grid.clientWidth;
  if (!gridW || gridW === heatmapFitWidth) return;
  heatmapFitWidth = gridW;
  const cellSize = Math.max(12, Math.min(28, Math.floor((gridW - (numCols - 1) * gap) / numCols)));
  const step = numCols > 1 ? (gridW - cellSize) / (numCols - 1) : cellSize;
  target.style.setProperty('--hm-size', cellSize + 'px');
  target.style.setProperty('--hm-gap', gap + 'px');
  target.style.setProperty('--hm-radius', (cellSize >= 18 ? 5 : cellSize >= 15 ? 4 : 3) + 'px');
  monthsEl.innerHTML = heatmapMonths.map(m => '<span style="left:' + (m.col * step).toFixed(1) + 'px">' + esc(m.label) + '</span>').join('');
}
function renderRecent(logs) {
  const rows = logs.filter(item => !item.probe).slice(0, 8).map(item => {
    const ok = item.status >= 200 && item.status < 400;
    return '<tr><td class="mono nowrap">' + esc(shortTime(item.time)) + '</td>'
      + '<td class="model-cell" title="' + esc(item.model || '') + '">' + esc(item.model) + '</td>'
      + '<td class="nowrap"><span class="badge ' + (ok ? 'ok' : 'bad') + '">' + item.status + '</span></td>'
      + '<td class="num">' + sec(item.duration_ms) + '</td>'
      + '<td class="num">' + compactTokens(item.prompt_tokens) + ' → ' + compactTokens(item.completion_tokens) + '</td>'
      + '<td class="num">' + money(itemCost(item)) + '</td></tr>';
  }).join('');
  const html = rows
    ? '<table><thead><tr><th>' + t('colTime') + '</th><th>' + t('colModel') + '</th><th>' + t('colStatus') + '</th><th class="num">' + t('colDuration') + '</th><th class="num">' + t('colInOut') + '</th><th class="num">' + t('colCost') + '</th></tr></thead><tbody>' + rows + '</tbody></table>'
    : '<div class="empty"><b>' + t('noRequests') + '</b>' + t('noRequestsHint') + '</div>';
  document.querySelector('#recent-table').innerHTML = html;
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
    + '<div class="field"><label>' + t('ignoreProxy') + '<span class="hint">' + t('ignoreProxyHint') + '</span></label><span class="switch"><input id="ignore-proxy" type="checkbox"' + (data.ignore_proxy ? ' checked' : '') + '> ' + t('on') + '</span></div>'
    + '</div><div class="form"><button class="primary" id="save-settings">' + t('saveSettings') + '</button><span class="muted">' + t('failoverOffHint') + '</span></div></div>';
  document.querySelector('#save-settings').onclick = () => applyAction({
    action:'put-settings', base_url:val('#base'), log_retention:Number(val('#keep')),
    stream_idle_seconds:Number(val('#idle')), failover:document.querySelector('#failover').checked,
    ignore_proxy:document.querySelector('#ignore-proxy').checked,
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
function renderLogs() {
  const logs = filteredLogs();
  const summary = computeStats(logs);
  const limit = Number(val('#log-limit') || 100);
  const rows = logs.slice(0, limit);
  document.querySelector('#log-summary').innerHTML =
    '<span>' + t('sumRequests') + '<b>' + num(summary.requests) + '</b></span>'
    + '<span>' + t('sumInput') + '<b>' + compact(summary.prompt_tokens) + '</b></span>'
    + '<span>' + t('sumOutput') + '<b>' + compact(summary.completion_tokens) + '</b></span>'
    + '<span>' + t('sumCached') + '<b>' + compact(summary.cached_tokens) + '</b></span>'
    + '<span>' + t('sumCost') + '<b>' + money(summary.total_cost) + '</b></span>'
    + '<span class="' + (summary.errors ? 'bad' : '') + '">' + t('sumErrors') + '<b>' + num(summary.errors) + '</b></span>'
    + '<span>' + t('sumRetries') + '<b>' + num(summary.retries) + '</b></span>'
    + '<span>' + t('sumFailovers') + '<b>' + num(summary.failovers) + '</b></span>'
    + '<span>' + t('sumProbes') + '<b>' + num(summary.probes) + '</b></span>';
  const headers = '<tr><th></th><th>' + t('colTime') + '</th><th>' + t('colModel') + '</th><th>' + t('colCredential') + '</th><th>' + t('colStatus') + '</th><th>' + t('colMode') + '</th><th class="num">' + t('colTtfb') + '</th><th class="num">' + t('colDuration') + '</th><th class="num">' + t('colSpeed') + '</th><th class="num">' + t('colInput') + '</th><th class="num">' + t('colOutput') + '</th><th class="num">' + t('colCached') + '</th><th class="num">' + t('colCacheRate') + '</th><th class="num">' + t('colCost') + '</th><th></th></tr>';
  const body = rows.map(item => logRowHtml(item)).join('');
  document.querySelector('#log-table').innerHTML = '<table><thead>' + headers + '</thead><tbody>'
    + (body || '<tr><td colspan="15" class="muted">' + t('noLogs') + '</td></tr>') + '</tbody></table>';
  document.querySelectorAll('[data-log-toggle]').forEach(row => {
    row.onclick = event => {
      if (event.target.closest('button')) return;
      toggleLogDetail(row.dataset.logToggle);
    };
  });
  document.querySelectorAll('[data-log-json]').forEach(btn => {
    btn.onclick = event => {
      event.stopPropagation();
      copyLogJson(btn.dataset.logJson);
    };
  });
  document.querySelectorAll('[data-filter-model]').forEach(btn => {
    btn.onclick = event => {
      event.stopPropagation();
      modelFilter = btn.dataset.filterModel;
      const sel = document.querySelector('#model-filter');
      if (sel) sel.value = modelFilter;
      renderOverview();
      renderLogs();
    };
  });
  document.querySelectorAll('[data-filter-cred]').forEach(btn => {
    btn.onclick = event => {
      event.stopPropagation();
      credFilter = btn.dataset.filterCred;
      const sel = document.querySelector('#cred-filter');
      if (sel) sel.value = credFilter;
      renderOverview();
      renderLogs();
    };
  });
}
function logRowHtml(item) {
  const ok = item.status >= 200 && item.status < 400;
  const input = Number(item.prompt_tokens || 0);
  const cached = Number(item.cached_tokens || 0);
  const rate = input > 0 ? cached / input : null;
  const id = item.request_id || (item.time + '|' + item.model);
  const mode = item.stream ? t('detailStreamYes') : t('detailStreamNo');
  const retry = item.attempts > 1 ? ' <span class="badge warn">' + item.attempts + 'x</span>' : '';
  const failover = item.failover ? ' <span class="badge indigo">⇄</span>' : '';
  const probe = item.probe ? ' <span class="badge probe">' + t('probeTag') + '</span>' : '';
  return '<tr class="log-row" data-log-toggle="' + esc(id) + '">'
    + '<td class="nowrap"><span class="chevron">▶</span></td>'
    + '<td class="mono nowrap">' + esc(timeText(item.time)) + '</td>'
    + '<td class="model-cell" title="' + esc(item.model || '') + '">' + esc(item.model) + '</td>'
    + '<td class="nowrap" title="' + esc(item.credential || '') + '">' + esc(item.credential) + '</td>'
    + '<td class="nowrap"><span class="badge ' + (ok ? 'ok' : 'bad') + '">' + item.status + '</span>' + probe + '</td>'
    + '<td class="nowrap muted">' + mode + retry + failover + '</td>'
    + '<td class="num">' + sec(item.ttfb_ms) + '</td>'
    + '<td class="num">' + sec(item.duration_ms) + '</td>'
    + '<td class="num">' + (Number(item.tokens_per_second || 0) > 0 ? Number(item.tokens_per_second).toFixed(0) : '—') + '</td>'
    + '<td class="num">' + tokens(item.prompt_tokens) + '</td>'
    + '<td class="num">' + tokens(item.completion_tokens) + '</td>'
    + '<td class="num" title="' + esc(cacheKeyHint(item)) + '">' + tokens(item.cached_tokens) + '</td>'
    + '<td class="num">' + (rate === null ? '—' : pct1(rate)) + '</td>'
    + '<td class="num">' + money(itemCost(item)) + '</td>'
    + '<td class="nowrap"><button class="ghost" data-log-json="' + esc(id) + '" title="' + t('copyJson') + '">JSON</button></td></tr>'
    + '<tr class="log-detail" data-detail="' + esc(id) + '" hidden><td colspan="15">' + logDetailHtml(item) + '</td></tr>';
}
function logDetailHtml(item) {
  const input = Number(item.prompt_tokens || 0);
  const cached = Number(item.cached_tokens || 0);
  const completion = Number(item.completion_tokens || 0);
  const uncached = Math.max(0, input - cached);
  const rate = input > 0 ? cached / input : 0;
  const generation = Math.max(0, Number(item.duration_ms || 0) - Number(item.ttfb_ms || 0));
  const id = item.request_id || (item.time + '|' + item.model);
  const ok = item.status >= 200 && item.status < 400;
  const facts = (list) => '<div class="detail-facts">' + list.map(pair => '<div class="detail-fact"><span>' + pair[0] + '</span><b>' + pair[1] + '</b></div>').join('') + '</div>';
  const cards = [
    '<div class="detail-card"><div class="detail-title">' + t('detailRequest') + '</div>' + facts([
      [t('detailRequestId'), esc(item.request_id || '—')],
      [t('detailTime'), esc(fullTime(item.time))],
      [t('detailFormat'), esc(item.format || '—')],
      [t('detailStream'), item.stream ? t('detailStreamYes') : t('detailStreamNo')],
      [t('detailAttempts'), String(item.attempts || 1)],
      [t('detailFailover'), item.failover ? t('yes') : t('no')],
      [t('detailInputBytes'), formatBytes(item.input_bytes)],
      [t('detailOutputBytes'), formatBytes(item.output_bytes)]
    ]) + '</div>',
    '<div class="detail-card"><div class="detail-title">' + t('detailModel') + '</div>' + facts([
      [t('detailClientModel'), esc(item.model || '—')],
      [t('detailUpstreamModel'), esc(item.upstream || '—')],
      [t('detailReasoning'), esc(item.reasoning || t('passthrough'))],
      [t('detailCredential'), esc(item.credential || '—')]
    ]) + '</div>',
    '<div class="detail-card"><div class="detail-title">' + t('detailTokens') + '</div>' + facts([
      [t('detailPromptTokens'), num(input)],
      [t('detailCachedTokens'), num(cached)],
      [t('detailUncachedTokens'), num(uncached)],
      [t('detailCompletionTokens'), num(completion)],
      [t('detailTotalTokens'), num(input + completion)],
      [t('detailCacheRate'), input > 0 ? pct1(rate) : '—'],
      [t('detailCacheSource'), esc(cacheKeyHint(item))],
      [t('detailCacheKey'), esc(item.cache_key || '—')],
      [t('detailPrefixHash'), esc(item.prefix_hash || '—')],
      [t('detailCost'), money(itemCost(item))]
    ]) + '</div>',
    '<div class="detail-card"><div class="detail-title">' + t('detailPerf') + '</div>' + facts([
      [t('detailTtfb'), sec(item.ttfb_ms)],
      [t('detailGeneration'), sec(generation)],
      [t('detailTotalTime'), sec(item.duration_ms)],
      [t('detailSpeed'), Number(item.tokens_per_second || 0) > 0 ? Number(item.tokens_per_second).toFixed(1) + ' tok/s' : '—']
    ]) + '</div>'
  ].join('');
  const error = !ok && item.error
    ? '<div class="detail-error"><div class="detail-title">' + t('detailError') + '</div><pre>' + esc(item.error) + '</pre><div class="hint">' + t('errorHint') + '</div></div>'
    : '';
  const actions = '<div class="detail-actions">'
    + '<button data-log-json="' + esc(id) + '">' + t('copyJson') + '</button>'
    + (item.model ? '<button data-filter-model="' + esc(item.model) + '">' + t('filterByModel') + '</button>' : '')
    + (item.credential ? '<button data-filter-cred="' + esc(item.credential) + '">' + t('filterByCredential') + '</button>' : '')
    + '</div>';
  return '<div class="detail-grid">' + cards + '</div>' + error + actions;
}
function toggleLogDetail(id) {
  const detail = document.querySelector('[data-detail="' + cssEscape(id) + '"]');
  const row = document.querySelector('[data-log-toggle="' + cssEscape(id) + '"]');
  if (!detail) return;
  detail.hidden = !detail.hidden;
  if (row) row.classList.toggle('expanded', !detail.hidden);
}
function cssEscape(value) {
  return String(value).replace(/["\\]/g, '\\$&');
}
function copyLogJson(id) {
  const item = allLogs().find(entry => (entry.request_id || (entry.time + '|' + entry.model)) === id);
  if (!item) return;
  const text = JSON.stringify(item, null, 2);
  const done = () => { msg.textContent = t('copied'); };
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text).then(done).catch(() => fallbackCopy(text, done));
  } else {
    fallbackCopy(text, done);
  }
}
function fallbackCopy(text, done) {
  const area = document.createElement('textarea');
  area.value = text;
  area.style.position = 'fixed';
  area.style.opacity = '0';
  document.body.appendChild(area);
  area.select();
  try { document.execCommand('copy'); done(); } catch (error) {}
  document.body.removeChild(area);
}
function exportCsv() {
  const rows = filteredLogs();
  const header = ['time','request_id','model','upstream','credential','format','status','ttfb_ms','duration_ms','tokens_per_second','prompt_tokens','cached_tokens','completion_tokens','total_tokens','cache_hit_rate','cost_usd','error'];
  const lines = [header.join(',')];
  rows.forEach(item => {
    const total = (item.prompt_tokens || 0) + (item.completion_tokens || 0);
    const values = [
      item.time, item.request_id || '', item.model || '', item.upstream || '', item.credential || '', item.format || '',
      item.status, item.ttfb_ms || 0, item.duration_ms || 0, item.tokens_per_second || 0,
      item.prompt_tokens || 0, item.cached_tokens || 0, item.completion_tokens || 0, total,
      item.cache_hit_rate || 0, itemCost(item).toFixed(6), item.error || ''
    ];
    lines.push(values.map(csvCell).join(','));
  });
  const blob = new Blob(['\ufeff' + lines.join('\r\n')], {type:'text/csv;charset=utf-8'});
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = 'mistral-requests-' + new Date().toISOString().slice(0, 10) + '.csv';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}
function csvCell(value) {
  const text = String(value === null || value === undefined ? '' : value);
  return /[",\n\r]/.test(text) ? '"' + text.replace(/"/g, '""') + '"' : text;
}
function cacheKeyHint(item) {
  const source = item.cache_key_source || '';
  if (source === 'client') return t('cacheKeyClient');
  if (source === 'affinity') return t('cacheKeyAffinity');
  if (source === 'root') return t('cacheKeyRoot');
  if (source === 'instructions') return t('cacheKeyInstructions');
  if (source === 'session') return t('cacheKeySession');
  if (source === 'user') return t('cacheKeyUser');
  return t('cacheKeyNone');
}
function val(sel) { const node = document.querySelector(sel); return node ? node.value.trim() : ''; }
async function refresh() {
  data = await api({action:'state'});
  render();
}
function scheduleAutoRefresh() {
  if (autoRefreshTimer) clearInterval(autoRefreshTimer);
  autoRefreshTimer = null;
  if (autoRefreshSeconds > 0) {
    autoRefreshTimer = setInterval(() => {
      refresh().catch(error => { msg.textContent = error.message; });
    }, autoRefreshSeconds * 1000);
  }
}
document.querySelectorAll('#time-seg button').forEach(button => button.onclick = () => {
  timeFilter = button.dataset.time;
  document.querySelectorAll('#time-seg button').forEach(item => item.classList.toggle('active', item === button));
  renderOverview();
  if (!document.querySelector('#logs').hidden) renderLogs();
});
document.querySelectorAll('#status-seg button').forEach(button => button.onclick = () => {
  statusFilter = button.dataset.status;
  document.querySelectorAll('#status-seg button').forEach(item => item.classList.toggle('active', item === button));
  renderOverview();
  if (!document.querySelector('#logs').hidden) renderLogs();
});
document.querySelectorAll('#chart-metric-seg button').forEach(button => button.onclick = () => {
  chartMetric = button.dataset.metric;
  document.querySelectorAll('#chart-metric-seg button').forEach(item => item.classList.toggle('active', item === button));
  renderChart(filteredLogs());
});
document.querySelectorAll('#chart-split-seg button').forEach(button => button.onclick = () => {
  chartSplit = button.dataset.split;
  document.querySelectorAll('#chart-split-seg button').forEach(item => item.classList.toggle('active', item === button));
  renderBreakdown(filteredLogs());
});
document.querySelectorAll('#heatmap-seg button').forEach(button => button.onclick = () => {
  heatmapMetric = button.dataset.hm;
  document.querySelectorAll('#heatmap-seg button').forEach(item => item.classList.toggle('active', item === button));
  renderHeatmap(heatmapLogs());
});
document.querySelector('#model-filter').onchange = event => {
  modelFilter = event.target.value;
  renderOverview();
  if (!document.querySelector('#logs').hidden) renderLogs();
};
document.querySelector('#cred-filter').onchange = event => {
  credFilter = event.target.value;
  renderOverview();
  if (!document.querySelector('#logs').hidden) renderLogs();
};
document.querySelector('#search').oninput = event => {
  searchQuery = event.target.value;
  renderOverview();
  if (!document.querySelector('#logs').hidden) renderLogs();
};
document.querySelector('#auto-refresh').onchange = event => {
  autoRefreshSeconds = Number(event.target.value) || 0;
  scheduleAutoRefresh();
};
document.querySelector('#refresh-btn').onclick = () => {
  refresh().catch(error => { msg.textContent = error.message; });
};
document.querySelector('#export-btn').onclick = exportCsv;
document.querySelector('#log-refresh').onclick = () => {
  refresh().catch(error => { msg.textContent = error.message; });
};
document.querySelector('#log-limit').onchange = () => { renderLogs(); };
document.querySelector('#ov-logs').onclick = () => switchTab('logs');
let resizeTimer = null;
window.addEventListener('resize', () => {
  if (resizeTimer) clearTimeout(resizeTimer);
  resizeTimer = setTimeout(() => {
    renderChart(filteredLogs());
    renderHeatmap(heatmapLogs());
  }, 120);
});
if (window.ResizeObserver) {
  const heatmapNode = document.querySelector('#heatmap');
  if (heatmapNode) new ResizeObserver(() => fitHeatmap()).observe(heatmapNode);
}
applyStaticI18n();
refresh().catch(error => { msg.textContent = error.message; });
</script>`
}
