# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
# Build static binary for Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o resume-tui .

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache openssh git bash && \
    ssh-keygen -A && \
    adduser -D -s /bin/bash cv && \
    mkdir -p /home/cv/.ssh && \
    chown -R cv:cv /home/cv
# SSH config for cv user
RUN echo 'PermitEmptyPasswords yes' >> /etc/ssh/sshd_config && \
    echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config && \
    echo 'Match User cv' >> /etc/ssh/sshd_config && \
    echo '    ForceCommand /home/cv/resume-tui' >> /etc/ssh/sshd_config
COPY --from=builder /app/resume-tui /home/cv/resume-tui
COPY resume.yaml /home/cv/resume.yaml
RUN chmod +x /home/cv/resume-tui && chown -R cv:cv /home/cv/resume-tui /home/cv/resume.yaml
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
