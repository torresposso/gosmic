# 08 - Log Operations (Method Override) 📝

> **Rank**: Commodore  
> **Skill**: REST, HTML Forms, and Method Override

**Mission**: Implement the core CRUD (Create, Read, Update, Delete) operations using standard web patterns.

## 🧠 Theory: The Limits of HTML Forms

Standard HTML forms only support `GET` and `POST` methods. They do not natively support `PUT`, `PATCH`, or `DELETE`. This is a historical limitation of the web specifications.

To build a RESTful system (`DELETE /posts/123`), we use a technique called **Method Override**.

### How Method Override Works
1.  The Form sends a `POST` request.
2.  The Form includes a hidden field named `_method` with the desired action (e.g., `DELETE`).
    ```html
    <input type="hidden" name="_method" value="DELETE"/>
    ```
3.  **Fiber Middleware** (`middleware.MethodOverride`) intercepts the request.
4.  It sees the `_method` field and rewrites the internal request method from `POST` to `DELETE`.
5.  The router matches `app.Delete(...)` instead of `app.Post(...)`.

## 📜 The Log Archivist (`handlers/posts.go`)

### 1. Recording Logs (Create)
We use a standard form submission. Fiber's CSRF middleware validates the request before it reaches our handler.

```go
func CreatePost() fiber.Handler {
    return func(c fiber.Ctx) error {
        // ... Validate Input ...
        
        // Persist to Cargo Hold
        // Note: We don't send the Author ID. PocketBase handles that securely.
        err := client.CreatePost(title, content, isPublic)
        
        // Post-Redirect-Get (PRG) Pattern
        return c.Redirect().To("/dashboard/posts")
    }
}
```

### 2. Purging Records (Delete)

In `views/posts.templ`, we create a form that looks like a button but acts like a secure delete request.

```go
<form method="POST" action={ templ.SafeURL("/dashboard/posts/" + post.ID) }>
    <!-- The Secret Handshake -->
    <input type="hidden" name="_method" value="DELETE"/>
    <input type="hidden" name="_csrf" value={ csrfToken }/>
    
    <button type="submit" class="danger">Purge Record</button>
</form>
```

When clicked:
1.  Browser sends `POST /dashboard/posts/123`.
2.  Fiber overrides method to `DELETE`.
3.  `handlers.DeletePost` executes.
4.  User is redirected back to the log list.

This approach ensures our application works robustly across all browsers and doesn't rely on JavaScript for critical data operations.

---
[Next: 09 - Simulation (Testing) →](./09-simulation.md)