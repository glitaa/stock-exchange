# ==========================================
# Build Stage
# ==========================================
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . .

# Compile the application as a statically linked binary (CGO_ENABLED=0)
# Target Linux OS to ensure compatibility with the Alpine base image
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/exchange-server ./cmd/server/main.go

# ==========================================
# Final Stage
# ==========================================
FROM alpine:3.23

# Install CA certificates for secure HTTPS connections if needed
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/exchange-server .

EXPOSE 8080

ENTRYPOINT ["./exchange-server"]
CMD ["-port", "8080"]