#!/bin/bash

echo "start make casino all"

go env -w GOPROXY=direct

internal_libs=""
internal_libs_db=""
internal_libs_logger=""
internal_libs_player=""
internal_libs_config=""

if [[ -z $1 ]]; then 
echo "Setting default master repo"
	internal_libs="github.com/ExtraWhy/internal-libs"
	internal_libs_db="github.com/ExtraWhy/internal-libs/db"
	internal_libs_logger="github.com/ExtraWhy/internal-libs/logger"
	internal_libs_player="github.com/ExtraWhy/internal-libs/player"
	internal_libs_config="github.com/ExtraWhy/internal-libs/config"
else
echo "Using a branch for all internlas: "$1
	internal_libs="github.com/ExtraWhy/internal-libs@"$1
	internal_libs_db="github.com/ExtraWhy/internal-libs/db@"$1
	internal_libs_logger="github.com/ExtraWhy/internal-libs/logger@"$1
	internal_libs_player="github.com/ExtraWhy/internal-libs/player@"$1
	internal_libs_config="github.com/ExtraWhy/internal-libs/config@"$1
fi

echo $internal_libs
echo $internal_libs_db
echo $internal_libs_logger
echo $internal_libs_player
echo $internal_libs_config


mkdir bin

cd login-service
go get $internal_libs
go get $internal_libs_db
go get $internal_libs_logger
go get $internal_libs_player
go get $internal_libs_config

go build -o login-service main.go
mv login-service ../bin
cp .env ../bin

cd ..

cd requests-go
go get $internal_libs
go get $internal_libs_db
go get $internal_libs_logger
go get $internal_libs_player
go get $internal_libs_config

go build  -o request-service main.go
cp gen-players.sh ../bin
mv request-service ../bin
cp config.yaml ../bin

echo "end make casino all"
