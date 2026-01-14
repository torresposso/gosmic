# 10 - Colonization (Production) 🛸

**Mission Phase**: Deployment  
**Objective**: Deploy the application to live servers.

## 🏭 Manufacturing (Build)

The ship must be assembled before launch to ensure structural integrity.

1.  **Build System**: Compile the binary.
    ```bash
    task build
    ```
2.  **Launch**:
    ```bash
    ./bin/app
    ```

## 📦 Containerization (Docker)

For deep space travel (cloud deployment), encapsulate the ship in a standard container.

```dockerfile
# Stage 1: Builder
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install build tools
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templates and build
RUN templ generate
RUN go build -o main .

# Stage 2: Runner
FROM alpine:latest
WORKDIR /root/

# Copy artifacts from builder
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

# Expose port and launch
EXPOSE 8080
CMD ["./main"]
```

## 🚩 Final Words from Command

You have successfully built the **Gosmic Code** mission log system. You have mastered:

- **Fiber v3**: High-performance HTTP routing.
- **Templ**: Type-safe, compiled HTML rendering.
- **PocketBase**: A self-contained backend solution.
- **Alpine.js**: Lightweight interactivity.

**Mission Accomplished.** 🚀