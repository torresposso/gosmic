# 03 - Navigation (Routing) 🗺️

**Mission Phase**: plotting Course  
**Objective**: Master the ship's navigation computer to handle URL routes.

## 📍 Coordinates (Basic Routes)

Navigate the ship to different sectors.

```go
// GET request: Retrieve sector info
app.Get("/sector/alpha", func(c fiber.Ctx) error {
    return c.SendString("Welcome to Sector Alpha")
})

// POST request: Send data to base
app.Post("/comms", func(c fiber.Ctx) error {
    return c.SendString("Message received")
})
```

## 🎯 Dynamic Coordinates (Parameters)

Capture variable segments like coordinates.

```go
// /system/:id (e.g., /system/sol)
app.Get("/system/:id", func(c fiber.Ctx) error {
    system := c.Params("id")
    return c.SendString("Approaching System: " + system)
})
```

## 🌌 Route Groups (Sectors)

Organize routes into logical sectors.

```go
// Public Sector
public := app.Group("/")
public.Get("/login", ShowLogin)

// Restricted Sector (Command)
command := app.Group("/command")
command.Get("/logs", ShowLogs)
```

## 🔄 Hyperjump (Redirection)

Redirect to another sector.

```go
app.Get("/jump", func(c fiber.Ctx) error {
    // Fiber v3 syntax
    return c.Redirect().To("/sector/beta")
})
```

---
[Next: 04 - Life Support (Middleware) →](./04-life-support.md)
