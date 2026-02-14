FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy libs first (dependencies)
COPY libs/ libs/

# Copy app files
COPY apps/quote/ apps/quote/

# Set working directory to the app
WORKDIR /workspace/apps/quote

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

# Build the Go app binary to a path that cannot collide with ./app source dir.
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o quote-server .

# Start a new minimal image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /workspace/apps/quote/quote-server ./quote-server
COPY --from=builder /workspace/apps/quote/templates ./templates

# Command to run the executable
CMD ["./quote-server"]
