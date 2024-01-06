 # Use the official Golang image with version 1.21.5 as the base image
FROM golang:1.21.5 as build-base

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

RUN useradd -u 1001 nonroot

# Copy the local package files to the container's workspace
COPY . .

# Build the Golang application
RUN go build \
    --ldflags="-linkmode external -extldflags -static" \
    -o myapp ./cmd/twitter-clone-http/main.go
 
FROM scratch

WORKDIR /

COPY --from=build-base /etc/passwd /etc/passwd

COPY --from=build-base --chown=nonroot:nonroot /app/config.json /

COPY --from=build-base --chown=nonroot:nonroot /app/myapp /

USER nonroot

# Command to run the executable
CMD ["/myapp"]
