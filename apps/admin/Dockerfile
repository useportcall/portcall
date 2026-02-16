FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy libs first (dependencies)
COPY libs/ libs/

# Copy app files
COPY apps/admin/ apps/admin/

# Set working directory to the app
WORKDIR /workspace/apps/admin

# Add replace directives to use local libs
RUN go mod edit -replace github.com/useportcall/portcall/libs/go/apix=../../libs/go/apix \
  -replace github.com/useportcall/portcall/libs/go/authx=../../libs/go/authx \
  -replace github.com/useportcall/portcall/libs/go/cryptox=../../libs/go/cryptox \
  -replace github.com/useportcall/portcall/libs/go/dbx=../../libs/go/dbx \
  -replace github.com/useportcall/portcall/libs/go/dogfoodx=../../libs/go/dogfoodx \
  -replace github.com/useportcall/portcall/libs/go/emailx=../../libs/go/emailx \
  -replace github.com/useportcall/portcall/libs/go/envx=../../libs/go/envx \
  -replace github.com/useportcall/portcall/libs/go/logx=../../libs/go/logx \
  -replace github.com/useportcall/portcall/libs/go/paymentx=../../libs/go/paymentx \
  -replace github.com/useportcall/portcall/libs/go/qx=../../libs/go/qx \
  -replace github.com/useportcall/portcall/libs/go/routerx=../../libs/go/routerx \
  -replace github.com/useportcall/portcall/libs/go/storex=../../libs/go/storex

# Download dependencies
RUN go mod tidy && go mod download

# Build the Go app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app

# Build Node.js frontend
FROM node:24-bookworm-slim AS frontend-builder
WORKDIR /workspace
COPY apps/admin/frontend ./apps/admin/frontend
WORKDIR /workspace/apps/admin/frontend
# IMPORTANT: Pass VITE_KEYCLOAK_URL as build arg
# Example: --build-arg VITE_KEYCLOAK_URL=https://auth.useportcall.com
ARG VITE_KEYCLOAK_URL
ENV VITE_KEYCLOAK_URL=${VITE_KEYCLOAK_URL}
RUN npm install
RUN npm run build

# Start a new minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates curl

# Install kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
  chmod +x kubectl && \
  mv kubectl /usr/local/bin/

WORKDIR /app
COPY --from=builder /workspace/apps/admin/app ./app
COPY --from=frontend-builder /workspace/apps/admin/frontend/dist ./frontend/dist

# Command to run the executable
CMD ["./app"]
