#!/bin/bash

cd $(dirname $0)


rm ./bin/gocro*

env GOOS=linux GOARCH=amd64 go build -o gocron-front -v ./src/frontend/
mv gocron-front  ./bin/

env GOOS=linux GOARCH=amd64 go build -o gocron-back -v ./src/backend/
mv gocron-back  ./bin/
