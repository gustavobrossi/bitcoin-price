# Stage 1: Build the application
# Start from the official Go image to create a build artifact
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download all dependencies
RUN go mod download

# Build the application to the binary named 'main'
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Build a minimal image from scratch
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy required libraries from the builder stage
COPY --from=builder /lib/x86_64-linux-gnu/libresolv.so.2 /lib/x86_64-linux-gnu/libresolv.so.2
COPY --from=builder /lib/x86_64-linux-gnu/libc.so.6 /lib/x86_64-linux-gnu/libc.so.6
COPY --from=builder /lib64/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2

# Copy CA certificates from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run the binary
CMD ["./main"]
