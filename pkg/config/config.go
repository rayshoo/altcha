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
	LogLevel      string
	RateLimit     float64
	Store        string
	SQLitePath   string
	RedisURL     string
	RedisCluster bool
	DemoPort     int

	// Analytics
	PostgresURL string
	GeoIPDB     string

	// Dashboard
	DashboardPort int

	// Auth
	AuthProvider              string
	AuthUsername              string
	AuthPassword              string
	AuthIssuer                string
	AuthClientID              string
	AuthClientSecret          string
	AuthPKCE                  bool
	AuthAuthorizationEndpoint string
	AuthTokenEndpoint         string
	AuthEndSessionEndpoint    string
	AuthJWKSURI               string
	AuthAllowedUsers          []string
	AuthAllowedGroups         []string
	AuthAllowedRoles          []string
}

func (c *Config) IsDebug() bool {
	return strings.EqualFold(c.LogLevel, "debug")
}

func (c *Config) AnalyticsEnabled() bool {
	return c.PostgresURL != ""
}

func Load() *Config {
	cfg := &Config{
		Port:          envInt("PORT", 3000),
		Secret:        envStr("SECRET", "$ecret.key"),
		ExpireMinutes: envInt("EXPIREMINUTES", 10),
		MaxNumber:     envInt("COMPLEXITY", 1000000),
		MaxRecords:    envInt("MAXRECORDS", 1000),
		CorsOrigin:    envList("CORS_ORIGIN", nil),
		Demo:          envBool("DEMO", false),
		LogLevel:      envStr("LOG_LEVEL", "info"),
		RateLimit:     envFloat("RATE_LIMIT", 0),
		Store:         envStr("STORE", "memory"),
		SQLitePath:    envStr("SQLITE_PATH", "data/altcha.db"),
		RedisURL:      envStr("REDIS_URL", "redis://localhost:6379"),
		RedisCluster:  envBool("REDIS_CLUSTER", false),
		DemoPort:      envInt("DEMO_PORT", 8000),

		// Analytics
		PostgresURL: envStr("POSTGRES_URL", ""),
		GeoIPDB:     envStr("GEOIP_DB", ""),

		// Dashboard
		DashboardPort: envInt("DASHBOARD_PORT", 9000),

		// Auth
		AuthProvider:              envStr("AUTH_PROVIDER", ""),
		AuthUsername:              envStr("AUTH_USERNAME", ""),
		AuthPassword:              envStr("AUTH_PASSWORD", ""),
		AuthIssuer:                envStr("AUTH_ISSUER", ""),
		AuthClientID:              envStr("AUTH_CLIENT_ID", ""),
		AuthClientSecret:          envStr("AUTH_CLIENT_SECRET", ""),
		AuthPKCE:                  envBool("AUTH_PKCE", true),
		AuthAuthorizationEndpoint: envStr("AUTH_AUTHORIZATION_ENDPOINT", ""),
		AuthTokenEndpoint:         envStr("AUTH_TOKEN_ENDPOINT", ""),
		AuthEndSessionEndpoint:    envStr("AUTH_END_SESSION_ENDPOINT", ""),
		AuthJWKSURI:               envStr("AUTH_JWKS_URI", ""),
		AuthAllowedUsers:          envList("AUTH_ALLOWED_USERS", nil),
		AuthAllowedGroups:         envList("AUTH_ALLOWED_GROUPS", nil),
		AuthAllowedRoles:          envList("AUTH_ALLOWED_ROLES", nil),
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

func envFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
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
