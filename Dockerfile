# Use the official Golang image to create a build artifact.
# This is a multi-stage build. In the first stage, the build environment is set up.
FROM golang:1.18-alpine as builder

# Install git, required for fetching Go dependencies.
RUN apk add --no-cache git

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code and configuration files into the container.
COPY . .

# Build the Go app.
RUN go build -o /bin/benzinga-webhook-receiver ./cmd/main.go

# Start a new stage from scratch using a smaller base image to reduce the final image size.
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /bin/benzinga-webhook-receiver .
COPY --from=builder /app/pkg/config/config.yaml /root/pkg/config/config.yaml

# Expose port 8080 to the outside world.
EXPOSE 8080

# Command to run the executable.
CMD ["./benzinga-webhook-receiver"]
