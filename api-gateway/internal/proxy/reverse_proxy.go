package proxy

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

// ReverseProxy forwards incoming requests to the specified target service.
// It uses Fiber's built-in proxy balancer to route traffic to the upstream server.
func ReverseProxy(target string) fiber.Handler {
	return proxy.Balancer(proxy.Config{
		Servers: []string{target},
		Timeout: 10 * time.Second,
	})
}
