# Step 1: Use an official Golang image to build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy ALL files
COPY . /app/

# Build the Go application
RUN ls -la
RUN CGO_ENABLED=1 go build -o main .

# Step 2: Use a minimal image to run the application
FROM ubuntu:24.04
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main ./
#COPY gorm.db /app/gorm.db
COPY templates/ /app/templates/
COPY assets/ /app/assets/

# Expose the port your application listens on
EXPOSE 8080

# IDK
#COPY main ./main

# Set the command to run the application
CMD ["./main"]
