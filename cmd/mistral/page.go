package main

func consolePage() string {
	return `<!doctype html>
<meta charset="utf-8">
<title>Mistral</title>
<style>
:root { color-scheme: light; --bg:#f6f7f9; --card:#fff; --ink:#1c1917; --muted:#57534e; --line:#e7e5e4; --accent:#9a3412; }
body { margin:0; font:14px/1.45 ui-sans-serif,system-ui,sans-serif; background:var(--bg); color:var(--ink); }
header { display:flex; justify-content:space-between; align-items:center; padding:18px 28px; background:var(--card); border-bottom:1px solid var(--line); }
h1 { font-size:18px; margin:0; }
nav { display:flex; gap:6px; padding:12px 28px; }
button, .tab { border:1px solid var(--line); background:var(--card); border-radius:8px; padding:7px 12px; cursor:pointer; }
.tab.active { background:var(--ink); color:#fff; }
main { padding:0 28px 40px; max-width:980px; }
.card { background:var(--card); border:1px solid var(--line); border-radius:12px; padding:16px; margin-top:12px; }
.row { display:flex; gap:10px; align-items:center; justify-content:space-between; padding:10px 0; border-top:1px solid var(--line); }
.muted { color:var(--muted); font-size:12px; }
input, select { padding:7px 8px; border:1px solid var(--line); border-radius:8px; }
table { width:100%; border-collapse:collapse; }
td, th { text-align:left; padding:8px 6px; border-bottom:1px solid var(--line); font-size:13px; }
.err { color:#9f1239; }
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
<section id="logs" class="card" hidden></section>
<section id="config" class="card" hidden></section>
<p id="msg" class="muted"></p>
</main>
<script>
const msg = document.querySelector('#msg');
let data = null;
let managementKey = readPanelKey();
function readPanelKey() {
  const raw = localStorage.getItem('cli-proxy-auth');
  if (!raw) return '';
  let parsed = null;
  if (raw.indexOf('enc::v2::') === 0 || raw.indexOf('enc::v1::') === 0) {
    const v1 = raw.indexOf('enc::v1::') === 0;
    const secret = 'cli-proxy-api-webui::secure-storage|' + location.host + (v1 ? '|' + navigator.userAgent : '');
    try {
      const binary = atob(raw.slice(9));
      const key = new TextEncoder().encode(secret);
      const out = new Uint8Array(binary.length);
      for (let i = 0; i < binary.length; i++) out[i] = binary.charCodeAt(i) ^ key[i % key.length];
      parsed = JSON.parse(new TextDecoder().decode(out));
    } catch (error) { return ''; }
  } else {
    try { parsed = JSON.parse(raw); } catch (error) { return ''; }
  }
  const state = parsed && parsed.state ? parsed.state : parsed;
  return state && state.managementKey ? state.managementKey : '';
}
document.querySelectorAll('.tab').forEach(tab => tab.onclick = () => {
  document.querySelectorAll('.tab').forEach(item => item.classList.toggle('active', item === tab));
  document.querySelectorAll('main section').forEach(section => section.hidden = section.id !== tab.dataset.tab);
  if (tab.dataset.tab === 'logs') loadLogs();
});
async function api(body) {
  if (!managementKey) managementKey = prompt('面板没有记住管理密钥。登录时勾选记住密码后再打开，或在这里粘贴一次。') || '';
  const response = await fetch('/v0/management/mistral/api', {
    method:'POST', headers:{'content-type':'application/json','authorization':'Bearer '+managementKey},
    body: JSON.stringify(body)
  });
  const text = await response.text();
  if (!response.ok) throw new Error(response.status+' '+text);
  return JSON.parse(text);
}
function render() {
  document.querySelector('#state').textContent = data.configured ? '可调度' : '还不能调度';
  document.querySelector('#overview').innerHTML = '<b>'+data.base_url+'</b><div class="muted">启用凭据按优先级取最高的一把。思考块会转成各协议自己的字段。</div>'
    + '<div class="row"><span>凭据</span><span>'+data.credentials.length+'</span></div>'
    + '<div class="row"><span>模型映射</span><span>'+data.models.length+'</span></div>'
    + '<div class="row"><span>流空闲截断</span><span>'+(data.stream_idle_seconds||0)+' 秒</span></div>';
  document.querySelector('#credentials').innerHTML = '<b>凭据</b>' + data.credentials.map(item =>
    '<div class="row"><div><b>'+item.label+'</b><div class="muted">'+item.id+' · '+item.mask+' · 优先级 '+item.priority+(item.disabled?' · 已停用':'')+'</div></div><div>'
    + '<button data-toggle="'+item.id+'">'+(item.disabled?'启用':'停用')+'</button> <button data-del="'+item.id+'">删除</button></div></div>').join('')
    + '<div class="row"><input id="c-label" placeholder="名称"><input id="c-key" placeholder="API key" type="password"><input id="c-pri" type="number" value="10" style="width:70px"><button id="c-add">添加</button></div>';
  document.querySelector('#models').innerHTML = '<b>模型映射</b><div class="muted">reasoning 留空则沿用客户端的思考档位。</div><table><tr><th>客户端</th><th>上游</th><th>思考</th><th></th></tr>'
    + data.models.map(item => '<tr><td>'+item.id+'</td><td>'+item.upstream_id+'</td><td>'+(item.reasoning||'透传')+'</td><td><button data-mdel="'+item.id+'">删除</button></td></tr>').join('')
    + '</table><div class="row"><input id="m-id" placeholder="客户端模型名"><input id="m-up" placeholder="上游模型 ID"><select id="m-reason"><option value="">透传</option><option>max</option><option>high</option><option>medium</option><option>low</option></select><button id="m-add">添加</button></div>';
  document.querySelector('#config').innerHTML = '<b>配置</b><div class="row"><label>Base URL</label><input id="base" value="'+data.base_url+'" style="width:360px"></div>'
    + '<div class="row"><label>日志保留</label><input id="keep" type="number" value="'+data.log_retention+'" style="width:80px"></div>'
    + '<div class="row"><label>流空闲秒数，0 为关闭</label><input id="idle" type="number" value="'+(data.stream_idle_seconds||0)+'" style="width:80px"></div>'
    + '<button id="save-settings">保存配置</button>';
  bind();
}
function bind() {
  const add = document.querySelector('#c-add');
  if (add) add.onclick = () => run(() => api({action:'put-credential', label:val('#c-label'), api_key:val('#c-key'), priority:Number(val('#c-pri')||10)}));
  document.querySelectorAll('[data-del]').forEach(btn => btn.onclick = () => run(() => api({action:'delete-credential', id:btn.dataset.del})));
  document.querySelectorAll('[data-toggle]').forEach(btn => btn.onclick = () => {
    const item = data.credentials.find(entry => entry.id === btn.dataset.toggle);
    run(() => api({action:'put-credential', id:item.id, label:item.label, disabled:!item.disabled, priority:item.priority}));
  });
  const madd = document.querySelector('#m-add');
  if (madd) madd.onclick = () => run(() => api({action:'put-model', id:val('#m-id'), upstream_id:val('#m-up'), reasoning:val('#m-reason')}));
  document.querySelectorAll('[data-mdel]').forEach(btn => btn.onclick = () => run(() => api({action:'delete-model', id:btn.dataset.mdel})));
  const save = document.querySelector('#save-settings');
  if (save) save.onclick = () => run(() => api({action:'put-settings', base_url:val('#base'), log_retention:Number(val('#keep')), stream_idle_seconds:Number(val('#idle'))}));
}
function val(sel) { const node = document.querySelector(sel); return node ? node.value.trim() : ''; }
async function run(fn) { msg.textContent=''; try { data = await fn(); render(); } catch (error) { msg.textContent = error.message; } }
async function loadLogs() {
  const result = await api({action:'logs', limit:40});
  document.querySelector('#logs').innerHTML = '<b>请求日志</b><table><tr><th>时间</th><th>模型</th><th>凭据</th><th>状态</th><th>首字</th><th>耗时</th><th>错误</th></tr>'
    + (result.logs||[]).map(item => '<tr><td>'+item.time.slice(11,19)+'</td><td>'+item.model+'</td><td>'+item.credential+'</td><td>'+item.status+'</td><td>'+item.ttfb_ms+'</td><td>'+item.duration_ms+'</td><td class="err">'+(item.error||'')+'</td></tr>').join('')
    + '</table>';
}
run(() => api({action:'state'})).catch(error => { msg.textContent = error.message; });
</script>`
}
