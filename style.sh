#!/bin/bash

./parse-sass.sh
go generate assets/*.go
go run -mod=vendor main.go -v -f log.txt
