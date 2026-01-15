# Build Stage
FROM golang:1.25.5-alpine AS builder

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templates
RUN templ generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for secure connections (PocketBase etc)
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /bin/app /app/app

# Copy static assets
COPY --from=builder /app/static /app/static

# Copy documentation chapters
COPY --from=builder /app/chapters /app/chapters

# Default environment variables
ENV PORT=8080
ENV GO_ENV=production

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"]

