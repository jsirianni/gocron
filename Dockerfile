FROM golang:1.12.1-alpine3.9 AS build-env
WORKDIR /src/gocron
ENV GOPATH=/
ADD . /src/gocron
RUN apk add git
RUN go get .
RUN env CGO_ENABLED=0 go test ./...
RUN env CGO_ENABLED=0 go build -o gocron


FROM alpine:latest
COPY --from=build-env /src/gocron/gocron /usr/local/bin/gocron
RUN \
    adduser -S gocron && \
    chmod +x /usr/local/bin/gocron && \
    chown gocron /usr/local/bin/gocron && \
    apk add ca-certificates
USER gocron
