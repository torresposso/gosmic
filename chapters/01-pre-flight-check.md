# 01 - Pre-Flight Check 🚀

**Mission Phase**: Initialization  
**Objective**: Prepare the development environment and understand the ship's systems.

## 🛠️ Prerequisites

Before boarding, ensure you have the following tools installed to maintain life support and propulsion:

1.  **Go 1.25+**: The core engine.
    *   *Verification*: `go version`
2.  **Task**: A task runner / build tool for executing launch commands.
    ```bash
    go install github.com/go-task/task/v3/cmd/task@latest
    ```
3.  **Templ**: For generating type-safe UI components.
    ```bash
    go install github.com/a-h/templ/cmd/templ@latest
    ```
4.  **PocketBase**: The backend data system (Cargo Hold).
    *   Download the appropriate binary for your OS.

## 🏗️ Initialize the Project

1.  **Clone the Mission Repo**
    ```bash
    git clone https://github.com/torresposso/gosmic.git
    cd gosmic
    ```

2.  **Install Dependencies**
    Download the required Go modules to the ship's computer.
    ```bash
    go mod download
    ```

## 📂 Ship Schematic (Project Structure)

Understanding the layout of the ship is crucial for survival. We follow an **Onion Architecture** to keep the core systems isolated and testable.

```text
gosmic/
├── main.go                    # Bridge: Application entry point and router configuration.
├── Taskfile.yml               # Launch Codes: Shortcuts for build, dev, and test commands.
├── .env                       # Signal Calibration: Environment variables (PORT, PB_URL).
├── handlers/                  # Crew: HTTP Handlers (Presentation Layer).
│   ├── auth.go               # Security: Login, Register, Logout logic.
│   ├── root.go               # Command: Dashboard and Landing page logic.
│   └── posts.go              # Operations: Mission Logs CRUD.
├── services/                  # Officers: Business Logic Layer.
│   ├── auth_service.go       # Auth rules and token management.
│   └── post_service.go       # Post creation and retrieval logic.
├── repositories/              # Engineering: Data Access Layer.
│   ├── auth_repo.go          # Database interactions for Users.
│   └── post_repo.go          # Database interactions for Posts.
├── middleware/                # Shields: Request interceptors (Auth, Flash, etc).
├── pb/
│   └── client.go             # Cargo Interface: Typed wrapper for PocketBase API.
├── views/                     # Windows: UI Components (Templ files).
│   ├── layout.templ          # Hull: Base HTML structure.
│   ├── home.templ            # Landing & Dashboard.
│   ├── auth.templ            # Login/Register Forms.
│   └── posts.templ           # Post lists and forms.
└── static/                    # Paint: CSS (Tailwind) and JS (Alpine.js/HTMX).
```

## 🚀 First Launch

1.  **Start PocketBase** (in a separate terminal)
    ```bash
    ./pocketbase serve --http=0.0.0.0:8090
    ```
    *Ensure the backend is reachable at port 8090.*

2.  **Ignite Main Thrusters**
    Use the task runner to start the application in development mode (with live reload).
    ```bash
    task dev
    ```

3.  **Access Command Center**
    Open your browser visualizer to `http://localhost:8080`.

## 🧠 Theory: Why Templ?

In standard Go `html/template`, errors are often discovered at **runtime**. A typo like `{{ .Titlee }}` causes the application to crash when the page is visited.

**Templ** solves this by compiling views into standard Go code.
*   **Compile-Time Safety**: If you reference a missing field, the code won't compile.
*   **Performance**: Templates are compiled to optimized Go functions, writing directly to the output buffer for maximum speed.
*   **Security**: Context-aware escaping is applied automatically, protecting against XSS (Cross-Site Scripting) anomalies.

---
[Next: 02 - Ignition (Hello Fiber) →](./02-ignition.md)