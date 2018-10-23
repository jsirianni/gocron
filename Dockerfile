FROM alpine:latest
COPY gocron /usr/local/bin/gocron
COPY example.config.yml /etc/gocron/config.yml
RUN \
    chmod +x /usr/local/bin/gocron; \
    adduser -S gocron
USER gocron
