package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/torresposso/gosmic/middleware"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/repositories"
	"github.com/torresposso/gosmic/services"
)

// MockRoundTripper allows us to mock HTTP responses
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req), nil
}

func TestLoginHandler(t *testing.T) {
	// Setup Mock Client
	mockTripper := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) *http.Response {
			if req.URL.Path == "/api/collections/users/auth-with-password" {
				respBody := map[string]interface{}{
					"token": "valid-pb-token",
					"record": map[string]interface{}{
						"id":    "user123",
						"email": "test@example.com",
						"name":  "Test User",
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
	authRepo := repositories.NewAuthRepository()
	authService := services.NewAuthService(authRepo)
	authHandler := NewAuthHandler(authService, pbClient)
	app.Post("/login", authHandler.Login())

	// Test Case: Successful Login
	form := url.Values{}
	form.Add("email", "test@example.com")
	form.Add("password", "password123")
	req := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode) // Redirects to /dashboard

	// Verify Cookie - now uses pb_auth instead of auth_token
	cookies := resp.Cookies()
	var authCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "pb_auth" {
			authCookie = c
			break
		}
	}
	assert.NotNil(t, authCookie)
	assert.Equal(t, "pb_auth", authCookie.Name)
	assert.Equal(t, "valid-pb-token", authCookie.Value) // Direct PB token
}

func TestAuthMiddleware(t *testing.T) {
	pbClient := pb.NewClient("http://mock-pb")

	app := fiber.New()
	app.Use(middleware.AuthMiddleware(pbClient))
	app.Get("/dashboard", func(c fiber.Ctx) error {
		return c.SendString("Protected Content")
	})

	// Test Case: Valid Cookie with PB token
	req := httptest.NewRequest("GET", "/dashboard", nil)
	req.AddCookie(&http.Cookie{Name: "pb_auth", Value: "valid-pb-token"})
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test Case: No Cookie
	req = httptest.NewRequest("GET", "/dashboard", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode) // Redirects to /login
}

func TestLogoutHandler(t *testing.T) {
	app := fiber.New()
	pbClient := pb.NewClient("http://mock-pb")
	authRepo := repositories.NewAuthRepository()
	authService := services.NewAuthService(authRepo)
	authHandler := NewAuthHandler(authService, pbClient)
	app.Get("/logout", authHandler.Logout())

	req := httptest.NewRequest("GET", "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "pb_auth", Value: "some-token"})

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, resp.StatusCode) // Redirects to /

	// Verify cookie is cleared
	cookies := resp.Cookies()
	for _, c := range cookies {
		if c.Name == "pb_auth" {
			// Cookie should be expired/cleared
			assert.True(t, c.MaxAge < 0 || c.Value == "")
		}
	}
}
