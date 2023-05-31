# Use the official Golang image as the base image
FROM golang:1.19-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module dependency files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o golang-rest-api .

# Expose the port the application listens on
EXPOSE 8080

# Define the command to run when the container starts
CMD ["./golang-rest-api","serve"]