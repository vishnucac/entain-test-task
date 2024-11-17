# Use Golang base image
FROM golang:1.23-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the host into the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Build the Go app
RUN go build -o main .

# Run unit tests
RUN go test ./... -v

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
