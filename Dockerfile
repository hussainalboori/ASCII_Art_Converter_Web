# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port on which the application will run (adjust as needed)
EXPOSE 8080

# Define environment variables (if needed)
# ENV MY_VARIABLE my_value

# Add a label to the Docker image
LABEL maintainer="hussain alboori hussain nabeel ali alsandi"
LABEL version="1.0"
LABEL description="This Go application creates a web-based ASCII art converter"

# Run the Go application
CMD ["./main"]
