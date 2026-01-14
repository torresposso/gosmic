package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// FlashMiddleware extracts flash messages from the session and puts them into context locals.
// It also clears them from the session so they are only displayed once.
func FlashMiddleware(store *session.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			// If we can't get the session, just continue
			return c.Next()
		}

		flash := sess.Get("flash")
		flashType := sess.Get("flash_type")

		if flash != nil {
			c.Locals("flash", flash.(string))
			sess.Delete("flash")
		}
		if flashType != nil {
			c.Locals("flash_type", flashType.(string))
			sess.Delete("flash_type")
		}

		// Save only if we modified the session
		if flash != nil || flashType != nil {
			if err := sess.Save(); err != nil {
				return c.Next()
			}
		}

		return c.Next()
	}
}
