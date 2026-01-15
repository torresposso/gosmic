package handlers

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/services"
	"github.com/torresposso/gosmic/views"
)

type DocHandler struct {
	Service      *services.DocService
	GlobalClient *pb.Client
}

func NewDocHandler(service *services.DocService, client *pb.Client) *DocHandler {
	return &DocHandler{
		Service:      service,
		GlobalClient: client,
	}
}

func (h *DocHandler) Index() fiber.Handler {
	return func(c fiber.Ctx) error {
		log.Println("DEBUG: DocHandler.Index called")
		chapters, err := h.Service.GetChapters()
		if err != nil || len(chapters) == 0 {
			log.Printf("DEBUG: No chapters found (Err: %v, Len: %d)", err, len(chapters))
			return c.Status(fiber.StatusNotFound).SendString("No documentation found")
		}
		// Redirect to the first chapter
		return c.Redirect().To("/docs/" + chapters[0].ID)
	}
}

func (h *DocHandler) Show() fiber.Handler {
	return func(c fiber.Ctx) error {
		chapterID := c.Params("chapter")

		chapters, err := h.Service.GetChapters()
		if err != nil {
			log.Printf("Error getting chapters: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		content, err := h.Service.GetChapterContent(chapterID)
		if err != nil {
			log.Printf("Apparent 404 for chapter %s: %v", chapterID, err)
			return c.Status(fiber.StatusNotFound).SendString("Chapter not found")
		}

		// Create request-scoped client to check auth status
		token := c.Cookies("pb_auth")
		userClient := h.GlobalClient.WithToken(token)

		return RenderLayout(c, "Documentation", userClient, views.Docs(chapters, chapterID, content))
	}
}
