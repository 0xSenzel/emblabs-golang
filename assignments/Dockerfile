# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy and download dependencies (for layer caching efficiency)
COPY go.mod ./
# Copy go.sum if it exists (using wildcard pattern makes it optional)
COPY go.su[m] ./ 
RUN go mod download

# Copy source code and build the single binary executable
COPY . .
# -o /app/api: names the output binary 'api' in the root of /app
RUN go build -o /app/api ./cmd/api

# Stage 2: Runner
FROM alpine:latest

# Expose the application port
EXPOSE 8080

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/api /app/api

# Command to run the executable
CMD ["/app/api"]