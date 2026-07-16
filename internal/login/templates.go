package login

import "html/template"

const BaseCSS = `
*,*::before,*::after{margin:0;padding:0;box-sizing:border-box}
html{height:100%;background:#08081a}
body{height:100%;background:#08081a;display:flex;align-items:center;justify-content:center;padding:20px;font-family:-apple-system,BlinkMacSystemFont,'SF Pro Display','Segoe UI',sans-serif;color:#e4e4ec;-webkit-font-smoothing:antialiased}
a{color:#9b6dff;text-decoration:none;font-weight:500}
.card{background:#0d0d22;border:1px solid rgba(106,61,255,0.08);border-radius:16px;padding:36px 28px 32px;width:100%;max-width:340px}
.card .logo{width:52px;height:52px;border-radius:14px;background:linear-gradient(135deg,#6a3dff,#3d7bff);display:flex;align-items:center;justify-content:center;margin:0 auto 20px;font-size:24px;font-weight:800;color:#fff;box-shadow:0 0 24px rgba(106,61,255,0.25)}
.card h1{font-size:20px;font-weight:700;text-align:center;color:#f2f2f6}
.card .sub{text-align:center;font-size:13px;color:#5c5c80;margin:6px 0 28px}
.field{margin-bottom:16px}
.field label{display:block;font-size:12px;font-weight:600;color:#5c5c80;margin-bottom:6px;text-transform:uppercase;letter-spacing:.5px}
.field input{width:100%;padding:12px 14px;font-size:15px;font-family:inherit;color:#e4e4ec;background:#060614;border:1px solid rgba(106,61,255,0.10);border-radius:10px;outline:none;transition:border-color .2s,box-shadow .2s;-webkit-appearance:none;appearance:none}
.field input::placeholder{color:#2a2a40}
.field input:focus{border-color:rgba(106,61,255,0.35);box-shadow:0 0 0 3px rgba(106,61,255,0.06)}
.btn{display:block;width:100%;padding:13px;font-size:15px;font-weight:600;font-family:inherit;color:#fff;background:linear-gradient(135deg,#6a3dff,#3d7bff);border:none;border-radius:10px;cursor:pointer;transition:opacity .2s,transform .1s;margin-top:20px}
.btn:hover{opacity:0.9}
.btn:active{transform:scale(0.98)}
.link-row{text-align:center;margin-top:18px;font-size:13px;color:#5c5c80}
footer{text-align:center;margin-top:18px;font-size:11px;color:#252540}
.alert{background:rgba(239,68,68,0.06);border:1px solid rgba(239,68,68,0.12);border-radius:8px;padding:10px 14px;text-align:center;font-size:13px;color:#ef4444;margin-top:12px}
.success{background:rgba(16,185,129,0.06);border:1px solid rgba(16,185,129,0.12);border-radius:8px;padding:10px 14px;text-align:center;font-size:13px;color:#10b981;margin-top:12px}
`

type LoginData struct {
	Error string
}

type RegisterData struct {
	Error   string
	Success string
}

type TokenData struct {
	Token  string
	GrowID string
}

var (
	TmplLogin    *template.Template
	TmplRegister *template.Template
	TmplToken    *template.Template
)

func InitTemplates() {
	TmplLogin = template.Must(template.New("login").Parse(LoginHTML))
	TmplRegister = template.Must(template.New("register").Parse(RegisterHTML))
	TmplToken = template.Must(template.New("token").Parse(TokenHTML))
}

const LoginHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no,viewport-fit=cover">
<meta name="color-scheme" content="dark">
<title>VIBETOPIA</title>
<style>` + BaseCSS + `</style>
</head>
<body>
<div>
<div class="card">
<div class="logo">V</div>
<h1>VIBETOPIA</h1>
<p class="sub">Vibe Coded GTPS</p>
<form id="loginForm" method="POST" action="/player/growid/login/validate">
<input type="hidden" name="_fmt" value="0">
<div class="field">
<label for="gid">GrowID</label>
<input type="text" name="growId" id="gid" autocomplete="username" placeholder="Enter your GrowID" required autofocus>
</div>
<div class="field">
<label for="pw">Password</label>
<input type="password" name="password" id="pw" autocomplete="current-password" placeholder="Enter your password" required>
</div>
<button type="submit" class="btn">Log In</button>
</form>
{{if .Error}}<div class="alert">{{.Error}}</div>{{end}}
<div class="link-row"><a href="/register">Create an account</a></div>
</div>
<footer>v5.50 &middot; VIBETOPIA</footer>
</div>
<script>
document.getElementById('loginForm').addEventListener('submit', async function(e) {
  e.preventDefault();
  const growId = document.getElementById('gid').value.trim();
  const password = document.getElementById('pw').value.trim();
  if (!growId || !password) return;
  try {
    const params = new URLSearchParams({growId, password});
    const resp = await fetch('/player/growid/login/validate?fmt=1', {
      method: 'POST', headers: {'Content-Type': 'application/x-www-form-urlencoded'}, body: params.toString()
    });
    const data = await resp.json();
    if (data.status === 'success' && data.token) {
      window.location = 'growtopia://login?token=' + encodeURIComponent(data.token) + '&growId=' + growId;
    } else {
      document.getElementById('errorBox').innerHTML = '<div class="alert">' + (data.message || 'Unknown error') + '</div>';
    }
  } catch(err) {
    var eb = document.createElement('div'); eb.className = 'alert'; eb.textContent = 'Connection error';
    document.querySelector('.card').appendChild(eb);
  }
});
</script>
</body>
</html>`

const RegisterHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no,viewport-fit=cover">
<meta name="color-scheme" content="dark">
<title>VIBETOPIA &middot; Register</title>
<style>` + BaseCSS + `</style>
</head>
<body>
<div>
<div class="card">
<div class="logo">+</div>
<h1>New Account</h1>
<p class="sub">Join VIBETOPIA</p>
<form method="POST" action="/register">
<div class="field"><label for="gid">GrowID</label><input type="text" name="tankIDName" id="gid" autocomplete="username" placeholder="Choose a GrowID" required autofocus></div>
<div class="field"><label for="pw">Password</label><input type="password" name="tankIDPass" id="pw" autocomplete="new-password" placeholder="Choose a password" required></div>
<div class="field"><label for="dn">Display Name</label><input type="text" name="displayName" id="dn" placeholder="Optional display name"></div>
<button type="submit" class="btn">Create Account</button>
</form>
{{if .Error}}<div class="alert">{{.Error}}</div>{{end}}
{{if .Success}}<div class="success">{{.Success}}</div>{{end}}
<div class="link-row"><a href="/">Already registered? Log in</a></div>
</div>
<footer>v5.50 &middot; VIBETOPIA</footer>
</div>
</body>
</html>`

const TokenHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,user-scalable=no,viewport-fit=cover">
<meta name="gt-token" content="{{.Token}}">
<meta name="gt-growid" content="{{.GrowID}}">
<title>VIBETOPIA — {{.GrowID}}</title>
<style>
body{background:#08081a;color:#e4e4ec;font-family:-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}
.card{text-align:center;max-width:340px;padding:32px}
h1{font-size:22px;margin:0 0 4px;color:#f2f2f6}
p{font-size:13px;color:#5c5c80;margin:4px 0 20px}
.meta-box{background:#0d0d22;border:1px solid rgba(106,61,255,0.15);border-radius:10px;padding:14px 18px;margin:16px 0;text-align:left;font-family:monospace;font-size:11px;color:#9b6dff;word-break:break-all;overflow-wrap:anywhere}
.meta-box code{display:block;margin:6px 0;color:#6a3dff}
.meta-box span{color:#3d7bff}
#status{font-size:12px;color:#5c5c80;margin-top:12px}
</style>
<script>
var TOKEN = "{{.Token}}";
var GROWID = "{{.GrowID}}";
var tried = [];

function log(s){var d=document.getElementById('status');d.textContent = (d.textContent||'') + '\\n' + s;}
function tryBridge(method) {tried.push(method);log(method);}

// ═══════ BRIDGE SHOTGUN — ALL KNOWN MECHANISMS ═══════
setTimeout(function(){

// 1. setInterval poller: terus retry meta refresh
log("⏳ Shotgunning bridges...");
tryBridge("meta-refresh");

// 2. window.location (works on iOS/Android WebView)
var gurl = "growtopia://login?token="+encodeURIComponent(TOKEN)+"&growId="+GROWID;
tryBridge("window.location");
window.location = gurl;

// 3. window.open (some WebViews intercept this)
setTimeout(function(){
  tryBridge("window.open");
  window.open(gurl, "_self");
}, 500);

// 4. document.location
setTimeout(function(){
  tryBridge("document.location");
  document.location = gurl;
}, 1000);

// 5. iframe injection (some apps poll iframe src)
setTimeout(function(){
  tryBridge("iframe-src");
  var iframe = document.createElement('iframe');
  iframe.style.display = 'none';
  iframe.src = gurl;
  document.body.appendChild(iframe);
}, 1500);

// 6. window.postMessage — for native handler polling
var msg = JSON.stringify({token:TOKEN,growId:GROWID,accountType:"growtopia",url:""});
tryBridge("postMessage-*");
window.postMessage(msg, "*");

// 7. Spam semua WKWebView message handler names
var handlers = ["openInBrowser","growtopia","gtLogin","gt","ubisoft","ibml","login","auth","native","callback","handler","onLogin","ubisoftLogin","growtopiaLogin","token","tokenReceiver","tokenHandler","bridge","Growtopia","IBML","UbisoftServices","UbiServices","ubiLogin","OpenURL"];
handlers.forEach(function(h){
  try {
    if(window.webkit && window.webkit.messageHandlers && window.webkit.messageHandlers[h]){
      tryBridge("messageHandler:"+h);
      window.webkit.messageHandlers[h].postMessage({token:TOKEN,growId:GROWID,accountType:"growtopia",url:""});
    }
  }catch(e){}
});

// 8. userContentController — script injection (Mac only)
try {
  if(window.webkit && window.webkit.messageHandlers){
    tryBridge("userContentController-spam");
    // Spam ALL properties
    var ks = Object.keys(window.webkit.messageHandlers);
    ks.forEach(function(k){
      try{ window.webkit.messageHandlers[k].postMessage({token:TOKEN,growId:GROWID}); }catch(e){}
    });
    log("handlers found: "+ks.join(","));
  }
}catch(e){log("no webkit.messageHandlers");}

// 9. Cookie bridge
tryBridge("cookie");
document.cookie = "gt_token="+TOKEN+"; path=/; max-age=120; SameSite=Lax";
document.cookie = "gt_growId="+GROWID+"; path=/; max-age=120; SameSite=Lax";
document.cookie = "_token="+TOKEN+"; path=/; max-age=120; SameSite=Lax";

// 10. sessionStorage + localStorage
tryBridge("storage");
try{sessionStorage.setItem("vibetopia_token", TOKEN);}catch(e){}
try{sessionStorage.setItem("vibetopia_growId", GROWID);}catch(e){}
try{localStorage.setItem("vibetopia_token", TOKEN);}catch(e){}
try{localStorage.setItem("vibetopia_growId", GROWID);}catch(e){}

// 11. BroadcastChannel (some Mac apps listen)
tryBridge("BroadcastChannel");
try{ new BroadcastChannel("growtopia").postMessage({token:TOKEN,growId:GROWID}); }catch(e){}
try{ new BroadcastChannel("gt_login").postMessage({token:TOKEN,growId:GROWID}); }catch(e){}

// 12. CustomEvent
tryBridge("CustomEvent");
try{ window.dispatchEvent(new CustomEvent("growtopia:login", {detail:{token:TOKEN,growId:GROWID}})); }catch(e){}
try{ window.dispatchEvent(new CustomEvent("gt:token", {detail:TOKEN})); }catch(e){}

// 13. Form POST ke localhost (some apps intercept localhost)
tryBridge("form-localhost");
var f = document.createElement('form');
f.method='POST';
f.action='http://127.0.0.1:17091/vibetopia/token';
f.style.display='none';
['token','growId'].forEach(function(n){
  var i = document.createElement('input');
  i.name=n; i.value=n==='token'?TOKEN:GROWID;
  f.appendChild(i);
});
document.body.appendChild(f);
setTimeout(function(){f.submit();}, 2000);

// 14. Retry periodic
var retryCount = 0;
var retryInt = setInterval(function(){
  retryCount++;
  if(retryCount > 30) { clearInterval(retryInt); log("⏹ Retries exhausted ("+retryCount+")"); return; }
  tryBridge("retry#"+retryCount);
  window.location = gurl;
}, 500);

}, 100);
</script>
</head>
<body>
<div class="card">
<h1>🎸 VIBETOPIA</h1>
<p>Welcome, {{.GrowID}}</p>
<div class="meta-box">
<code>Token:</code> <span>{{.Token}}</span><br>
<code>GrowID:</code> <span>{{.GrowID}}</span>
</div>
<p id="status">⏳ Bridging to game...</p>
<p style="font-size:10px;color:#252540;margin-top:20px">If stuck, your Mac client does not support this login method.<br>Try Android or Windows instead.</p>
</div>
</body>
</html>`

