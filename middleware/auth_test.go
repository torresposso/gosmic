package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/torresposso/gosmic/pb"
)

func TestAuthMiddleware(t *testing.T) {
	globalClient := pb.NewClient("http://localhost:8090")
	app := fiber.New()
	app.Use(AuthMiddleware(globalClient))

	app.Get("/protected", func(c fiber.Ctx) error {
		client := GetPBClient(c)
		if client == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.JSON(fiber.Map{
			"id":    client.GetUserID(),
			"token": client.AuthToken,
		})
	})

	t.Run("RedirectIfNoCookie", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusSeeOther, resp.StatusCode)
		assert.Equal(t, "/login", resp.Header.Get("Location"))
	})

	t.Run("AuthenticatedWithValidToken", func(t *testing.T) {
		userID := "user-123"
		claims := map[string]string{"id": userID}
		claimsJSON, _ := json.Marshal(claims)
		payload := base64.RawURLEncoding.EncodeToString(claimsJSON)
		token := "header." + payload + ".signature"

		req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "pb_auth",
			Value: token,
		})

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var resBody map[string]string
		json.NewDecoder(resp.Body).Decode(&resBody)
		assert.Equal(t, userID, resBody["id"])
		assert.Equal(t, token, resBody["token"])
	})

	t.Run("InvalidTokenPayloadDoesNotSetID", func(t *testing.T) {
		token := "header.invalidpayload.signature"
		req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "pb_auth",
			Value: token,
		})

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var resBody map[string]string
		json.NewDecoder(resp.Body).Decode(&resBody)
		assert.Equal(t, "", resBody["id"])
		assert.Equal(t, token, resBody["token"])
	})
}
