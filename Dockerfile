FROM golang:1.22.4-alpine
WORKDIR /app

# Copy go.mod and go.sum first for dependency management
COPY go.mod go.sum ./
RUN go mod download

# Copy Go source files
COPY ./*.go ./
COPY internal ./internal

# Copy the other required files and directories
COPY .env ./
COPY index.html ./
COPY assets ./assets

# Build the Go application
RUN go build -o web-server .

EXPOSE 8080

# Command to run the built binary
CMD ["/app/web-server"]
