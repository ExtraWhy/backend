#!/bin/bash

echo "start make casino all"

mkdir bin

cd login-service
go build -o login-service main.go
mv login-service ../bin
cp .env ../bin

cd ..

cd requests-go
go build  -o request-service main.go
mv request-service ../bin
cp config.yaml ../bin

echo "end make casino all"
