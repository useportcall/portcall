FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app

# Start a new minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /
COPY --from=builder /app /app

# Command to run the executable
CMD ["./app/app"]
