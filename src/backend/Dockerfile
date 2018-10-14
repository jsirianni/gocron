FROM alpine:latest


RUN wget https://github.com/jsirianni/gocron/releases/download/v3.0.5/gocron-back -O /usr/local/bin/gocron-back
RUN chmod +x /usr/local/bin/gocron-back


RUN adduser -S gocron
USER gocron


CMD ["/usr/local/bin/gocron-front --verbose]
