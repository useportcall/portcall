FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Install git (required for go mod download)
RUN apk add --no-cache git

# Copy libs first (dependencies)
COPY libs/ libs/

# Copy app files
COPY apps/file/ apps/file/

# Set working directory to the app
WORKDIR /workspace/apps/file

# Add replace directives to use local libs
RUN go mod edit -replace github.com/useportcall/portcall/libs/go/apix=../../libs/go/apix \
  -replace github.com/useportcall/portcall/libs/go/authx=../../libs/go/authx \
  -replace github.com/useportcall/portcall/libs/go/cryptox=../../libs/go/cryptox \
  -replace github.com/useportcall/portcall/libs/go/dbx=../../libs/go/dbx \
  -replace github.com/useportcall/portcall/libs/go/emailx=../../libs/go/emailx \
  -replace github.com/useportcall/portcall/libs/go/envx=../../libs/go/envx \
  -replace github.com/useportcall/portcall/libs/go/logx=../../libs/go/logx \
  -replace github.com/useportcall/portcall/libs/go/paymentx=../../libs/go/paymentx \
  -replace github.com/useportcall/portcall/libs/go/qx=../../libs/go/qx \
  -replace github.com/useportcall/portcall/libs/go/ratelimitx=../../libs/go/ratelimitx \
  -replace github.com/useportcall/portcall/libs/go/routerx=../../libs/go/routerx \
  -replace github.com/useportcall/portcall/libs/go/storex=../../libs/go/storex

# Download dependencies
RUN go mod tidy && go mod download

# Build the Go app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server

# Start a new minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /workspace/apps/file/server ./server
COPY --from=builder /workspace/apps/file/templates ./templates

# Command to run the executable
CMD ["./server"]
