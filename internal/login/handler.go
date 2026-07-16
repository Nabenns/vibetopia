package login

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	InitTemplates()
	return &Handler{DB: db}
}

// ─── Token ───────────────────────────────────────────────────────────────────

func generateToken(growID string) string {
	raw := fmt.Sprintf("_token=%d&growId=%s", time.Now().Unix(), growID)
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func safeGrowID(s string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, s)
}

func playerExists(db *sql.DB, growID string) bool {
	var id int64
	err := db.QueryRow("SELECT id FROM players WHERE LOWER(growid) = LOWER($1)", growID).Scan(&id)
	return err == nil
}

func checkPassword(db *sql.DB, growID, password string) bool {
	var pw string
	err := db.QueryRow("SELECT password FROM players WHERE LOWER(growid) = LOWER($1)", growID).Scan(&pw)
	return err == nil && pw == password
}

// ─── HTTP Handlers ───────────────────────────────────────────────────────────

func (h *Handler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	growID := r.URL.Query().Get("growId")
	if growID == "" {
		growID = r.URL.Query().Get("user")
	}

	// Token page — shown after successful login via 302 redirect
	if token != "" && growID != "" {
		log.Printf("[vibetopia/dashboard] TOKEN PAGE: growId=%s", growID)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		TmplToken.Execute(w, TokenData{Token: token, GrowID: growID})
		return
	}

	// Login form with optional error
	errorMsg := ""
	errType := r.URL.Query().Get("error")
	switch errType {
	case "wrong_password":
		errorMsg = "Wrong password for " + r.URL.Query().Get("growId") + "."
	case "not_found":
		errorMsg = "Account '" + r.URL.Query().Get("growId") + "' not found."
	case "session_expired":
		errorMsg = "Session expired. Please login again."
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	TmplLogin.Execute(w, LoginData{Error: errorMsg})
}

func (h *Handler) HandleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{
			"status": "error", "message": "Method not allowed",
			"token": "", "url": "", "accountType": "growtopia",
		})
		return
	}

	r.ParseForm()
	growID := strings.TrimSpace(r.FormValue("growId"))
	password := r.FormValue("password")
	log.Printf("[vibetopia/validate] growId=%s", growID)

	if growID == "" || password == "" {
		writeJSON(w, map[string]string{
			"status": "error", "message": "GrowID and password required.",
			"token": "", "url": "", "accountType": "growtopia",
		})
		return
	}

	if !playerExists(h.DB, growID) {
		log.Printf("[vibetopia/validate] NOT FOUND: %s", growID)
		http.Redirect(w, r, "/player/login/dashboard?error=not_found&growId="+growID, http.StatusFound)
		return
	}

	if !checkPassword(h.DB, growID, password) {
		log.Printf("[vibetopia/validate] WRONG PASSWORD: %s", growID)
		http.Redirect(w, r, "/player/login/dashboard?error=wrong_password&growId="+growID, http.StatusFound)
		return
	}

	token := generateToken(growID)
	log.Printf("[vibetopia/validate] SUCCESS: %s", growID)

	// Update last_login
	h.DB.Exec("UPDATE players SET last_login = NOW(), last_ip = $1 WHERE LOWER(growid) = LOWER($2)",
		extractIP(r), growID)

	fmtMode := r.URL.Query().Get("fmt")
	if fmtMode == "" {
		fmtMode = r.FormValue("_fmt")
	}
	if fmtMode == "" {
		fmtMode = "1" // Default to JSON for JS fetch
	}

	switch fmtMode {
	case "1":
		writeJSON(w, map[string]string{
			"status": "success", "message": "Account Validated.",
			"token": token, "url": "", "accountType": "growtopia",
		})
	case "2":
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "requestedName|%s\ntankIDName|%s\ntankIDPass|%s\n_token=%d&growId=%s&password=%s&reg=0",
			growID, growID, password, time.Now().Unix(), growID, password)
	default:
		log.Printf("[vibetopia/validate] 302 → dashboard growId=%s", growID)
		http.Redirect(w, r, fmt.Sprintf("/player/login/dashboard?token=%s&growId=%s",
			url.QueryEscape(token), growID), http.StatusFound)
	}
}

