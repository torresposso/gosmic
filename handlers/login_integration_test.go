package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/stretchr/testify/assert"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/repositories"
	"github.com/torresposso/gosmic/services"
)

func TestLoginFlowWithCSRF(t *testing.T) {
	// Setup Mock Client
	mockTripper := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.Path, "auth-with-password") {
				respBody := map[string]interface{}{
					"token": "valid-pb-token",
					"record": map[string]interface{}{
						"id":    "user123",
						"email": "test@example.com",
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

	// Add CSRF Middleware similar to main.go
	app.Use(csrf.New(csrf.Config{
		Extractor:      extractors.FromForm("_csrf"),
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		IdleTimeout:    1 * time.Hour,
	}))

	// Initialize Services and Handlers for testing
	authRepo := repositories.NewAuthRepository()
	authService := services.NewAuthService(authRepo)
	authHandler := NewAuthHandler(authService, pbClient)

	app.Get("/login", authHandler.ShowLogin())
	app.Post("/login", authHandler.Login())

	// Step 1: GET /login to get CSRF token
	req1 := httptest.NewRequest("GET", "/login", nil)
	resp1, err := app.Test(req1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	// Extract CSRF cookie
	cookies := resp1.Cookies()
	var csrfCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "csrf_" {
			csrfCookie = c
			break
		}
	}
	assert.NotNil(t, csrfCookie, "CSRF cookie not found")

	// Extract CSRF token from HTML body
	bodyBytes, _ := io.ReadAll(resp1.Body)
	bodyStr := string(bodyBytes)

	// Simple regex to find the token value
	re := regexp.MustCompile(`name="_csrf" value="([^"]+)"`)
	matches := re.FindStringSubmatch(bodyStr)
	assert.True(t, len(matches) > 1, "CSRF token input not found in HTML")
	csrfToken := matches[1]

	// Step 2: POST /login with CSRF token
	form := url.Values{}
	form.Add("email", "test@example.com")
	form.Add("password", "password123")
	form.Add("_csrf", csrfToken)

	req2 := httptest.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req2.AddCookie(csrfCookie) // Important: Send back the cookie

	resp2, err := app.Test(req2)
	assert.NoError(t, err)

	// Debug output if it fails
	if resp2.StatusCode != http.StatusSeeOther {
		bodyBytes2, _ := io.ReadAll(resp2.Body)
		t.Logf("Response Body: %s", string(bodyBytes2))
	}

	assert.Equal(t, http.StatusSeeOther, resp2.StatusCode, "Login failed or CSRF rejected")
}
