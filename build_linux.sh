#!/bin/bash
cd $(dirname $0)

rm bin/gocron
cd ./src
env GOOS=linux GOARCH=amd64 go build
mv src  ../bin/gocron
