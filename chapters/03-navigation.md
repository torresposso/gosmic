# 03 - Navigation (Routing) 🗺️

**Mission Phase**: Plotting Course  
**Objective**: Master the ship's navigation computer to handle URL routes using Fiber's router and Middleware groups.

## 📍 Coordinates (Basic Routes)

Navigate the ship to different sectors.

```go
// GET request: Retrieve sector info
app.Get("/", rootHandler.Home())

// POST request: Send data to base
app.Post("/login", authHandler.Login())
```

## 🎯 Dynamic Coordinates (Parameters)

Capture variable segments like coordinates.

```go
// /posts/:id (e.g., /posts/123)
app.Get("/posts/:id", postHandler.Get())
```

In the handler:
```go
func (h *PostHandler) Get() fiber.Handler {
    return func(c fiber.Ctx) error {
        id := c.Params("id")
        return c.SendString("Scanning Log: " + id)
    }
}
```

## 🌌 Route Groups (Sectors)

Organize routes into logical sectors, especially to apply security protocols (Middleware).

```go
// 1. Initialize Handlers (Dependency Injection)
authHandler := handlers.NewAuthHandler(authService, globalClient)
postHandler := handlers.NewPostHandler(postService, sessStore)

// 2. Public Sector
app.Get("/", rootHandler.Home())
app.Get("/login", authHandler.ShowLogin())

// 3. Protected Sector (Requires Authentication)
// All routes in this group will pass through AuthMiddleware first
protected := app.Group("/dashboard", middleware.AuthMiddleware(globalClient))

protected.Get("", rootHandler.Dashboard())
protected.Get("/posts", postHandler.List())
protected.Post("/posts", postHandler.Create())
```

## 🔄 Hyperjump (Redirection)

Redirect to another sector.

```go
// Direct redirect
return c.Redirect().To("/dashboard")

// With Flash Message (via Session)
utils.SetFlash(c, "success", "Hyperjump successful!")
return c.Redirect().To("/system/beta")
```

---
[Next: 04 - Life Support (Middleware) →](./04-life-support.md)
