# Multi-stage build for minimal, secure image
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o watchtower-masterbot .

# Final minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Create non-root user and data directory
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /app/data && \
    chown -R appuser:appgroup /app

USER appuser

COPY --from=builder /app/watchtower-masterbot .
COPY --from=builder /app/config ./config

CMD ["./watchtower-masterbot"]