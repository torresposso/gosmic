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
Our handlers accept a `*pb.Client`. To test them, we need to mock the *layer below it*—the HTTP Transport.

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

### 2. The Verification Logic (`handlers/login_integration_test.go`)

We use `app.Test(req)` from Fiber to inject requests directly into the router.

```go
func TestLoginFlowWithCSRF(t *testing.T) {
    // A. Setup the Simulation (Mock Response)
    mockTripper := &MockRoundTripper{
        RoundTripFunc: func(req *http.Request) *http.Response {
            if strings.Contains(req.URL.Path, "auth-with-password") {
                // Return a fake token and user record
                return &http.Response{
                    StatusCode: http.StatusOK,
                    Body:       io.NopCloser(bytes.NewBuffer(fakeSuccessJSON)),
                }
            }
            return &http.Response{StatusCode: 404}
        },
    }

    // B. Inject the Mock
    pbClient := pb.NewClient("http://mock-pb")
    pbClient.HTTPClient.Transport = mockTripper

    // C. Run the Scenario
    // 1. GET /login (Retrieve CSRF)
    // 2. POST /login (Submit Credentials + CSRF)
    
    // D. Assert Outcome
    assert.Equal(t, http.StatusSeeOther, resp.StatusCode) // Redirect to Dashboard
}
```

### 3. View Testing (`views/views_test.go`)
Since Templ components are just functions that write to a buffer, we can test our UI without a browser.

```go
func TestDashboardView(t *testing.T) {
    buf := new(bytes.Buffer)
    err := Dashboard("Shepard", "shepard@sr2.com", 5, "token").Render(context.Background(), buf)
    assert.NoError(t, err)

    content := buf.String()
    assert.Contains(t, content, "Welcome aboard, Commander Shepard!")
}
```

### 4. API Client Testing (`pb/client_test.go`)
We verify our `pb.Client` by spawning a local `httptest.Server` that acts like a fake PocketBase.

```go
func TestAuthWithPassword(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(authResponse{Token: "test-token"})
    }))
    defer server.Close()

    client := NewClient(server.URL)
    token, _, _ := client.AuthWithPassword("test@email.com", "pass")
    assert.Equal(t, "test-token", token)
}
```

## 跑 Running Simulations

We use our launch codes (`Taskfile.yml`) to run the suite with filtered output, hiding configuration noise and focusing on the core systems.

```bash
# Run all verified tests
task test

# Run tests in verbose mode
task test:v
```

> [!NOTE]
> **Commander's Log**: Real-world tests often require setting up middleware (like CSRF) within the test function to match the production environment, as seen in `login_integration_test.go`.

---
[Next: 10 - Colonization (Production & Deployment) →](./10-colonization.md)