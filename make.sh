#!/bin/bash

echo "start make casino all"

output="bin"

ls -l | grep bin > /dev/null
result=$?
if [[ $result == 0 ]]; then
	echo "Deleting "$output" folder..."
	rm -rf $output
fi


login_service_name="user-service"
request_service_name="reuests-service"

internal_libs=""
internal_libs_db=""
internal_libs_logger=""
internal_libs_player=""
internal_libs_config=""


go env -w GOPROXY=direct


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
echo "--------------------------------------------------------------------------------"
echo $internal_libs
echo $internal_libs_db
echo $internal_libs_logger
echo $internal_libs_player
echo $internal_libs_config
echo "--------------------------------------------------------------------------------"

mkdir $output

echo "Prepare login service"
cd login-service
go get $internal_libs
go get $internal_libs_db
go get $internal_libs_logger
go get $internal_libs_player
go get $internal_libs_config
go build -o $login_service_name main.go
mv $login_service_name "../"$output
cp .env "../"$output
echo "finished"
cd ..

echo "Preparing requests service"
cd requests-go
go get $internal_libs
go get $internal_libs_db
go get $internal_libs_logger
go get $internal_libs_player
go get $internal_libs_config
go build  -o $request_service_name main.go
cp gen-players.sh "../"$output
mv $request_service_name "../"$output
cp config.yaml "../"$output
echo "finished"

echo "end make casino all :-)"
