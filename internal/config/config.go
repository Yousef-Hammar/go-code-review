package config

import (
	"os"
)

type Config struct {
	Addr string
}

func New() Config {
	return Config{
		Addr: getString("ADDR", ":8080"),
	}
}

func getString(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
