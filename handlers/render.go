package handlers

import (
	"github.com/torresposso/gosmic/views"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
)

func Render(c fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

// RenderLayout renders a component wrapped in the main layout
func RenderLayout(c fiber.Ctx, title string, pbClient interface{ IsAuthenticated() bool }, component templ.Component) error {
	csrfToken := csrf.TokenFromContext(c)

	flash, _ := c.Locals("flash").(string)
	flashType, _ := c.Locals("flash_type").(string)

	return Render(c, views.Layout(title, pbClient.IsAuthenticated(), csrfToken, flash, flashType, component))
}
