# Use the official Golang image as the base image
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main /app/cmd/app/main.go

# Use a minimal image for running the application
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/.env .
COPY --from=builder /app/fonts ./fonts
COPY --from=builder /app/images ./images

# Expose the application's port
EXPOSE 8080

# Run the application
CMD ["./main"]