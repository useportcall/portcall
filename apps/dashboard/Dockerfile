FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy libs first (dependencies)
COPY libs/ libs/

# Copy sdks (for Go SDK dependency)
COPY sdks/go/ sdks/go/

# Copy app files (excluding frontend)
COPY apps/dashboard/ apps/dashboard/
COPY apps/billing/ apps/billing/

# Set working directory to the app
WORKDIR /workspace/apps/dashboard

# Add replace directives to use local libs and SDKs
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

# build nodejs frontend
FROM node:24-bookworm-slim AS frontend-builder
WORKDIR /workspace
COPY apps/dashboard/frontend ./apps/dashboard/frontend
WORKDIR /workspace/apps/dashboard/frontend
# IMPORTANT: Pass VITE_KEYCLOAK_URL as build arg
# Example: --build-arg VITE_KEYCLOAK_URL=https://auth.useportcall.com
ARG VITE_KEYCLOAK_URL
ENV VITE_KEYCLOAK_URL=${VITE_KEYCLOAK_URL}
RUN npm install
RUN npm run build

# Start a new minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /workspace/apps/dashboard/server ./server
COPY --from=frontend-builder /workspace/apps/dashboard/frontend/dist ./frontend/dist

# Command to run the executable
CMD ["./server"]
