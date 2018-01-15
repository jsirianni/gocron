#!/bin/bash

cd $(dirname $0)


rm ./bin/gocro*

env GOOS=linux GOARCH=amd64 go build -v ./src/frontend/gocron-front.go
mv gocron-front  ./bin/

env GOOS=linux GOARCH=amd64 go build -v ./src/backend/gocron-back.go
mv gocron-back  ./bin/
