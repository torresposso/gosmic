# 05 - Cargo Hold (PocketBase) 📦

**Mission Phase**: Supply Chain  
**Objective**: Establish a secure link with the ship's cargo hold (Database).

## 📡 The Uplink (Custom Client)

While official SDKs exist, building a custom typed client allows for greater control and reduced dependencies. We implement a lightweight wrapper around `net/http` to communicate with PocketBase.

### 1. The Client Structure
Our client manages the connection pool and the authentication state (token).

```go
type Client struct {
    BaseURL    string
    HTTPClient *http.Client
    AuthToken  string // The JWT key to the cargo hold
}
```

### 2. Authentication (Handshake)
To open the hold, we authenticate with credentials. This returns the raw token and user record.

```go
func (c *Client) AuthWithPassword(email, password string) (string, *User, error) {
    // ... Send POST to /api/collections/users/auth-with-password
    // ... Returns token and user struct
}
```

### 3. Request Scoping (Critical Architecture)
We use a **Request Scoping** pattern. The `globalClient` holds the connection pool, but for each request, we create a lightweight clone that carries the specific user's token. This prevents race conditions where one user's token could overwrite another's in a shared client.

```go
// WithToken creates a shallow copy of the client with the user's token.
// It shares the underlying HTTP connection pool for efficiency.
func (c *Client) WithToken(token string) *Client {
    return &Client{
        BaseURL:    c.BaseURL,
        HTTPClient: c.HTTPClient, 
        AuthToken:  token, // Request-specific token
    }
}
```

### 4. Typed Methods
We create strongly-typed methods for our data models.

```go
func (c *Client) CreatePost(post *Post) error {
    // ... Send POST to /api/collections/posts/records
    // ... Attach Authorization header using c.AuthToken
}
```

This manual control ensures our application remains lightweight and type-safe.

---
[Next: 06 - Hyperdrive UI (Alpine.js) →](./06-hyperdrive-ui.md)