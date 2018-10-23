FROM alpine:latest
COPY gocron /usr/local/bin/gocron
RUN \
    chmod +x /usr/local/bin/gocron; \
    adduser -S gocron
USER gocron
