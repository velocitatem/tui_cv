FROM alpine:latest

RUN apk add --no-cache openssh git bash && \
    mkdir -p /home/cv/.ssh && \
    ssh-keygen -A && \
    adduser -D -s /bin/bash cv && \
    chown -R cv:cv /home/cv && \
    echo 'PermitEmptyPasswords yes' >> /etc/ssh/sshd_config && \
    echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config && \
    echo 'Match User cv\nForceCommand /home/cv/resume-tui' >> /etc/ssh/sshd_config

COPY resume-tui /home/cv/resume-tui
COPY resume.yaml /home/cv/resume.yaml
RUN chmod +x /home/cv/resume-tui && chown cv:cv /home/cv/*

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D"]
