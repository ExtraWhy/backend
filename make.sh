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
proto_service="proto-player-service"
proto_client="proto-player-client"


internal_libs=""
internal_libs_db=""
internal_libs_logger=""
internal_libs_player=""
internal_libs_config=""
internal_libs_proto_models=""


go env -w GOPROXY=direct


if [[ -z $1 ]]; then 
echo "Setting default master repo"
	internal_libs="github.com/ExtraWhy/internal-libs"
	internal_libs_db="github.com/ExtraWhy/internal-libs/db"
	internal_libs_logger="github.com/ExtraWhy/internal-libs/logger"
	internal_libs_player="github.com/ExtraWhy/internal-libs/player"
	internal_libs_config="github.com/ExtraWhy/internal-libs/config"
	internal_libs_proto_models="github.com/ExtraWhy/internal-libs/proto-models"

else
echo "Using a branch for all internlas: "$1
	internal_libs="github.com/ExtraWhy/internal-libs@"$1
	internal_libs_db="github.com/ExtraWhy/internal-libs/db@"$1
	internal_libs_logger="github.com/ExtraWhy/internal-libs/logger@"$1
	internal_libs_player="github.com/ExtraWhy/internal-libs/player@"$1
	internal_libs_config="github.com/ExtraWhy/internal-libs/config@"$1
	internal_libs_proto_models="github.com/ExtraWhy/internal-libs/proto-models@"$1
	
fi
echo "--------------------------------------------------------------------------------"
echo $internal_libs
echo $internal_libs_db
echo $internal_libs_logger
echo $internal_libs_player
echo $internal_libs_config
echo $internal_libs_proto_models
echo "--------------------------------------------------------------------------------"

mkdir $output

echo "Prepare proto service "
cd proto-player-serv
go get $internal_libs
go get $internal_libs_db
go get $internal_libs_logger
go get $internal_libs_player
go get $internal_libs_config
go get $internal_libs_proto_models
go mod tidy

go build -o $proto_service main.go
mv $proto_service "../"$output
echo "finished"
cd ..



echo "Prepare proto client "
cd proto-player-client
go get $internal_libs
go get $internal_libs_db
go get $internal_libs_logger
go get $internal_libs_player
go get $internal_libs_config
go get $internal_libs_proto_models
go mod tidy

go build -o $proto_client main.go
mv $proto_client "../"$output
echo "finished"
cd ..



echo "Prepare login service"
cd login-service
go mod tidy
go get $internal_libs
go get $internal_libs_db
go get $internal_libs_logger
go get $internal_libs_player
go get $internal_libs_config
go get $internal_libs_proto_models
go mod tidy

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
go get $internal_libs_proto_models
go mod tidy

go build  -o $request_service_name main.go
cp gen-players.sh "../"$output
mv $request_service_name "../"$output
cp config.yaml "../"$output
echo "finished"

echo "end make casino all :-)"
