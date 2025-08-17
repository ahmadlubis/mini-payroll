FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application - Fixed: point to cmd directory, not cmd/main.go file
RUN CGO_ENABLED=0 GOOS=linux go build -o payslip-system ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/payslip-system .

EXPOSE 8080

CMD ["./payslip-system"]