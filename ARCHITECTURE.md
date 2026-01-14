# Gosmic Code - Architecture & Design Document

## 1. System Overview

**Gosmic Code** is a server-side rendered (SSR) web application designed as a "Mission Log System" for interstellar explorers. It leverages the performance of **Go (Fiber v3)**, the type-safety of **Templ**, and the rapid backend capabilities of **PocketBase**.

### High-Level Stack
*   **Language:** Go (Golang) 1.25+
*   **Web Framework:** Fiber v3 (Express-inspired, high performance)
*   **UI Engine:** Templ (Type-safe HTML generation)
*   **Database & Auth:** PocketBase (SQLite + Realtime API)
*   **Frontend Interactivity:** Alpine.js (Lightweight DOM manipulation)
*   **CSS Framework:** Pico.css (Semantic, minimal)

## 2. Business Logic & Entities

The application revolves around two core domains: **Crew Management** (Auth) and **Mission Logging** (Posts).

### PocketBase Schema

#### A. Users Collection (`users`)
*   *Default System Collection*
*   **Fields:**
    *   `id`: Unique Identifier.
    *   `email`: Communication ID.
    *   `name`: Commander Name.
    *   `avatar`: (Optional) Profile visual.

#### B. Posts Collection (`posts`)
Represents the mission logs recorded by the crew.
*   **Fields:**
    *   `title` (Text, Required): The subject of the log entry.
    *   `content` (Text): The detailed mission report.
    *   `author` (Relation -> `users`, Required): The officer who wrote the log.
    *   `public` (Boolean):
        *   `true`: Broadcast to deep space (visible to public).
        *   `false`: Encrypted (visible only to author).
*   **API Rules (Security):**
    *   **Create:** `author = @request.auth.id` (Prevents spoofing).
    *   **Update/Delete:** `author = @request.auth.id` (Ownership enforcement).
    *   **View/List:** `public = true || author = @request.auth.id`.

## 3. Application Architecture

We follow a **Modular Monolith** approach adapted for Go, emphasizing "Package by Layer" for simplicity in this scale, while adhering to separation of concerns.

### Directory Structure

```text
/
├── cmd/                    # (Optional) Application entry points
├── config/                 # Configuration loading (Env vars)
├── handlers/               # [Interface Adapter] HTTP Transport Layer
│   ├── auth.go             # Handle Login/Register flows
│   ├── posts.go            # Handle CRUD operations
│   └── ...
├── middleware/             # [Cross-Cutting Concerns]
│   ├── auth.go             # Token validation & Context injection
│   └── ...
├── pb/                     # [Infrastructure] External Service Gateway
│   └── client.go           # Typed client for PocketBase API
├── views/                  # [Presentation] UI Logic
│   ├── components/         # Reusable atoms (Buttons, Inputs)
│   ├── layouts/            # Base HTML shells
│   └── pages/              # Full page compositions
├── static/                 # Public assets (CSS, JS, Images)
└── main.go                 # Composition Root (Wiring)
```

### Architectural Patterns

#### 1. The Presentation Layer (Views)
*   Uses **Templ** to define UI components as Go functions.
*   **Logic:** Strictly display logic. No business rules or database calls inside `.templ` files.
*   **Data Flow:** Handlers pass "View Models" (structs or simple types) to Views.

#### 2. The Transport Layer (Handlers)
*   **Role:** Receive HTTP requests (Fiber Context).
*   **Responsibility:**
    1.  Parse Input (Forms, Query Params).
    2.  Validate Input (Basic formatting).
    3.  Call Infrastructure/Service Layer.
    4.  Select the appropriate View to render.
*   **Constraint:** Handlers should not contain complex business rules (e.g., "Calculate tax").

#### 3. The Infrastructure Layer (PocketBase Client)
*   **Role:** Communicate with external systems (DB).
*   **Pattern:** **Proxy Authentication**.
    *   The Go app acts as a proxy between the User and PocketBase.
    *   It holds a `globalClient` for general connection pooling.
    *   It creates a `request-scoped client` using `WithToken()` for every authenticated request.

## 4. Security Design

*   **Authentication:**
    *   Stateless JWT (JSON Web Token) issued by PocketBase.
    *   Stored in **HttpOnly, Secure, SameSite=Lax** cookies.
    *   **No** LocalStorage usage to prevent XSS token theft.
*   **CSRF Protection:**
    *   Double Submit Cookie pattern (managed by Fiber middleware).
    *   Tokens embedded in forms via `<input type="hidden" name="_csrf">`.
*   **Input Sanitization:**
    *   Handled automatically by Templ (Context-aware escaping).
*   **Authorization:**
    *   Enforced at the Database level (PocketBase API Rules).
    *   Enforced at the Route level (Auth Middleware).

## 5. Future Scalability (Clean Architecture Evolution)

As the mission grows, the architecture should evolve:

1.  **Observability:**
    *   Implement structured logging and metrics.
    *   Add distributed tracing.

2.  **Caching Layer:**
    *   Implement Redis for session storage and data caching.
