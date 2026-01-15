# 04 - Life Support (Middleware) 🛡️

**Mission Phase**: Systems Check  
**Objective**: Intercept signals for security, logging, and performance.

## 🛡️ Systems Pipeline

Middleware sits between the user and the core handlers.

```text
Signal → [Compress] → [Logger] → [Security] → [Handlers]
```

## 🔧 Essential Systems

### 1. Static Assets (Hull Plating)
Serve CSS, JS, and images with advanced caching policies.

```go
app.Use("/static", static.New("./static", static.Config{
    Compress: true,
    ModifyResponse: func(c fiber.Ctx) error {
        c.Set("Cache-Control", "public, max-age=31536000")
        return nil
    },
}))
```

### 2. Logger (Flight Recorder)
Record every interaction.
```go
app.Use(logger.New(logger.Config{
    Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
    TimeFormat: "2006-01-02 15:04:05",
}))
```

## 🔐 Security Field (Custom Auth)

We use a custom middleware to protect the Command Center. It extracts the authentication cookie, validates it against PocketBase, and injects a scoped client into the request context.

```go
// middleware/auth.go
func AuthMiddleware(globalClient *pb.Client) fiber.Handler {
    return func(c fiber.Ctx) error {
        token := c.Cookies("pb_auth")
        if token == "" {
            return c.Redirect().To("/login")
        }

        // Validate token with PocketBase (creates a request-scoped client)
        user, err := globalClient.WithToken(token).AuthRefresh()
        if err != nil {
            return c.Redirect().To("/login")
        }

        // Store confirmed user and scoped client in locals
        c.Locals("user", user)
        c.Locals("pb_client", globalClient.WithToken(token))

        return c.Next()
    }
}
```

This ensures that `protected` routes always have access to a valid, authenticated PocketBase client.

---
[Next: 05 - Cargo Hold (PocketBase) →](./05-cargo-hold.md)
