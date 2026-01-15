# Gosmic Code - Architecture & Design Document

## 1. System Overview

**Gosmic Code** is a server-side rendered (SSR) web application designed as a "Mission Log System" for interstellar explorers. It leverages the performance of **Go (Fiber v3)**, the type-safety of **Templ**, and the rapid backend capabilities of **PocketBase**.

### High-Level Stack
*   **Language:** Go (Golang) 1.25+
*   **Web Framework:** Fiber v3 (Express-inspired, high performance)
*   **UI Engine:** Templ (Type-safe HTML generation)
*   **Database & Auth:** PocketBase v0.25+ (SQLite + Realtime API)
*   **Frontend Interactivity:** Alpine.js (Lightweight DOM manipulation)
*   **CSS Framework:** Tailwind CSS v4.0 + DaisyUI v5.0 (Modern, utility-first)

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

## 3. Application Architecture (Onion Model)

We follow an **Onion Architecture** approach, ensuring that the core business logic is independent of external concerns (like the DB or the Web Framework).

### Directory Structure

```text
/
├── handlers/               # [Interface Adapter] HTTP Transport Layer
│   ├── auth.go             # Handle Login/Register flows
│   ├── posts.go            # Handle CRUD operations
│   └── root.go             # Main dashboard/home logic
├── services/               # [Application Layer] Business Logic
│   ├── auth_service.go     # High-level auth workflows
│   └── post_service.go     # High-level post operations
├── repositories/           # [Infrastructure Adapter] Data Persistence
│   ├── auth_repository.go  # PocketBase auth communication
│   └── post_repository.go  # PocketBase post queries
├── middleware/             # [Cross-Cutting Concerns]
│   ├── auth.go             # Token validation & Context injection
│   ├── flash.go            # Flash message propagation
│   └── method_override.go  # PUT/DELETE support in forms
├── pb/                     # [Infrastructure] External Service Gateway
│   └── client.go           # Typed client for PocketBase API
├── views/                  # [Presentation] UI Logic (Templ)
│   ├── layout.templ        # Base HTML shells
│   ├── home.templ          # Landing page
│   └── posts.templ         # Mission logs dashboard
├── static/                 # Public assets (CSS, JS)
└── main.go                 # Composition Root (Dependency Injection)
```

### Architectural Patterns

#### 1. Inversion of Control & DI
*   **Repositories** implement data access.
*   **Services** consume Repositories to perform business logic.
*   **Handlers** consume Services to fulfill HTTP requests.
*   Wiring happens in `main.go`, ensuring loose coupling and easy testing.

#### 2. The Presentation Layer (Views)
*   Uses **Templ** to define UI components.
*   **Data Flow:** Handlers pass models to Views.
*   **UX:** Enhanced with **Alpine.js** for client-side state and **HTMX** for seamless updates.

#### 3. Proxy Authentication Pattern
*   The Go app acts as a secure buffer between the User and PocketBase.
*   **Request-Scoped Client:** Uses `middleware.AuthMiddleware` to inject an authenticated PB client into the context for each request, solving shared-state concurrency bugs.

## 4. Security Design

*   **Authentication:** HttpOnly, Secure, SameSite=Lax cookies for PB Tokens.
*   **CSRF Protection:** Double Submit Cookie pattern via Fiber CSRF middleware.
*   **Input Sanitization:** Context-aware escaping by Templ.
*   **A11Y (Accessibility):** Semantic HTML5, Skip links, and ARIA attributes for screen readers.

## 5. Deployment & Infrastructure

*   **Docker:** Multi-stage build (Alpine-based) with Go and Bun.
*   **Port:** Standardized to `8080`.
*   **Environment:** Production-ready configuration via `GO_ENV=production`.
*   **Tunnels:** Integration with `cloudflared` for local development exposure.
