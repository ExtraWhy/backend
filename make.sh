#!/bin/bash

echo "start make casino all"

mkdir bin

cd login-service
go get github.com/ExtraWhy/internal-libs
go get github.com/ExtraWhy/internal-libs/db
go get github.com/ExtraWhy/internal-libs/logger
go get github.com/ExtraWhy/internal-libs/player
go get github.com/ExtraWhy/internal-libs/config

go build -o login-service main.go
mv login-service ../bin
cp .env ../bin

cd ..

cd requests-go
go get github.com/ExtraWhy/internal-libs
go get github.com/ExtraWhy/internal-libs/db
go get github.com/ExtraWhy/internal-libs/logger
go get github.com/ExtraWhy/internal-libs/player
go get github.com/ExtraWhy/internal-libs/config

go build  -o request-service main.go
cp gen-players.sh ../bin
mv request-service ../bin
cp config.yaml ../bin

echo "end make casino all"
