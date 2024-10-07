# Use the official Golang image
FROM golang:alpine

# Install necessary packages
RUN apk add --no-cache bash

# Create user mt
RUN adduser -D mt

# Set the working directory inside the container
WORKDIR /app

# Copy the Go file into the container
COPY ./app/main.go .

# Change ownership of the working directory to user mt
RUN chown -R mt:mt /app

# Switch to user mt
USER mt

# Build the Go application
RUN #go build -o build/main main.go

EXPOSE 9099
# Specify the command to run the application
#CMD ["./build/main"]
