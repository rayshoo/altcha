package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"altcha/pkg/config"
	"altcha/pkg/dashboard"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	if cfg.PostgresURL == "" {
		fmt.Println("[DASHBOARD]: POSTGRES_URL is required")
		os.Exit(1)
	}
	if cfg.AuthProvider == "" {
		fmt.Println("[DASHBOARD]: AUTH_PROVIDER is required (basic or keycloak)")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		fmt.Printf("[DASHBOARD]: Failed to connect to postgres: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("[DASHBOARD]: Failed to ping postgres: %v\n", err)
		os.Exit(1)
	}

	srv, err := dashboard.NewServer(cfg, db)
	if err != nil {
		fmt.Printf("[DASHBOARD]: Failed to create server: %v\n", err)
		os.Exit(1)
	}

	go func() {
		addr := fmt.Sprintf("0.0.0.0:%d", cfg.DashboardPort)
		fmt.Printf("[DASHBOARD]: Dashboard is running at http://localhost:%d\n", cfg.DashboardPort)
		if err := srv.Start(addr); err != nil {
			fmt.Printf("[DASHBOARD]: Server stopped: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("[DASHBOARD]: Shutting down...")
}
