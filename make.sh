#!/bin/bash

echo "start make casino all"

output="bin"

ls -l | grep $output > /dev/null
result=$?
if [[ $result == 0 ]]; then
	echo "Deleting "$output" folder..."
	rm -rf $output
fi


login_service_name="user-service"
request_service_name="requests-service"
proto_service="game-service"


internal_libs=""
internal_libs_db=""
internal_libs_logger=""
internal_libs_player=""
internal_libs_user=""
internal_libs_games=""
internal_libs_config=""
internal_libs_proto_models=""
arg1=

function defaults() {
	echo "Setting default master repo"
	internal_libs="github.com/ExtraWhy/internal-libs"
	internal_libs_db="github.com/ExtraWhy/internal-libs/db"
	internal_libs_logger="github.com/ExtraWhy/internal-libs/logger"
	internal_libs_player="github.com/ExtraWhy/internal-libs/models/player"
	internal_libs_games="github.com/ExtraWhy/internal-libs/models/games"
	internal_libs_user="github.com/ExtraWhy/internal-libs/models/user"
	internal_libs_config="github.com/ExtraWhy/internal-libs/config"
	internal_libs_proto_models="github.com/ExtraWhy/internal-libs/proto-models"
}

function update_go() {
if [[ $arg1 == "-n" ]]; then
	echo "Updating go modules off"
else
	echo "Updating go modules on"
	go get $internal_libs
	go get $internal_libs_db
	go get $internal_libs_logger
	go get $internal_libs_player
	go get $internal_libs_user
	go get $internal_libs_config
	go get $internal_libs_proto_models
	go mod tidy	
fi
	
}

go env -w GOPROXY=direct

if [[ -z $1 ]]; then
	defaults
else
	arg1=$1
	if [[ $arg1 == "-n" ]]; then 
		echo "No update modules"
		defaults
	else
		echo "Using a branch for all internlas: "$1
		internal_libs="github.com/ExtraWhy/internal-libs@"$1
		internal_libs_db="github.com/ExtraWhy/internal-libs/db@"$1
		internal_libs_logger="github.com/ExtraWhy/internal-libs/logger@"$1
		internal_libs_player="github.com/ExtraWhy/internal-libs/models/player@"$1
		internal_libs_user="github.com/ExtraWhy/internal-libs/models/user@"$1
		internal_libs_games="github.com/ExtraWhy/internal-libs/models/games@"$1
		internal_libs_config="github.com/ExtraWhy/internal-libs/config@"$1
		internal_libs_proto_models="github.com/ExtraWhy/internal-libs/proto-models@"$1
	fi
	
fi
echo "--------------------------------------------------------------------------------"
echo $internal_libs
echo $internal_libs_db
echo $internal_libs_logger
echo $internal_libs_player
echo $internal_libs_user
echo $internal_libs_config
echo $internal_libs_games
echo $internal_libs_proto_models
echo "--------------------------------------------------------------------------------"

mkdir $output

echo "Prepare game service "
cd game-service
update_go

go build -o $proto_service main.go
mv $proto_service "../"$output
cp *.sh "../"$output
echo "finished"
cd ..


echo "Prepare user service"
cd user-service
update_go

go build -o $login_service_name main.go
mv $login_service_name "../"$output
cp *.yaml "../"$output
echo "finished"
cd ..

echo "Preparing requests service"
cd requests-service
update_go

go build  -o $request_service_name main.go
cp gen-players.sh "../"$output
mv $request_service_name "../"$output
cp *.yaml "../"$output
cp *.html "../"$output
echo "finished"

echo "end make casino all :-)"
