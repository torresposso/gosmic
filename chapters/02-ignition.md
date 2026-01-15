# 02 - Ignition 🚀

> **Rank**: Ensign  
> **Skill**: Handlers, Context, and Middleware

**Mission**: Launch the web server and understand the "Why" behind Fiber's speed.

## ⚡ Technical Deep Dive: Fiber v3 vs. The World

Standard Go web frameworks (like `net/http`) are robust but can be verbose. **Fiber** is an Express-inspired web framework built on top of **`fasthttp`**, the fastest HTTP engine for Go.

### The Zero-Allocation Philosophy
In standard Go `net/http`, every request allocates significant memory for headers, URLs, and bodies. This creates "Garbage" that the "Garbage Collector" (GC) must clean up, leading to potential latency spikes.

Fiber acts like a specialized racing engine:
1.  **Memory Reuse**: It aggressively reuses memory for context and request handling.
2.  **Performance**: It is optimized for high throughput and low memory footprint.

## 🔌 The Main Thrusters (`main.go`)

Let's examine the production-ready entry point.

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	// ... imports
)

func main() {
	// 1. Load Configurations
	port := getEnv("PORT", "8080")

	// 2. Initialize the App
	app := fiber.New(fiber.Config{
		AppName:       "Gosmic Code",
		StrictRouting: true,
		CaseSensitive: true,
	})

	// 3. Middleware Chain (The Shield Generator)
	app.Use(recover.New()) 
	app.Use(logger.New())
	app.Use(compress.New()) // Gzip/Brotli

	// 🛡️ CSRF Protection
	app.Use(csrf.New(csrf.Config{
		// ... config
	}))

	// 📂 Static Assets with Cache Control
	app.Use("/static", static.New("./static", static.Config{
		Compress: true,
		ModifyResponse: func(c fiber.Ctx) error {
			c.Set("Cache-Control", "public, max-age=31536000")
			return nil
		},
	}))

	// 4. Launch (Hyperspace Jump)
	log.Fatal(app.Listen(":" + port))
}
```

## 🧱 The View System: Type-Safe Rendering

We avoid string concatenation for HTML to prevent XSS attacks and logic errors. We use **Templ**.

### Why Templ?
*   **Compilation**: Errors are caught at build time, not runtime.
*   **Security**: Context-aware escaping (auto-sanitizes inputs).
*   **Integration**: Components are just Go functions, making them easy to test and compose.

`views/hello.templ`:
```go
package views

templ Hello(name string) {
	<div>Welcome aboard, Commander { name }!</div>
}
```

When you run `templ generate`, this becomes efficient Go code that writes directly to the response buffer.

---
[Next: 03 - Navigation (Routing) →](./03-navigation.md)