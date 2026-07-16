package main

import (
	"fmt"
	"vibetopia/internal/config"
)

func main() {
	cfg := config.Load()
	fmt.Printf("VIBETOPIA game-server starting on UDP :%s\n", cfg.GamePort)
	select {} // keep alive
}
