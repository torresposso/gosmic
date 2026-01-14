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
    AuthRecord *User  // Cached user details
}
```

### 2. Authentication (Handshake)
To open the hold, we authenticate with credentials. This method returns the raw token and the user record, allowing the handler to manage the session (e.g., via cookies).

```go
func (c *Client) AuthWithPassword(email, password string) (string, *User, error) {
    body := map[string]interface{}{
        "identity": email,
        "password": password,
    }
    
    // ... Send POST to /api/collections/users/auth-with-password
    // ... Decode response
    
    return authResp.Token, &authResp.Record, nil
}
```

### 3. Request Scoping
A critical pattern in our architecture is **Request Scoping**. The `globalClient` holds the connection pool, but for each request, we create a lightweight clone that carries the specific user's token.

```go
// WithToken creates a shallow copy of the client with the user's token.
// It shares the underlying HTTP connection pool for efficiency.
func (c *Client) WithToken(token string) *Client {
    return &Client{
        BaseURL:    c.BaseURL,
        HTTPClient: c.HTTPClient, 
        AuthToken:  token,
    }
}
```

### 4. Listing Manifests (Fetch Data)
Retrieve the list of mission logs using the scoped client.

```go
func (c *Client) ListPosts() ([]Post, error) {
    // ... Send GET to /api/collections/posts/records
    // ... Attach Authorization header using c.AuthToken
    // ... Decode JSON into []Post
}
```

This manual control gives us precision and ensures our application remains lightweight and fast.

---
[Next: 06 - Hyperdrive UI (Alpine.js) →](./06-hyperdrive-ui.md)