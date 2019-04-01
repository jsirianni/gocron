# build environment runs unit tests and compiles
# gocron with the official golang image
FROM golang:1.12.1-alpine3.9 AS build-env
WORKDIR /src/gocron
ENV GOPATH=/
ADD . /src/gocron
RUN apk add git
RUN go get .
RUN env CGO_ENABLED=0 go test ./...
RUN env CGO_ENABLED=0 go build -o gocron


# test environment installs postgresql and tests
# gocron at runtime
FROM ubuntu:bionic AS test-env
RUN apt-get update && apt-get install -y iproute2 postgresql curl sudo
COPY --from=build-env /src/gocron/gocron /usr/local/bin/gocron
RUN chmod +x /usr/local/bin/gocron
COPY docker/build_tests.sh /build_test.sh
RUN chmod +x /build_test.sh
RUN bash /build_test.sh


# final artifact contains the gocron binary
#
FROM alpine:latest
COPY --from=build-env /src/gocron/gocron /usr/local/bin/gocron
RUN \
    adduser -S gocron && \
    chmod +x /usr/local/bin/gocron && \
    chown gocron /usr/local/bin/gocron && \
    apk add ca-certificates
USER gocron
