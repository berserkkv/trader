# Build stage
FROM golang:1.24-alpine as build

WORKDIR /app

# Copy the Go modules file and install dependencies first (to leverage Docker cache)
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the entire source code
COPY . .

# Build the application. Ensure the binary is outputted as "trading-server"
RUN go build -o trader ./cmd/main.go

# Final stage: create the runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=build /app/trader .
COPY --from=build /app/users.db /root
COPY --from=build /app/config.yaml /root

# Expose the port the application will listen on
EXPOSE 8080

# Set the default command to run the app
CMD ["./trader"]
