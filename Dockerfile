FROM alpine:latest AS gocron-base
COPY gocron /usr/local/bin/gocron
RUN \
    chmod +x /usr/local/bin/gocron; \
    adduser -S gocron
USER gocron


FROM gocron-base AS gocron-front
expose 8080
ENTRYPOINT ["/usr/local/bin/gocon", "frontend", "verbose"]


FROM gocron-base AS gocron-back
ENTRYPOINT ["/usr/local/bin/gocron", "backend", "verbose"]
