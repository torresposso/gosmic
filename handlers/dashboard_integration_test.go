package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/stretchr/testify/assert"
	"github.com/torresposso/gosmic/pb"
)

func TestDashboardAccess(t *testing.T) {
	// Setup Mock Client
	mockTripper := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.Path, "/api/collections/posts/records") {
				respBody := map[string]interface{}{
					"items": []map[string]interface{}{
						{"id": "post1", "title": "Test Post", "content": "Content"},
					},
				}
				body, _ := json.Marshal(respBody)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				}
			}
			return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(bytes.NewBufferString(""))}
		},
	}

	pbClient := pb.NewClient("http://mock-pb")
	pbClient.HTTPClient.Transport = mockTripper

	app := fiber.New()

	app.Use(csrf.New(csrf.Config{
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		IdleTimeout:    1 * time.Hour,
	}))

	// Mimic Auth Middleware logic manually for test simplicity
	app.Use(func(c fiber.Ctx) error {
		// Mock authenticated user
		userClient := pbClient.WithToken("mock-token")
		userClient.AuthRecord = &pb.User{
			ID:    "user123",
			Email: "test@example.com",
			Name:  "Test User",
		}
		c.Locals("pb", userClient)
		return c.Next()
	})

	app.Get("/dashboard", Dashboard())

	req := httptest.NewRequest("GET", "/dashboard", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "Welcome aboard")
	assert.Contains(t, bodyStr, "Commander Test User")
	assert.Contains(t, bodyStr, "test@example.com")
	// assert.Contains(t, bodyStr, "1") // Post count check, brittle if format changes
}
