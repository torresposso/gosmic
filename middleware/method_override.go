package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

// MethodOverride middleware allows method overriding via header or form value
func MethodOverride() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Only override on POST requests
		if c.Method() != fiber.MethodPost {
			return c.Next()
		}

		// Check query/form value "_method"
		method := c.FormValue("_method")
		if method == "" {
			// Check header "X-HTTP-Method-Override"
			method = c.Get("X-HTTP-Method-Override")
		}

		if method != "" {
			// Validate method
			upperMethod := strings.ToUpper(method)
			switch upperMethod {
			case fiber.MethodPut, fiber.MethodDelete, fiber.MethodPatch:
				c.Method(upperMethod)
			}
		}

		return c.Next()
	}
}
