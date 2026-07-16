package main

import (
	"fmt"
	"vibetopia/internal/config"
)

func main() {
	cfg := config.Load()
	fmt.Printf("VIBETOPIA login-server starting on %s\n", cfg.ListenAddr)
	select {} // keep alive
}
