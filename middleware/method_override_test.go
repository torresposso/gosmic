package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestMethodOverride(t *testing.T) {
	app := fiber.New()
	app.Use(MethodOverride())

	app.All("/", func(c fiber.Ctx) error {
		return c.SendString(c.Method())
	})

	t.Run("IgnoreNonPost", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, "/?_method=DELETE", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Should still be GET
		req2 := httptest.NewRequest(fiber.MethodGet, "/", nil)
		resp2, _ := app.Test(req2)
		assert.Equal(t, fiber.MethodGet, getResponseBody(resp2))
	})

	t.Run("OverrideViaFormValue", func(t *testing.T) {
		data := url.Values{}
		data.Set("_method", "DELETE")
		req := httptest.NewRequest(fiber.MethodPost, "/", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, _ := app.Test(req)
		assert.Equal(t, "DELETE", getResponseBody(resp))
	})

	t.Run("OverrideViaHeader", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/", nil)
		req.Header.Set("X-HTTP-Method-Override", "PUT")

		resp, _ := app.Test(req)
		assert.Equal(t, "PUT", getResponseBody(resp))
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/?_method=INVALID", nil)
		resp, _ := app.Test(req)
		// Should remain POST
		assert.Equal(t, "POST", getResponseBody(resp))
	})
}

func getResponseBody(resp *http.Response) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	return buf.String()
}
