FROM golang:1.23.1-alpine3.20
WORKDIR /gocrud

# Install any required packages (like git for fetching dependencies)
# RUN apk add --no-cache git

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o go-app
# Expose the application port (optional, adjust based on your app)
EXPOSE 8083

# Run the built binary
CMD ["./go-app"]