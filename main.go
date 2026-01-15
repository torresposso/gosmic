package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/torresposso/gosmic/handlers"
	"github.com/torresposso/gosmic/middleware"
	"github.com/torresposso/gosmic/pb"
	"github.com/torresposso/gosmic/repositories"
	"github.com/torresposso/gosmic/services"
)

func main() {
	pbURL := getEnv("PB_URL", "http://localhost:8090")
	port := getEnv("PORT", "8080")
	baseURL := getEnv("BASE_URL", "http://localhost:"+port)
	isProd := os.Getenv("GO_ENV") == "production"

	// Create a global PocketBase client (connection pool is shared)
	globalClient := pb.NewClient(pbURL)

	app := fiber.New(fiber.Config{
		AppName:       "Fiber v3 + PocketBase Tutorial",
		StrictRouting: true,
		CaseSensitive: true,
	})

	sessStore := session.NewStore()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	// Static files
	app.Use("/static", static.New("./static", static.Config{
		CacheDuration: 3600 * time.Second,
		Compress:      true,
	}))

	// CSRF Protection
	app.Use(csrf.New(csrf.Config{
		Extractor:      extractors.FromForm("_csrf"),
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		CookieSecure:   isProd,
		CookieHTTPOnly: true,
		IdleTimeout:    1 * time.Hour,
	}))

	// Support _method input for PUT/DELETE in forms
	app.Use(middleware.MethodOverride())
	app.Use(middleware.FlashMiddleware(sessStore))

	// Initialize Repositories
	postRepo := repositories.NewPostRepository()
	authRepo := repositories.NewAuthRepository()

	// Initialize Services
	postService := services.NewPostService(postRepo)
	authService := services.NewAuthService(authRepo)

	// Initialize Handlers
	postHandler := handlers.NewPostHandler(postService, sessStore)
	authHandler := handlers.NewAuthHandler(authService, globalClient)
	rootHandler := handlers.NewRootHandler(globalClient, postService)

	// Public routes
	app.Get("/", rootHandler.Home())
	app.Get("/login", authHandler.ShowLogin())
	app.Post("/login", authHandler.Login())
	app.Get("/register", authHandler.ShowRegister())
	app.Post("/register", authHandler.Register())
	app.Get("/logout", authHandler.Logout())

	// Protected routes - middleware creates request-scoped client
	protected := app.Group("/dashboard", middleware.AuthMiddleware(globalClient))
	protected.Get("", rootHandler.Dashboard())
	protected.Get("/posts", postHandler.List())
	protected.Post("/posts", postHandler.Create())
	protected.Get("/posts/:id", postHandler.Get())
	protected.Get("/posts/:id/edit", postHandler.Edit())
	protected.Put("/posts/:id", postHandler.Update())
	protected.Delete("/posts/:id", postHandler.Delete())

	// API routes
	api := app.Group("/api", middleware.AuthMiddleware(globalClient))
	api.Post("/posts/:id/toggle", postHandler.Toggle())

	log.Printf("Server starting on %s", baseURL)
	log.Printf("PocketBase: %s", pbURL)
	log.Fatal(app.Listen(":" + port))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
