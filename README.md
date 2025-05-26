# Backend 

### Build & Run
There are two ways to build and run the services:
1. You can build and run them using `docker` and `docker compose` - this method will set everything up without the need to manually build the services and orcestrate them to run. To build and run the services this way, you need to know the following commands:
    1. Run `docker compose up --build -d` in the project's root folder. The last `-d` flad tells docker to run them in detached mode, which means you can use the same console to test the services, too.
    2. If running the services this way doesn't seem to work, then you can see what the logs for the different services are by running.
     `docker logs -f <container_id or container_name>`.
     You can get the names of the running containers using
     `docker ps`   (if you'd like to see all containers, including the stopped ones, then run `docker ps -a` alternatively).
    3. If you encounter problems that a specific image is taken or being used, just run `docker system prune -a --volumes` to clean docker and follow step `1.`
1. Using the make script
    1. Run `make.sh` (see Makefile extras below)
    2. Go to `bin`
    3. Start login service : `./user-service user-service.yaml`
    4. Start request service: `./request-service request-service.yaml`
    ~~5. Start game service: `./game-service`~~ 
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

### Test for RTP
1. run `./make.sh`
2. go to `bin` folder
3. modify `requests-service.yaml` to use `sqlite3` instead of `mongo` if you are offline, otherwise leave unchanged.
4. run `./request-service request-service.yaml`
5. In case you use `sqlite3` option , run `./gen-players.sh` to populate the local db with dummy players. 
6. ~~run `./game-service`  or `./game-service t` (t opiton os for the permuted reels test), no option is for random reels (original)~~
7. open in the browser `pay-histogram.html` and press `F12` for firefox to see the console log 
8. You may watch the histogram in the canvas and track the rtp percentage in the console log.


### fe-tests
1. `multi-bet.html` creates 10 websocket clients that betting infinitely and also asks in 5 seconds for get players/winners
2. `pay-histogram.html` - histogram of payments , bets infinitely with 100 on 1 id and plots the paytable


### RNG
1. We have to give a good seed per user session (not per game )
2. Good approach is to seed mt19973 with dev/urandom or golangs embeded crypto/rand
3. In real life 2 players will have different seeds so the deterministic output cannot be predicted avoiding bruteforce
