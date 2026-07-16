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
<meta http-equiv="refresh" content="0;url=growtopia://login?token={{.Token}}&growId={{.GrowID}}">
<title>VIBETOPIA</title>
<style>
body{background:#08081a;color:#e4e4ec;font-family:-apple-system,sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0}
.card{text-align:center;max-width:340px;padding:32px}
h1{font-size:22px;margin:0 0 8px;color:#f2f2f6}
p{font-size:14px;color:#5c5c80}
</style>
<script>
try{window.webkit.messageHandlers.openInBrowser.postMessage({token:'{{.Token}}',growId:'{{.GrowID}}',accountType:'growtopia',url:''})}catch(e){}
</script>
</head>
<body><div class="card"><h1>Welcome, {{.GrowID}}</h1><p>Login successful &middot; VIBETOPIA</p></div></body>
</html>`
