# Build stage
FROM golang:1.16.2-alpine AS builder

# Install gcc and related build utilities
RUN apk add build-base

# Change working directory
WORKDIR /app

# Copy go dependency files
COPY go.* ./

# Download dependencies
RUN go mod download

# Copy rest of the source code
COPY . ./

# Build the executable
RUN go build -o server

# Final stage
FROM alpine

# Change working directory
WORKDIR /app

# Copy built executable from build stage
COPY --from=builder /app/server /app/

# Set the entrypoint as the executable
ENTRYPOINT ./server