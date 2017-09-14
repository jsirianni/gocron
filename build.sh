#!/bin/bash

cd $(dirname $0)

go build src/gocron.go
mv gocron  bin/
