FROM golang:1.25.3-alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

# Move to working directory /build
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o server 


# Copy binary from build to main folder
# RUN cp /build/main .

# Build a small image
FROM alpine:3.23

# Move to /dist directory as the place for resulting binary folder
WORKDIR /app

# Add non-root user
RUN adduser -D appuser

COPY --from=builder /app/server .
COPY ./database/data.json /app/database/data.json
COPY ./static /app/static

USER appuser

EXPOSE 9000

# Command to run
CMD ["./server"]