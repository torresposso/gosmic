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

## 🏗️ Initialize the Project

1.  **Clone the Mission Repo**
    ```bash
    git clone https://github.com/torresposso/gosmic-code.git
    cd gosmic-code
    ```

2.  **Install Dependencies**
    Download the required Go modules to the ship's computer.
    ```bash
    go mod download
    ```

## 📂 Ship Schematic (Project Structure)

Understanding the layout of the ship is crucial for survival.

```text
gosmic/
├── main.go                    # Bridge: Application entry point and router configuration.
├── Taskfile.yml               # Launch Codes: Shortcuts for build, dev, and test commands.
├── .env                       # Signal Calibration: Environment variables (PORT, PB_URL).
├── handlers/                  # Crew: Core logic handling HTTP requests.
│   ├── render.go             # Translator: Helper to render Templ components.
│   ├── auth.go               # Security: Login, Register, and Session management.
│   ├── home.go               # Command: Dashboard and Landing page logic.
│   └── posts.go              # Operations: CRUD logic for Mission Logs.
├── pb/
│   └── client.go             # Cargo Interface: Typed wrapper for the PocketBase API.
├── views/                     # Windows: UI Components (Templ files).
│   ├── layout.templ          # Hull: Base HTML structure (Head, Nav, Footer).
│   ├── home.templ            # Views for public landing and user dashboard.
│   ├── auth.templ            # Forms for identification and enlistment.
│   └── posts.templ           # Views for listing and editing logs.
└── static/                    # Paint: CSS (Pico.css) and JS (Alpine.js) assets.
```

## 🚀 First Launch

1.  **Start PocketBase** (in a separate terminal)
    ```bash
    ./pocketbase serve
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