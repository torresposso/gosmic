package middleware

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/torresposso/gosmic/pb"

	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware validates the PocketBase token from cookie and creates a request-scoped client.
func AuthMiddleware(globalClient *pb.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("pb_auth")
		if token == "" {
			return c.Redirect().To("/login")
		}

		// Create a request-scoped client with the user's token
		userClient := globalClient.WithToken(token)

		// Extract User ID from JWT (Payload is the 2nd part)
		// We trust PocketBase to verify the signature on the next request.
		// We just need the ID to form the request correctly.
		parts := strings.Split(token, ".")
		if len(parts) == 3 {
			payload, _ := base64.RawURLEncoding.DecodeString(parts[1])
			var claims struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(payload, &claims); err == nil && claims.ID != "" {
				userClient.AuthRecord = &pb.User{ID: claims.ID}
			}
		}

		// Store in context for handlers to use
		c.Locals("pb", userClient)

		return c.Next()
	}
}

// GetPBClient retrieves the request-scoped PocketBase client from context.
// Returns nil if not authenticated.
func GetPBClient(c fiber.Ctx) *pb.Client {
	if client, ok := c.Locals("pb").(*pb.Client); ok {
		return client
	}
	return nil
}
