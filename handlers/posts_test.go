package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/services"
)

func TestPostHandler_List(t *testing.T) {
	app := fiber.New()
	mockService := new(services.MockPostService)
	store := session.NewStore()
	handler := NewPostHandler(mockService, store)

	app.Get("/posts", func(c fiber.Ctx) error {
		// Mock inject PB client
		c.Locals("pb", &pb.Client{})
		return handler.List()(c)
	})

	t.Run("Success", func(t *testing.T) {
		posts := []pb.Post{{ID: "1", Title: "Test"}}
		mockService.On("List", mock.Anything, mock.Anything, "").Return(posts, nil).Once()

		req := httptest.NewRequest("GET", "/posts", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		appNoClient := fiber.New()
		appNoClient.Get("/posts", handler.List())

		req := httptest.NewRequest("GET", "/posts", nil)
		resp, err := appNoClient.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
	})
}

func TestPostHandler_Create(t *testing.T) {
	app := fiber.New()
	mockService := new(services.MockPostService)
	store := session.NewStore()
	handler := NewPostHandler(mockService, store)

	app.Post("/posts", func(c fiber.Ctx) error {
		c.Locals("pb", &pb.Client{})
		return handler.Create()(c)
	})

	t.Run("Success", func(t *testing.T) {
		mockService.On("Create", mock.Anything, mock.Anything, "New Post", "Content", true).Return(nil).Once()

		form := url.Values{}
		form.Add("title", "New Post")
		form.Add("content", "Content")
		form.Add("public", "on")

		req := httptest.NewRequest("POST", "/posts", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/posts", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestPostHandler_Delete(t *testing.T) {
	app := fiber.New()
	mockService := new(services.MockPostService)
	store := session.NewStore()
	handler := NewPostHandler(mockService, store)

	app.Delete("/posts/:id", func(c fiber.Ctx) error {
		c.Locals("pb", &pb.Client{})
		return handler.Delete()(c)
	})

	t.Run("StandardSuccess", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, mock.Anything, "1").Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/posts/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("HTMXSuccess", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, mock.Anything, "1").Return(nil).Once()

		req := httptest.NewRequest("DELETE", "/posts/1", nil)
		req.Header.Set("HX-Request", "true")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestPostHandler_Toggle(t *testing.T) {
	app := fiber.New()
	mockService := new(services.MockPostService)
	store := session.NewStore()
	handler := NewPostHandler(mockService, store)

	app.Post("/posts/:id/toggle", func(c fiber.Ctx) error {
		c.Locals("pb", &pb.Client{})
		return handler.Toggle()(c)
	})

	t.Run("HTMX_Success", func(t *testing.T) {
		mockService.On("TogglePublic", mock.Anything, mock.Anything, "1").Return(nil).Once()
		mockService.On("Get", mock.Anything, mock.Anything, "1").Return(&pb.Post{ID: "1", Title: "T"}, nil).Once()

		req := httptest.NewRequest("POST", "/posts/1/toggle", nil)
		req.Header.Set("HX-Request", "true")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("HTMX_Error_Service", func(t *testing.T) {
		mockService.On("TogglePublic", mock.Anything, mock.Anything, "1").Return(assert.AnError).Once()

		req := httptest.NewRequest("POST", "/posts/1/toggle", nil)
		req.Header.Set("HX-Request", "true")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode) // Views.FlashMessage returns 200 with OOB
		mockService.AssertExpectations(t)
	})
}

func TestPostHandler_GetAndEditErrors(t *testing.T) {
	app := fiber.New()
	mockService := new(services.MockPostService)
	store := session.NewStore()
	handler := NewPostHandler(mockService, store)

	app.Get("/posts/:id", func(c fiber.Ctx) error {
		c.Locals("pb", &pb.Client{})
		return handler.Get()(c)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.On("Get", mock.Anything, mock.Anything, "999").Return(nil, assert.AnError).Once()
		req := httptest.NewRequest("GET", "/posts/999", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestPostHandler_GetAndEdit(t *testing.T) {
	app := fiber.New()
	mockService := new(services.MockPostService)
	store := session.NewStore()
	handler := NewPostHandler(mockService, store)

	app.Get("/posts/:id", func(c fiber.Ctx) error {
		c.Locals("pb", &pb.Client{})
		return handler.Get()(c)
	})
	app.Get("/posts/:id/edit", func(c fiber.Ctx) error {
		c.Locals("pb", &pb.Client{})
		return handler.Edit()(c)
	})

	t.Run("GetSuccess", func(t *testing.T) {
		mockService.On("Get", mock.Anything, mock.Anything, "1").Return(&pb.Post{ID: "1", Title: "T"}, nil).Once()
		req := httptest.NewRequest("GET", "/posts/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("EditSuccess", func(t *testing.T) {
		mockService.On("Get", mock.Anything, mock.Anything, "1").Return(&pb.Post{ID: "1", Title: "T"}, nil).Once()
		req := httptest.NewRequest("GET", "/posts/1/edit", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestPostHandler_Update(t *testing.T) {
	app := fiber.New()
	mockService := new(services.MockPostService)
	store := session.NewStore()
	handler := NewPostHandler(mockService, store)

	app.Post("/posts/:id/update", func(c fiber.Ctx) error {
		c.Locals("pb", &pb.Client{})
		return handler.Update()(c)
	})

	t.Run("Success", func(t *testing.T) {
		mockService.On("Update", mock.Anything, mock.Anything, "1", "Updated", "Content", false).Return(nil).Once()

		form := url.Values{}
		form.Add("title", "Updated")
		form.Add("content", "Content")

		req := httptest.NewRequest("POST", "/posts/1/update", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
