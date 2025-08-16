# ---- Builder ----
FROM golang:1.23.1-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

# Copy go.mod + go.sum để cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source
COPY . .

# Build binary
RUN go build -o hub-service main.go

# ---- Runner ----
FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/hub-service .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./hub-service"]