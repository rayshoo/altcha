package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port          int
	Secret        string
	ExpireMinutes int
	MaxNumber     int
	MaxRecords    int
	CorsOrigin    []string
	Demo          bool
}

func Load() *Config {
	cfg := &Config{
		Port:          envInt("PORT", 3000),
		Secret:        envStr("SECRET", "$ecret.key"),
		ExpireMinutes: envInt("EXPIREMINUTES", 10),
		MaxNumber:     envInt("MAXNUMBER", 1000000),
		MaxRecords:    envInt("MAXRECORDS", 1000),
		CorsOrigin:    envList("CORS_ORIGIN", nil),
		Demo:          envBool("DEMO", false),
	}

	if cfg.Secret == "$ecret.key" {
		fmt.Println(" [WARNING] CHANGE ALTCHA SECRET KEY - its still default !!! ")
	}

	return cfg
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		return strings.EqualFold(v, "true")
	}
	return fallback
}

func envList(key string, fallback []string) []string {
	if v := os.Getenv(key); v != "" {
		parts := strings.Split(v, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	return fallback
}
