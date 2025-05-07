# Backend 

### Build & Run
1. Run `make.sh` (see Makefile extras below)
2. Go to `bin`
3. Start login service : `./user-service user-service.yaml`
4. Start request service: `./request-service request-service.yaml`
5. Start game service: `./game-service` 
~~5. [test] run `./gen-players.sh` to insert 10 dummy players~~ (still possible but not recommended avoid it, use the mongo)


### Todo:
1. ~~Learn uber's zap logger and use it accordingly in the project~~
2. ~~Replace all printlines with `logger.ZapperLog` and use it for everything~~


### Protobuf server-client [deprecated soon]
~~1. Those are just a template usages of grpc~~
~~2. Use as reference not for final product we will need much more~~
~~3. Fixed also make to deal with those 2~~

### Makefile extras
1. Run `./make.sh` with no args to build and update go modules of all projects
2. Run `./make.sh <branchname>` to build and update go modules from `internal-libs` on your branchname from `internal-libs` branch
3. Run `./make.sh -n` to only build projects with no update of go modules of `internal-libs` (local build faster)

### fe-tests
1. `multi-bet.html` creates 10 websocket clients that betting infinitely and also asks in 5 seconds for get players/winners
2. `pay-histogram.html` - histogram of payments , bets infinitely with 100 on 1 id and plots the paytable