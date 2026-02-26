package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"altcha/pkg/config"
	"altcha/pkg/server"
	"altcha/pkg/store"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	s, err := initStore(cfg)
	if err != nil {
		fmt.Printf("[ALTCHA]: Failed to initialize store (%s): %v\n", cfg.Store, err)
		os.Exit(1)
	}
	defer s.Close()

	fmt.Printf("[ALTCHA]: Using %s store\n", cfg.Store)

	apiServer := server.NewAPIServer(cfg, s)
	go func() {
		addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
		fmt.Printf("[ALTCHA]: Captcha Server is running at http://localhost:%d\n", cfg.Port)
		if err := apiServer.Start(addr); err != nil {
			fmt.Printf("[ALTCHA]: API server stopped: %v\n", err)
		}
	}()

	if cfg.Demo {
		demoServer := server.NewDemoServer(cfg)
		go func() {
			addr := fmt.Sprintf("0.0.0.0:%d", cfg.DemoPort)
			fmt.Printf("[ALTCHA]: Captcha Test Server is running at http://localhost:%d\n", cfg.DemoPort)
			if err := demoServer.Start(addr); err != nil {
				fmt.Printf("[ALTCHA]: Demo server stopped: %v\n", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("[ALTCHA]: Shutting down...")
}

func initStore(cfg *config.Config) (store.Store, error) {
	switch cfg.Store {
	case "sqlite":
		return store.NewSQLiteStore(cfg.SQLitePath, cfg.MaxRecords)
	case "redis":
		return store.NewRedisStore(cfg.RedisURL, cfg.RedisCluster, cfg.ExpireMinutes)
	default:
		return store.NewMemoryStore(cfg.MaxRecords), nil
	}
}
