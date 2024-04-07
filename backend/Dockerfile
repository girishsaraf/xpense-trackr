FROM golang:1.22-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go clean -modcache
RUN go mod tidy
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/expense-tracker/main.go
RUN go build -o recurring_cron ./cmd/recurring-cron/main.go
# Run the Go application

# Start cron and run the Go application
CMD ["/bin/sh", "-c", "./main & ./recurring_cron && cron -f"]
