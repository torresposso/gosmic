package handlers

import (
	"github.com/torresposso/gosmic/middleware"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/views"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

// Home renders the landing page (public route)
func Home(globalClient *pb.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Check if user has a token (for showing login/logout in nav)
		token := c.Cookies("pb_auth")
		userClient := globalClient.WithToken(token)
		return RenderLayout(c, "Fiber v3 + PocketBase Tutorial", userClient, views.Index(userClient.IsAuthenticated()))
	}
}

// Dashboard renders the user's dashboard (protected route)
func Dashboard() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Redirect().To("/login")
		}

		posts, err := client.ListPosts()
		if err != nil {
			// Return empty list instead of failing page load? Or handle error gracefully
			posts = []pb.Post{}
		}

		csrfToken := csrf.TokenFromContext(c)
		return RenderLayout(c, "Dashboard", client,
			views.Dashboard(client.GetCurrentUserName(), client.GetCurrentUserEmail(), len(posts), csrfToken))
	}
}
