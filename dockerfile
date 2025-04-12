# Build stage
FROM golang:1.20-alpine as build

WORKDIR /app

# Copy the Go modules file and install dependencies first (to leverage Docker cache)
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire source code
COPY . .

# Build the application. Ensure the binary is outputted as "trading-server"
RUN go build -o trading-server ./cmd/main.go

# Final stage: create the runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=build /app/trading-server .

# Expose the port the application will listen on
EXPOSE 8080

# Set the default command to run the app
CMD ["./trading-server"]
