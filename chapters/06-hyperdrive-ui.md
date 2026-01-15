# 06 - Hyperdrive UI (HTMX & Alpine.js) ⚡

**Mission Phase**: Engine Tuning  
**Objective**: Enhance user interaction with lightweight JavaScript.

## 🚀 The Philosophy: Server-Driven UI

In this mission, we prioritize **Stability** and **Simplicity**. We rely on standard HTML Forms and Server-Side Rendering (SSR) for the core navigation. However, we boost the user experience using **HTMX** and **Alpine.js**.

## 🧬 HTMX (Sublight Communication)

HTMX allows us to access AJAX, CSS Transitions, WebSockets and Server Sent Events directly in HTML.

### 1. Boosting Links
In `layout.templ`, we enable `hx-boost` on the body. This turns standard anchors and forms into AJAX requests, avoiding full page reloads and making the app feel like an SPA.

```html
<body hx-boost="true"> ... </body>
```

### 2. Live Toggles
For the "Toggle Public/Private" feature on posts, we use HTMX to swap just the button, not the whole page.

```html
<button hx-post={ "/api/posts/" + post.Id + "/toggle" } 
        hx-target="closest div" 
        hx-swap="outerHTML">
    Make Private
</button>
```

## 🤖 Droids (Alpine.js)

Alpine.js allows us to sprinkle interactivity directly into our HTML markup for client-side state.

### Example: Auto-Dismissing Alerts
Flash messages should appear briefly and then vanish.

In `views/layout.templ`:

```go
templ FlashMessage(flash string, flashType string) {
    if flash != "" {
        <div id="flash-message" 
             hx-swap-oob="true" 
             x-data="{ show: true }" 
             x-show="show" 
             x-init="setTimeout(() => show = false, 8000)">
            
            <div class={ "alert " + flashType }>
                <span>{ flash }</span>
            </div>
        </div>
    }
}
```

**How it works:**
1.  **`x-data`**: Initializes local state `show = true`.
2.  **`x-init`**: Sets a timer to hide the alert after 8 seconds.
3.  **`hx-swap-oob`**: Allows the server to inject this message even if it wasn't the main target of the HTMX request.

This combination of **HTMX** for server communication and **Alpine.js** for UI state gives us a powerful "Hyperdrive" experience without the weight of React or Vue.

---
[Next: 07 - Security Clearance (Auth System) →](./07-security-clearance.md)