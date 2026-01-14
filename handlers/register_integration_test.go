package handlers

import (
	"bytes"
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

func TestRegisterFlow(t *testing.T) {
	// Setup Mock Client
	mockTripper := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) *http.Response {
			if strings.Contains(req.URL.Path, "/api/collections/users/records") {
				// Simulate successful creation
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("{}")),
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

	app.Get("/register", authHandler.ShowRegister())
	app.Post("/register", authHandler.Register())

	// Step 1: GET /register to get CSRF token
	req1 := httptest.NewRequest("GET", "/register", nil)
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

	re := regexp.MustCompile(`name="_csrf" value="([^"]+)"`)
	matches := re.FindStringSubmatch(bodyStr)
	assert.True(t, len(matches) > 1, "CSRF token input not found in HTML")
	csrfToken := matches[1]

	// Step 2: POST /register with VALID data
	form := url.Values{}
	form.Add("email", "newuser@example.com")
	form.Add("password", "password123")
	form.Add("passwordConfirm", "password123")
	form.Add("name", "New User")
	form.Add("_csrf", csrfToken)

	req2 := httptest.NewRequest("POST", "/register", strings.NewReader(form.Encode()))
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req2.AddCookie(csrfCookie)

	resp2, err := app.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, resp2.StatusCode, "Register failed or CSRF rejected")
	assert.Equal(t, "/login?registered=true", resp2.Header.Get("Location"))

	// Step 3: POST /register with INVALID data (mismatched passwords)
	formInvalid := url.Values{}
	formInvalid.Add("email", "fail@example.com")
	formInvalid.Add("password", "password123")
	formInvalid.Add("passwordConfirm", "wrongpass")
	formInvalid.Add("name", "Fail User")
	formInvalid.Add("_csrf", csrfToken)

	req3 := httptest.NewRequest("POST", "/register", strings.NewReader(formInvalid.Encode()))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req3.AddCookie(csrfCookie)

	resp3, err := app.Test(req3)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp3.StatusCode) // Should render form again with error

	bodyBytes3, _ := io.ReadAll(resp3.Body)
	assert.Contains(t, string(bodyBytes3), "Passwords do not match")
}
