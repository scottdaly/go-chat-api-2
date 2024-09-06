# Use the official Golang image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

RUN GOPROXY=https://goproxy.cn go get github.com/klauspost/compress

# Download dependencies
RUN go mod download

# Now copy the rest of the source code
COPY . .

# Build the application
RUN go build -o main .

# Create a volume for the SQLite database
VOLUME /app/data

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]