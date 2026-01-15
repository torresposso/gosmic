package handlers

import (
	"github.com/torresposso/gosmic/middleware"
	"github.com/torresposso/gosmic/services"
	"github.com/torresposso/gosmic/views"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

type PostHandler struct {
	postService services.PostService
	sessStore   *session.Store
}

func NewPostHandler(ps services.PostService, ss *session.Store) *PostHandler {
	return &PostHandler{postService: ps, sessStore: ss}
}

func (h *PostHandler) List() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Redirect().To("/login")
		}

		query := c.Query("q")
		posts, err := h.postService.List(c.Context(), client, query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to load posts")
		}

		csrfToken := csrf.TokenFromContext(c)
		return RenderLayout(c, "Posts", client, views.Posts(posts, csrfToken))
	}
}

func (h *PostHandler) Create() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Redirect().To("/login")
		}

		title := c.FormValue("title")
		content := c.FormValue("content")
		isPublic := c.FormValue("public") == "on"

		if title == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Title is required")
		}

		err := h.postService.Create(c.Context(), client, title, content, isPublic)
		sess, _ := h.sessStore.Get(c)
		if err != nil {
			sess.Set("flash", "Failed to create post")
			sess.Set("flash_type", "error")
			sess.Save()
			return c.Redirect().To("/dashboard/posts")
		}

		sess.Set("flash", "Mission log recorded successfully")
		sess.Set("flash_type", "success")
		sess.Save()
		return c.Redirect().To("/dashboard/posts")
	}
}

func (h *PostHandler) Get() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Redirect().To("/login")
		}

		post, err := h.postService.Get(c.Context(), client, c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Post not found")
		}
		csrfToken := csrf.TokenFromContext(c)
		return Render(c, views.PostItem(*post, csrfToken))
	}
}

func (h *PostHandler) Edit() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Redirect().To("/login")
		}

		post, err := h.postService.Get(c.Context(), client, c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Post not found")
		}
		csrfToken := csrf.TokenFromContext(c)
		return RenderLayout(c, "Edit Log", client, views.EditPostForm(*post, csrfToken))
	}
}

// Update updates an existing post
func (h *PostHandler) Update() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Redirect().To("/login")
		}

		id := c.Params("id")
		title := c.FormValue("title")
		content := c.FormValue("content")
		isPublic := c.FormValue("public") == "on"

		err := h.postService.Update(c.Context(), client, id, title, content, isPublic)
		sess, _ := h.sessStore.Get(c)
		if err != nil {
			sess.Set("flash", "Failed to update log")
			sess.Set("flash_type", "error")
			sess.Save()
			return c.Redirect().To("/dashboard/posts")
		}

		sess.Set("flash", "Mission log updated successfully")
		sess.Set("flash_type", "success")
		sess.Save()
		return c.Redirect().To("/dashboard/posts")
	}
}

func (h *PostHandler) Delete() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Redirect().To("/login")
		}

		postID := c.Params("id")
		err := h.postService.Delete(c.Context(), client, postID)

		if c.Get("HX-Request") == "true" {
			c.Set("Content-Type", "text/html")
			if err != nil {
				return views.FlashMessage("Failed to purge log", "error").Render(c.Context(), c.Response().BodyWriter())
			}
			// When deleting with hx-target="#post-ID", returning only the OOB flash
			// effectively clears the target element.
			return views.FlashMessage("Mission log purged successfully", "success").Render(c.Context(), c.Response().BodyWriter())
		}

		sess, _ := h.sessStore.Get(c)
		if err != nil {
			sess.Set("flash", "Failed to purge log")
			sess.Set("flash_type", "error")
			sess.Save()
			return c.Redirect().To("/dashboard/posts")
		}

		sess.Set("flash", "Mission log purged successfully")
		sess.Set("flash_type", "success")
		sess.Save()
		return c.Redirect().To("/dashboard/posts")
	}
}

func (h *PostHandler) Toggle() fiber.Handler {
	return func(c fiber.Ctx) error {
		client := middleware.GetPBClient(c)
		if client == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		id := c.Params("id")
		err := h.postService.TogglePublic(c.Context(), client, id)

		if c.Get("HX-Request") == "true" {
			c.Set("Content-Type", "text/html")
			if err != nil {
				return views.FlashMessage("Failed to toggle visibility", "error").Render(c.Context(), c.Response().BodyWriter())
			}

			post, err := h.postService.Get(c.Context(), client, id)
			if err != nil {
				return views.FlashMessage("Log disappeared during transmission", "error").Render(c.Context(), c.Response().BodyWriter())
			}

			csrfToken := csrf.TokenFromContext(c)
			views.PostItem(*post, csrfToken).Render(c.Context(), c.Response().BodyWriter())
			return views.FlashMessage("Visibility matrix updated", "success").Render(c.Context(), c.Response().BodyWriter())
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to toggle visibility"})
		}

		return c.Redirect().To("/dashboard/posts")
	}
}
