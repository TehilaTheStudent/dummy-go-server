# Start from the official Golang image as a builder
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o server .

# Start a minimal image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server ./server
COPY openapi.yaml ./openapi.yaml
EXPOSE 8080
ENV PORT=8080
CMD ["./server"]
