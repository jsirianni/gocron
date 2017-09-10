#!/bin/bash
cd $(dirname $0)

go build src/cronserver.go
mv cronserver bin/
