# ---- Builder ----
FROM golang:1.23.1-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o hub-service main.go

# ---- Runner ----
FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/hub-service .

EXPOSE 8080

CMD ["./hub-service"]