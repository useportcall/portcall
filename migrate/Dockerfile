FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy libs first (dependencies)
COPY libs/ libs/

# Copy migrate tool
COPY migrate/ migrate/

# Set working directory to migrate
WORKDIR /workspace/migrate

# Add replace directives to use local libs
RUN go mod edit -replace github.com/useportcall/portcall/libs/go/dbx=../libs/go/dbx \
  -replace github.com/useportcall/portcall/libs/go/envx=../libs/go/envx

# Download dependencies
RUN go mod tidy && go mod download

# Build the migration tool
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o migrate

# Start a new minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /workspace/migrate/migrate ./migrate

# Command to run the migration
CMD ["./migrate"]
