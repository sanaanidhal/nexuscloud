# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Download dependencies first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o product-api ./cmd/main.go

# Stage 2: Minimal runtime image
FROM alpine:3.19 AS production

# Add CA certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

RUN addgroup -S nexuscloud && adduser -S appuser -G nexuscloud

WORKDIR /app

# Copy ONLY the compiled binary — nothing else
COPY --from=builder /app/product-api .

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./product-api"]