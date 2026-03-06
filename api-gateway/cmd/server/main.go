package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nhathuych/api-gateway/internal/api/router"
	"github.com/nhathuych/api-gateway/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	app := fiber.New(fiber.Config{
		AppName:      "Microservices API Gateway",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Setup routes
	router.SetupRoutes(app, cfg)

	// Handle graceful shutdown in a goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			log.Panic(err)
		}
	}()

	fmt.Printf("Gateway is running on port %s...\n", cfg.Port)

	// Wait for termination signal (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Shutting down Gateway...")
	_ = app.Shutdown()
}
