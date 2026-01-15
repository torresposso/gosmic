# 09 - Simulation (Testing) 🧪

> **Rank**: Admiral  
> **Skill**: Unit Testing, Mocking, and Dependency Injection

**Mission**: Verify ship systems in a controlled environment (Test Suite) without relying on external variables (PocketBase).

## 🧩 The Strategy: Interfaces & Mocks

Why check systems in a simulation?
1.  **Speed**: No network latency. Tests run in milliseconds.
2.  **Reliability**: No checking if the Database is "up".
3.  **Isolation**: We test *our* logic, not PocketBase's logic.

### Dependency Injection (DI)
Our handlers accept Repositories or Services interfaces. To test them, we can use mocks.

However, for `pb.Client` which is a struct, we mock the **HTTP Transport**.

### 1. The Mock Uplink (`MockRoundTripper`)
Go's `http.Client` uses a `RoundTripper` interface to execute requests. By replacing this, we can intercept and control the responses.

```go
type MockRoundTripper struct {
    RoundTripFunc func(req *http.Request) *http.Response
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.RoundTripFunc(req), nil
}
```

### 2. Testing Handlers (`handlers/posts_test.go`)

We use `app.Test(req)` from Fiber to inject requests directly into the router.

```go
func TestCreatePost(t *testing.T) {
    // A. Setup
    mockService := new(mocks.PostService)
    handler := handlers.NewPostHandler(mockService, store)
    app := fiber.New()
    app.Post("/posts", handler.Create())

    // B. Execute
    req := httptest.NewRequest("POST", "/posts", body)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, _ := app.Test(req)

    // C. Verify
    assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
}
```

### 3. View Testing (`views/views_test.go`)
Since Templ components are just functions that write to a buffer, we can test our UI without a browser.

```go
func TestDashboardView(t *testing.T) {
    buf := new(bytes.Buffer)
    // Render the component to buffer
    err := Dashboard("Shepard", ...).Render(context.Background(), buf)
    
    assert.NoError(t, err)
    assert.Contains(t, buf.String(), "Welcome aboard, Commander Shepard!")
}
```

## 🏃 Running Simulations

We use our launch codes (`Taskfile.yml`) to run the suite with filtered output, hiding configuration noise and focusing on the core systems.

```bash
# Run all verified tests
task test

# Run tests in verbose mode
task test:v
```

> [!NOTE]
> **Commander's Log**: We recently conducted a full audit (`TESTING_AUDIT.md`) and significantly improved coverage for `handlers/posts.go`. Always keep your simulations up to date!

---
[Next: 10 - Colonization (Production & Deployment) →](./10-colonization.md)