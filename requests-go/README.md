### Backend 

#### install:

1. Install `golang` for your machine (see Win32/X/MacOs) for linux you may use local script `install-me-go.sh`
2. Install `gin package` from [here](https://pkg.go.dev/github.com/gin-gonic/gin#section-readme)


#### Build:
1. In root folder run `go build main.go`
2. Run the binary - defaults:
`localhost:8080`

### Run:
1. When it runs you may open browser to `localhost:8080/players` to see all
2. run in terminal 
`curl http://localhost:8080/players  --include  --header "Content-Type: application/json" --request "POST" --data '{"id" : 6, "money" : 10000, "name" : "Lubakamadafaka"}'` 
to add new player

### Current api

1. GET 	`localhost:8081/players` - get all players 
2. GET 	`localhost:8081/players/<id>` get player by id 
3. GET  `localhost:8081/players/winners` get last winners (todo criteria)
4. POST	`localhost:8081/players` - post a new player (see Run)


### Current player json format 
```
   {
        "id": 1,
        "money": 123456,
        "name": "Lubaka F"
    }
```

### Howto SQL
1. Install `sqlite3`
2. No setup needed 

### Dummy data
1. After `make.sh` go to `bin`
2. Run the service 
3. run `gen-players.sh` to add 10 dummy players 
4. run `sqlite3 players.db` in `bin`
5. in sqlite shell run `select * from players;`
You should see the test data:
```
1|10000|Lubaka
2|10000|Kalniq
3|10000|Gandalf
4|10000|Krasena
5|10000|Ekstramena
6|10000|Shto?
7|10000|Kucheto
8|10000|Bonbonev
9|999999999|Skalata
10|10000|Robota
```
