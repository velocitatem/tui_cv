# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
# Build static binary for Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o resume-tui .

# Runtime stage
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/resume-tui /app/resume-tui
COPY resume.yaml /app/resume.yaml
RUN chmod +x /app/resume-tui
EXPOSE 22
CMD ["/app/resume-tui", "serve"]
