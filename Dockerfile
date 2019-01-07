FROM golang:alpine AS build-env
WORKDIR /src/gocron
ENV GOPATH=/
ADD . /src/gocron
RUN apk add git
RUN go get . && go build -o gocron


FROM alpine:latest
COPY --from=build-env /src/gocron/gocron /usr/local/bin/gocron
COPY example.config.yml /etc/gocron/config.yml
RUN \
    adduser -S gocron && \
    chmod +x /usr/local/bin/gocron && \
    chown gocron /usr/local/bin/gocron && \
    apk add ca-certificates
USER gocron
