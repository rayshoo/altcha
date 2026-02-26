package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"altcha/pkg/config"
	"altcha/pkg/server"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	apiServer := server.NewAPIServer(cfg)
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		fmt.Printf("[ALTCHA]: Captcha Server is running at http://localhost:%d\n", cfg.Port)
		if err := apiServer.Start(addr); err != nil {
			fmt.Printf("[ALTCHA]: API server stopped: %v\n", err)
		}
	}()

	if cfg.Demo {
		demoServer := server.NewDemoServer(cfg)
		go func() {
			fmt.Println("[ALTCHA]: Captcha Test Server is running at http://localhost:8080")
			if err := demoServer.Start(":8080"); err != nil {
				fmt.Printf("[ALTCHA]: Demo server stopped: %v\n", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("[ALTCHA]: Shutting down...")
}