func (h *Handler) HandleCheckToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	refreshToken := r.FormValue("refreshToken")
	if refreshToken == "" {
		refreshToken = r.FormValue("token")
	}
	if refreshToken == "" {
		refreshToken = r.FormValue("requestData")
	}

	// Fallback: scan POST body
	if refreshToken == "" && r.Body != nil {
		buf, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(strings.NewReader(string(buf)))
		body := string(buf)
		for _, field := range []string{"refreshToken=", "token=", "requestData="} {
			if idx := strings.Index(body, field); idx >= 0 {
				end := strings.Index(body[idx:], "&")
				if end < 0 {
					end = len(body) - idx
				}
				refreshToken = body[idx+len(field):][:end-len(field)]
				break
			}
		}
	}

	var growID string
	if refreshToken != "" {
		if decoded, err := url.QueryUnescape(refreshToken); err == nil {
			refreshToken = decoded
		}
		if padding := len(refreshToken) % 4; padding != 0 {
			refreshToken += strings.Repeat("=", 4-padding)
		}
		if decoded, err := base64.StdEncoding.DecodeString(refreshToken); err == nil {
			parts := strings.Split(string(decoded), "&")
			for _, p := range parts {
				kv := strings.SplitN(p, "=", 2)
				if len(kv) == 2 && kv[0] == "growId" {
					growID = kv[1]
				}
			}
		}
	}

	if growID == "" {
		growID = "guest"
	}
	if !playerExists(h.DB, growID) {
		growID = "guest"
	}

	newToken := generateToken(growID)
	log.Printf("[vibetopia/checktoken] growID=%s", growID)

	http.Redirect(w, r, fmt.Sprintf("/player/login/dashboard?token=%s&growId=%s",
		url.QueryEscape(newToken), growID), http.StatusFound)
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		TmplRegister.Execute(w, RegisterData{})
		return
	}

	r.ParseForm()
	growID := strings.TrimSpace(r.FormValue("tankIDName"))
	password := r.FormValue("tankIDPass")
	displayName := strings.TrimSpace(r.FormValue("displayName"))

	if growID == "" || password == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		TmplRegister.Execute(w, RegisterData{Error: "GrowID and password are required."})
		return
	}

	safeID := safeGrowID(growID)
	if safeID != growID || len(growID) < 3 || len(growID) > 24 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		TmplRegister.Execute(w, RegisterData{Error: "GrowID must be 3-24 letters/numbers."})
		return
	}

	if len(password) < 3 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		TmplRegister.Execute(w, RegisterData{Error: "Password must be at least 3 characters."})
		return
	}

	if playerExists(h.DB, growID) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		TmplRegister.Execute(w, RegisterData{Error: "This GrowID is already taken."})
		return
	}

	if displayName == "" {
		displayName = growID
	}

	_, err := h.DB.Exec(
		"INSERT INTO players (growid, password, display) VALUES ($1, $2, $3)",
		growID, password, displayName,
	)
	if err != nil {
		log.Printf("[vibetopia/register] ERROR: %v", err)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		TmplRegister.Execute(w, RegisterData{Error: "Server error. Try again."})
		return
	}

	// Give starter items
	var newPlayerID int64
	h.DB.QueryRow("SELECT id FROM players WHERE growid = $1", growID).Scan(&newPlayerID)
	if newPlayerID > 0 {
		starterItems := map[int]int{
			2: 1,    // Dirt
			18: 1,   // Door
			32: 1,   // Wood Block
			112: 1,  // Lock
			242: 1,  // Fist
			1388: 1, // Torch
		}
		for itemID, qty := range starterItems {
			h.DB.Exec("INSERT INTO player_inventory (player_id, item_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (player_id, item_id) DO NOTHING",
				newPlayerID, itemID, qty)
		}
	}

	log.Printf("[vibetopia/register] NEW: %s (%s)", growID, displayName)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	TmplRegister.Execute(w, RegisterData{
		Success: `Account created! <a href="/" style="color:#10b981;font-weight:600;">Go to Login</a>`,
	})
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func extractIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	if fwd := r.Header.Get("X-Real-IP"); fwd != "" {
		return fwd
	}
	// Strip port from RemoteAddr
	host := r.RemoteAddr
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}
	return host
}
