package config

import "os"

type Config struct {
	DatabaseURL  string
	RedisURL     string
	ListenAddr   string
	GamePort     string
	GameServerIP string
}

func Load() *Config {
	return &Config{
		DatabaseURL:  envOrDefault("DATABASE_URL", "postgres://vibetopia:vibetopia123@localhost:5432/vibetopia?sslmode=disable"),
		RedisURL:     envOrDefault("REDIS_URL", "redis://localhost:6379"),
		ListenAddr:   envOrDefault("LISTEN_ADDR", ":8080"),
		GamePort:     envOrDefault("GAME_PORT", "17091"),
		GameServerIP: envOrDefault("GAME_SERVER_IP", "103.253.213.178"),
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
