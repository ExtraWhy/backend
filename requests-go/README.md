### Backend (rest service)

#### install:

1. Install `golang` for your machine (see Win32/X/MacOs) for linux you may use local script `install-me-go.sh`
~~2. Install `gin package` from [here](https://pkg.go.dev/github.com/gin-gonic/gin#section-readme)~~


#### Build:
1. In root folder run `go build main.go requests-service.yaml`
2. Run the binary - defaults:
`localhost:8081`

#### Run:
1. When it runs you may open browser to `localhost:8080/players` to see all
2. run in terminal 
`curl http://localhost:8081/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 6, "money" : 10000, "name" : "Lubakamadafaka"}'` 
to add new player

#### Current api

1. GET 	`localhost:8081/players` - get all players 
2. GET 	`localhost:8081/players/<id>` get player by id 
3. GET  `localhost:8081/players/winners` get last winners (todo criteria)
4. POST	`localhost:8081/players` - post a new player (see Run)
5. GET	`localhost:8081/players/<id>/play` get a player and an amout won (MOCK)


#### Current player json format 
```
   {
        "id": 1,
        "money": 0,
        "name": "Lubaka F"
    }
```

#### Howto SQL (optional)
1. Install `sqlite3`
2. No setup needed 
3. Edit `requests-service.yaml` -> `database_type: sqlite3`

#### Howto MONGODB
1. Start [here](https://account.mongodb.com/account/login)
2. Follow the guide
3. [cheetsheet](https://www.mongodb.com/developer/products/mongodb/cheat-sheet/)
4. Edit `requests-service.yaml` -> `database_type: mongo`


#### Dummy data
1. After `make.sh` go to `bin`
2. Run the service 
3. run `gen-players.sh` to add 10 dummy players 
4. run `sqlite3 players.db` in `bin`
5. in sqlite shell run `select * from players;`
You should see the test data:
```
1|0|Lubaka
2|0|Kalniq
3|0|Gandalf
4|0|Krasena
5|0|Ekstramena
6|0|Shto?
7|0|Kucheto
8|0|Bonbonev
9|0|Skalata
10|0|Robota
```
6. Optional to test the random gameplay:
- run `proto-player-serv` in `bin` folder copied with `make.sh`
- use `GET players/id/play` to recieve a random win/lose between 0 and 1

