package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestLogger логує вхідні HTTP-запити
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Продовжуємо обробку запиту
		err := c.Next() // Виклик наступного обробника в ланцюжку

		// Логуємо після обробки
		latency := time.Since(start)
		status := c.Response().StatusCode() // <-- CORRECTED LINE HERE
		method := c.Method()
		path := c.Path()
		ip := c.IP()

		log.Printf("[%s] %s %s %s %d - %s",
			time.Now().Format("2006-01-02 15:04:05"),
			ip,
			method,
			path,
			status,
			latency,
		)
		return err
	}
}
