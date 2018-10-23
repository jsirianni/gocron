#!/bin/bash
cd $(dirname $0)

env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build && sudo docker build -t gocron:latest .
