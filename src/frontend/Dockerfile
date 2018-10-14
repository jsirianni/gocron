FROM alpine:latest


RUN wget https://github.com/jsirianni/gocron/releases/download/v3.0.5/gocron-front -O /usr/local/bin/gocron-front
RUN chmod +x /usr/local/bin/gocron-front


RUN adduser -S gocron
USER gocron


CMD ["/usr/local/bin/gocron-front --verbose --port 8080"]
