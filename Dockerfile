# Use the official Golang image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download all dependencies
RUN go mod download

# Build the application
RUN sudo go build -o main .

# Create a volume for the SQLite database
VOLUME /app/data

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]