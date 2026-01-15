# 07 - Security Clearance (Auth System) 🔐

> **Rank**: Captain  
> **Skill**: AppSec, Tokens, and Cookies

**Mission**: Implement a secure authentication system that handles identity without exposing users to hijackers.

## 🛡️ The Security Matrix: Proxy Authentication

We implement a **Proxy Authentication** pattern. PocketBase acts as our Identity Provider (IdP), while our Fiber application manages the user session via secure cookies.

### The Identity Mechanism: PocketBase Tokens
When a user logs in, PocketBase issues a token (JWT). Instead of validating this token locally (which would require sharing secrets), we treat it as an opaque credential. We store it securely and pass it back to PocketBase on every request.

### The Transport Mechanism: HttpOnly Cookies
Where do we store this token?
*   **LocalStorage**: ❌ VULNERABLE to XSS. Malicious scripts can read this.
*   **HttpOnly Cookie**: ✅ SECURE. JavaScript cannot access this cookie; the browser sends it automatically.

## 🚧 Defense Against Dark Arts (CSRF)

Using cookies introduces a vulnerability: **Cross-Site Request Forgery (CSRF)**.
*   *Scenario*: A malicious site tricks a user's browser into sending a POST request to your app. Since the browser automatically attaches cookies, the request appears authenticated.
*   *Defense*: **CSRF Tokens**. The server generates a unique secret per session (`_csrf`) that must be included in every state-changing request (POST/PUT/DELETE). The malicious site cannot read this token, so its forged requests are rejected.

## 👨‍💻 Implementation (`handlers/auth.go`)

### 1. The Login Handler
We authenticate against PocketBase, receive a token, and plant it in a secure cookie.

```go
func (h *AuthHandler) Login() fiber.Handler {
    return func(c fiber.Ctx) error {
        // ... Retrieve credentials ...
        
        // 1. Authenticate with PocketBase via Service
        token, user, err := h.AuthService.Login(email, password)
        if err != nil {
             // ... handle error
        }

        // 2. Set the Secure Cookie
        c.Cookie(&fiber.Cookie{
            Name:     "pb_auth",
            Value:    token,
            HTTPOnly: true,  // JS cannot access
            SameSite: "Lax", // CSRF mitigation
            Path:     "/",
            Secure:   isProd, // SSL only in production
        })
        
        return c.Redirect().To("/dashboard")
    }
}
```

### 2. The Middleware (`middleware/auth.go`)
This gatekeeper ensures only officers with valid badges enter the command center.

```go
func AuthMiddleware(globalClient *pb.Client) fiber.Handler {
    return func(c fiber.Ctx) error {
        // 1. Retrieve Token from Cookie
        token := c.Cookies("pb_auth")
        if token == "" {
            return c.Redirect().To("/login")
        }
        
        // 2. Check Validity with PocketBase
        // This validates the token AND refreshes user data
        user, err := globalClient.WithToken(token).AuthRefresh()
        if err != nil {
            // Token expired or invalid
            return c.Redirect().To("/login")
        }
        
        // 3. Hydrate Context
        c.Locals("user", user)
        c.Locals("pb_client", globalClient.WithToken(token))
        
        return c.Next()
    }
}
```

## 🧠 Security Best Practice: Author Assignment
You might be tempted to send the `author` ID from the client when creating a post. **Don't.**
Trusting client input for ownership is dangerous (ID Spoofing). Instead, configure PocketBase **API Rules** to automatically assign the author:
`author = @request.auth.id`

This ensures that a user can only create records attributed to themselves.

---
[Next: 08 - Log Operations (Method Override) →](./08-log-operations.md)