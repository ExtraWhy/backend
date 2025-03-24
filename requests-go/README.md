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

1. GET 	`localhost:8080/players` - get all players 
2. GET 	`localhost:8080/players/<id>` get player by id 
3. POST	`localhost:8080/players` - post a new player (see Run)

### Current player json format 
```
   {
        "id": 1,
        "money": 123456,
        "name": "Lubaka F"
    }
```