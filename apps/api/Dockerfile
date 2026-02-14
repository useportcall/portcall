FROM golang:1.25-alpine AS builder

# Install git (needed for go mod download)
RUN apk add --no-cache git

WORKDIR /workspace

# Copy libs first (dependencies)
COPY libs/ libs/

# Copy sdks (for Go SDK dependency)
COPY sdks/go/ sdks/go/

# Copy app files (including the new frontend directory)
COPY apps/api/ apps/api/
COPY apps/billing/ apps/billing/

# Set working directory to the app
WORKDIR /workspace/apps/api

# Add replace directives to use local libs
RUN go mod edit -replace github.com/useportcall/portcall/libs/go/apix=../../libs/go/apix \
  -replace github.com/useportcall/portcall/apps/billing=../../apps/billing \
  -replace github.com/useportcall/portcall/libs/go/authx=../../libs/go/authx \
  -replace github.com/useportcall/portcall/libs/go/cryptox=../../libs/go/cryptox \
  -replace github.com/useportcall/portcall/libs/go/dbx=../../libs/go/dbx \
  -replace github.com/useportcall/portcall/libs/go/discordx=../../libs/go/discordx \
  -replace github.com/useportcall/portcall/libs/go/dogfoodx=../../libs/go/dogfoodx \
  -replace github.com/useportcall/portcall/libs/go/emailx=../../libs/go/emailx \
  -replace github.com/useportcall/portcall/libs/go/envx=../../libs/go/envx \
  -replace github.com/useportcall/portcall/libs/go/logx=../../libs/go/logx \
  -replace github.com/useportcall/portcall/libs/go/paymentx=../../libs/go/paymentx \
  -replace github.com/useportcall/portcall/libs/go/qx=../../libs/go/qx \
  -replace github.com/useportcall/portcall/libs/go/ratelimitx=../../libs/go/ratelimitx \
  -replace github.com/useportcall/portcall/libs/go/routerx=../../libs/go/routerx \
  -replace github.com/useportcall/portcall/libs/go/services=../../libs/go/services \
  -replace github.com/useportcall/portcall/libs/go/storex=../../libs/go/storex \
  -replace github.com/useportcall/portcall/sdks/go=../../sdks/go

# Download dependencies
RUN go mod tidy && go mod download

# Build the Go app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server

# Build frontend
FROM node:24-bookworm-slim AS frontend-builder
WORKDIR /workspace
COPY apps/api/frontend ./apps/api/frontend
WORKDIR /workspace/apps/api/frontend
RUN npm install
RUN npm run build

# Start a new minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /workspace/apps/api/server ./server
COPY --from=builder /workspace/apps/api/openapi.yaml ./openapi.yaml
COPY --from=frontend-builder /workspace/apps/api/frontend/dist ./frontend/dist

# Command to run the executable
CMD ["./server"]
