#!/bin/bash
cd $(dirname $0)

rm bin/gocron
cd ./src
go build
mv src  ../bin/gocron
