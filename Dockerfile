FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o fetch-kit .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/fetch-kit .

# Expose port for API server
EXPOSE 8080

ENTRYPOINT ["/app/fetch-kit"]

CMD [] 