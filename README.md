# 🚀 Gosmic Code: The Hypermedia-Driven Mission

![License](https://img.shields.io/badge/license-MIT-blue.svg) ![Go Version](https://img.shields.io/badge/go-1.22%2B-cyan) ![Fiber](https://img.shields.io/badge/fiber-v3-green)

**Welcome to the flight deck, Commander.**

You are about to embark on a mission to build a **production-grade**, **hyper-fast**, and **maintainable** web application. We aren't just building a "To-Do list"; we are engineering a **Mission Log System** capable of withstanding the rigors of deep space (and high traffic).

## 🌌 The Mission Objective

To master the **Go-HTMX-PocketBase** stack—a potent combination that delivers the interactivity of a Single Page App (SPA) with the simplicity and performance of Server-Side Rendering (SSR).

### Why this Stack? (The Engineering Brief)

Modern web development often feels like navigating an asteroid field of complexity: separate backends, hydration issues, state synchronization hell. **We choose a different path.**

| Component | The Role | Why It's Superior |
|-----------|----------|-------------------|
| **Fiber v3** | **Thrusters** | Built on `fasthttp`, it offers **zero-allocation** performance. It's the fastest Go framework, period. |
| **PocketBase** | **Cargo Hold** | An embedded, real-time backend (SQLite + Go). Zero latency database access, built-in Auth, and easy deployment. |
| **Templ** | **Blueprints** | Type-safe HTML. If your code compiles, your UI is correct. No more runtime template errors. |
| **Alpine.js** | **Droids** | For the 10% of interactivity that *must* happen instantly firmly on the client (dropdowns, modals). |

## 🏗️ Technical Architecture

We adhere to a **Server-Side Rendering (SSR)** architecture.

```mermaid
sequenceDiagram
    participant User
    participant Browser
    participant Server (Fiber+Templ)
    participant Database (PocketBase)

    User->>Browser: Clicks "Archive Log"
    Browser->>Server (Fiber+Templ): POST /logs/123 (Method Override: DELETE)
    Server (Fiber+Templ)->>Database (PocketBase): Delete Record
    Database (PocketBase)-->>Server (Fiber+Templ): Confirm
    Server (Fiber+Templ)-->>Browser: Return HTML Page Redirect/Content
    Browser->>Browser: Full Page Reload/Render
```

**Key Takeaway**: The server returns complete **HTML** pages.

## 📚 Flight Manual (Curriculum)

This course is structured to take you from Cadet to Commander.

| Rank | Mission Segment | Skill Acquired |
|------|-----------------|----------------|
| **Cadet** | [01 - Pre-Flight Check](./chapters/01-pre-flight-check.md) | **Architecture Design**: Understanding the folder structure and module system. |
| **Ensign** | [02 - Ignition](./chapters/02-ignition.md) | **Fiber Deep Dive**: Handlers, zero-allocation context, and middleware chains. |
| **Lt. JG** | [03 - Navigation](./chapters/03-navigation.md) | **Routing**: Groups, parameters, and constraints. |
| **Lieutenant** | [04 - Life Support](./chapters/04-life-support.md) | **Middleware**: Building custom interceptors for logging and observability. |
| **Lt. Cmdr** | [05 - Cargo Hold](./chapters/05-cargo-hold.md) | **Data Integration**: Embedding PocketBase and type-safe record handling. |
| **Commander** | [06 - UI Components](./chapters/06-ui-components.md) | **Templ**: Creating reusable, composable UI components. |
| **Captain** | [07 - Security Clearance](./chapters/07-security-clearance.md) | **AppSec**: JWT signing, HttpOnly cookies, and **CSRF protection**. |
| **Commodore** | [08 - Log Operations](./chapters/08-log-operations.md) | **Full CRUD**: Implementing standard RESTful handlers and forms. |
| **Admiral** | [09 - Simulation](./chapters/09-simulation.md) | **Testing**: Unit testing handlers and mocking HTTP transports. |
| **Fleet Admiral** | [10 - Colonization](./chapters/10-colonization.md) | **DevOps**: Dockerizing and deploying to the cloud. |

## 🚀 Launch Sequence

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/torresposso/gosmic-code.git
    cd gosmic-code
    ```

2.  **Ignite Engines**:
    ```bash
    # Install dependencies and tools
    task install
    
    # Start the development server (Air + Templ + PocketBase)
    task dev
    ```

3.  **Access Command Center**:
    Open `http://localhost:8080`.

---

**"Ad Astra Per Aspera"** (To the stars through difficulties)
*But with this stack, the difficulties are significantly reduced.* 🌠
