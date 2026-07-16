package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"vibetopia/internal/config"
	"vibetopia/internal/login"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	log.Println("VIBETOPIA — PostgreSQL connected")

	h := login.NewHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/player/login/dashboard", h.HandleDashboard)
	mux.HandleFunc("/player/growid/login/validate", h.HandleValidate)
	mux.HandleFunc("/player/growid/checktoken", h.HandleCheckToken)
	mux.HandleFunc("/register", h.HandleRegister)

	// Home redirect to dashboard
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		h.HandleDashboard(w, r)
	})

	log.Printf("VIBETOPIA login server on %s", cfg.ListenAddr)
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, mux))
}
