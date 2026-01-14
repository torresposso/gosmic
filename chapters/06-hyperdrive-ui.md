# 06 - Hyperdrive UI (Alpine.js) ⚡

**Mission Phase**: Engine Tuning  
**Objective**: Enhance user interaction with lightweight JavaScript.

## 🚀 The Philosophy: Server-Driven UI

In this mission, we prioritize **Stability** and **Simplicity**. We rely on standard HTML Forms and Server-Side Rendering (SSR) for the core navigation and data mutation. This ensures the ship works even if the navigation computer (JavaScript) is damaged.

However, for ephemeral interactions—like dismissing alerts or toggling visibility—we employ **Alpine.js**.

## 🤖 Droids (Alpine.js)

Alpine.js allows us to sprinkle interactivity directly into our HTML markup without writing complex JavaScript bundles. It is the "jQuery for the modern web."

### Example: Auto-Dismissing Alerts
Flash messages (like "Login Successful") should appear briefly and then vanish.

In `views/layout.templ`:

```go
templ FlashMessage(flash string, flashType string) {
    if flash != "" {
        // x-data: Initialize state
        // x-show: Bind visibility to state
        // x-init: Execute logic on load
        <div class={ "flash " + flashType } 
             x-data="{ show: true }" 
             x-show="show" 
             x-transition 
             x-init="setTimeout(() => show = false, 5000)">
            { flash }
        </div>
    }
}
```

**How it works:**
1.  **`x-data="{ show: true }"`**: Defines a local state variable `show` initialized to `true`.
2.  **`x-show="show"`**: The element is visible only when `show` is true.
3.  **`x-init="..."`**: When the element mounts, it sets a timer. After 5000ms, `show` becomes `false`, hiding the element.
4.  **`x-transition`**: Alpine automatically applies CSS transitions for a smooth fade-out.

This approach keeps our logic co-located with our view, making the codebase easier to maintain and reason about.

---
[Next: 07 - Security Clearance (Auth System) →](./07-security-clearance.md)