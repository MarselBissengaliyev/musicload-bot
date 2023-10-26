# Use the official Go image
FROM golang:1.21.3  

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project to the working directory
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Set the command to run the binary
CMD ["./main"]
