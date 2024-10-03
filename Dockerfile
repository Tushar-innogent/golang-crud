FROM golang:1.23.1-alpine3.20
# WORKDIR /gocrud
# COPY . /gocrud
# RUN go build /gocrud
# EXPOSE 8081
# ENTRYPOINT [ "./go-crud" ]
# CMD [ "./go-crud" ]

# Use a Go image as the base image
# FROM golang:1.23.1-alpinem

# Install any required packages (like git for fetching dependencies)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o myapp

# Expose the application port (optional, adjust based on your app)
EXPOSE 8081

# Run the built binary
CMD ["./myapp"]