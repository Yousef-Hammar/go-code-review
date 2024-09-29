package api

import (
	"os"
	"strconv"
)

type Config struct {
	Addr string
	Port int
}

func New() *Config {
	return &Config{
		Addr: getString("addr", ":8080"),
	}
}

func getString(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if parsedValue, err := strconv.Atoi(value); err == nil {
			return parsedValue
		}
	}
	return fallback
}
