 # Use the official Golang image with version 1.21.5 as the base image
FROM golang:1.21.5

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download

# Copy the local package files to the container's workspace
COPY . .

# Build the Golang application
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o myapp ./cmd/twitter-clone-http/main.go

# Expose the port that the application will run on (change to your desired port)
EXPOSE 9090

# Command to run the executable
CMD ["./myapp"]
