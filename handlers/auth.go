package handlers

import (
	"github.com/torresposso/gosmic/middleware"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/services"
	"github.com/torresposso/gosmic/views"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

type AuthHandler struct {
	authService  services.AuthService
	globalClient *pb.Client
}

func NewAuthHandler(as services.AuthService, gc *pb.Client) *AuthHandler {
	return &AuthHandler{
		authService:  as,
		globalClient: gc,
	}
}

// ShowLogin renders the login form
func (h *AuthHandler) ShowLogin() fiber.Handler {
	return func(c fiber.Ctx) error {
		registered := c.Query("registered") == "true"
		csrfToken := csrf.TokenFromContext(c)
		return RenderLayout(c, "Login", h.globalClient.WithToken(""), views.Login("", "", registered, csrfToken))
	}
}

// Login authenticates a user and sets the PocketBase token in a cookie
func (h *AuthHandler) Login() fiber.Handler {
	return func(c fiber.Ctx) error {
		email := c.FormValue("email")
		password := c.FormValue("password")
		csrfToken := csrf.TokenFromContext(c)

		token, err := h.authService.Login(c.Context(), h.globalClient, email, password)
		if err != nil {
			return RenderLayout(c, "Login", h.globalClient.WithToken(""), views.Login(err.Error(), email, false, csrfToken))
		}

		c.Cookie(&fiber.Cookie{
			Name:     "pb_auth",
			Value:    token,
			HTTPOnly: true,
			Secure:   c.Protocol() == "https", // Automatically detect
			SameSite: "Lax",
			Path:     "/",
		})

		return c.Redirect().To("/dashboard")
	}
}

// ShowRegister renders the registration form
func (h *AuthHandler) ShowRegister() fiber.Handler {
	return func(c fiber.Ctx) error {
		csrfToken := csrf.TokenFromContext(c)
		return RenderLayout(c, "Register", h.globalClient.WithToken(""), views.Register("", "", "", csrfToken))
	}
}

// Register creates a new user account
func (h *AuthHandler) Register() fiber.Handler {
	return func(c fiber.Ctx) error {
		email := c.FormValue("email")
		password := c.FormValue("password")
		passwordConfirm := c.FormValue("passwordConfirm")
		name := c.FormValue("name")
		csrfToken := csrf.TokenFromContext(c)

		if password != passwordConfirm {
			return RenderLayout(c, "Register", h.globalClient.WithToken(""), views.Register("Passwords do not match", email, name, csrfToken))
		}

		err := h.authService.Register(c.Context(), h.globalClient, email, password, name)
		if err != nil {
			return RenderLayout(c, "Register", h.globalClient.WithToken(""), views.Register(err.Error(), email, name, csrfToken))
		}

		return c.Redirect().To("/login?registered=true")
	}
}

// Logout clears the authentication cookie
func (h *AuthHandler) Logout() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.ClearCookie("pb_auth")
		return c.Redirect().To("/")
	}
}

type RootHandler struct {
	globalClient *pb.Client
	postService  services.PostService
}

func NewRootHandler(gc *pb.Client, ps services.PostService) *RootHandler {
	return &RootHandler{
		globalClient: gc,
		postService:  ps,
	}
}

func (h *RootHandler) Home() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("pb_auth")
		userClient := h.globalClient.WithToken(token)
		return RenderLayout(c, "Fiber v3 + PocketBase Tutorial", userClient, views.Index(userClient.IsAuthenticated()))
	}
}

func (h *RootHandler) Dashboard() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Redirect().To("/login")
		}

		posts, err := h.postService.List(c.Context(), client, "")
		if err != nil {
			posts = []pb.Post{}
		}

		csrfToken := csrf.TokenFromContext(c)
		return RenderLayout(c, "Dashboard", client,
			views.Dashboard(client.GetCurrentUserName(), client.GetCurrentUserEmail(), len(posts), csrfToken))
	}
}
