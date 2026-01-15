package middleware

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/stretchr/testify/assert"
)

func TestFlashMiddleware(t *testing.T) {
	app := fiber.New()
	store := session.NewStore()

	app.Use(FlashMiddleware(store))

	app.Get("/set", func(c fiber.Ctx) error {
		sess, _ := store.Get(c)
		sess.Set("flash", "hello")
		sess.Set("flash_type", "success")
		sess.Save()
		return c.SendString("set")
	})

	app.Get("/get", func(c fiber.Ctx) error {
		flash := c.Locals("flash")
		flashType := c.Locals("flash_type")

		if flash == nil || flashType == nil {
			return c.SendString("none")
		}
		return c.SendString(flash.(string) + ":" + flashType.(string))
	})

	// 1. Set flash
	req := httptest.NewRequest("GET", "/set", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	setCookie := resp.Header.Get("Set-Cookie")
	cookieParts := strings.Split(setCookie, ";")
	cookie := cookieParts[0]

	// 2. Get flash (should be present in locals)
	req = httptest.NewRequest("GET", "/get", nil)
	req.Header.Set("Cookie", cookie)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "hello:success", string(respBody))

	// 3. Get flash again (should be cleared)
	req = httptest.NewRequest("GET", "/get", nil)
	req.Header.Set("Cookie", cookie)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	respBody, _ = io.ReadAll(resp.Body)
	assert.Equal(t, "none", string(respBody))
}
