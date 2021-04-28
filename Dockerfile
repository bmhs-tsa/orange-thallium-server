# Build container
FROM golang:1.16-alpine AS build

# Working directory
WORKDIR /src

# Copy source code
COPY . .

# Install Go dependencies
RUN go get -d -v .
RUN go install -v .

# Build the server
RUN go build -o /out/app .

# Production container
FROM alpine:latest

# Metadata
LABEL security="wakefulcloud@protonmail.com"

# Create a system user
RUN addgroup -S go && adduser -S go -G go

# Working directory
WORKDIR /go

# Copy binary
COPY --chown=go:go --from=build /out/app .

# Copy default config
COPY config/default.toml config/

# Switch to the system user
USER go

# Start the server
CMD [ "/go/app" ]