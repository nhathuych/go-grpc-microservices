package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/nhathuych/api-gateway/internal/config"
	"github.com/nhathuych/api-gateway/internal/proxy"
)

var startTime = time.Now()

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${locals:requestid}] ${time} | ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Use(timeout.NewWithContext(func(c *fiber.Ctx) error {
		return c.Next()
	}, 15*time.Second))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD",
	}))

	// Health check endpoint (used by Kubernetes, Docker, or monitoring systems).
	// Must not be protected by rate limiting or authentication.
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "api-gateway",
			"uptime":  time.Since(startTime).String(),
			"time":    time.Now().UTC(),
		})
	})

	// Global rate limiter (in-memory)
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too Many Requests",
				"message": "Request rate limit exceeded. Please retry later.",
			})
		},
	}))

	api := app.Group("/api/v1")

	accountProxy := proxy.ReverseProxy(cfg.AccountServiceURL)
	orderProxy := proxy.ReverseProxy(cfg.OrderServiceURL)

	api.All("/accounts/*", accountProxy)
	api.All("/orders/*", orderProxy)
}
