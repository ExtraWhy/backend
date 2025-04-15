package server

import (
	playercache "casino/rest-backend/player-cache"
	server "casino/rest-backend/proto-client"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/ExtraWhy/internal-libs/config"
	"github.com/ExtraWhy/internal-libs/db"
	"github.com/ExtraWhy/internal-libs/models/player"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var allNodesWaitGroup sync.WaitGroup

type Server struct {
	Host    string
	Port    uint16
	router  *gin.Engine
	dbiface db.DbIface
	winReq  server.WinRequest
}

func (srv *Server) DoRun(conf *config.RequestService) error {

	if conf.DatabaseType == "mongo" {
		srv.dbiface = &db.NoSqlConnection{}
		srv.dbiface.Init("Cluster0", "cryptowincryptowin:EfK0weUUe7t99Djx")
		srv.router = gin.Default()

	} else {
		srv.dbiface = &db.DBSqlConnection{}
		srv.dbiface.Init("sqlite3", "players.db")
		srv.router = gin.Default()
		defer srv.dbiface.Deinit()
		srv.dbiface.(*db.DBSqlConnection).CreatePlayersTable()
	}
	srv.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"}, // Next.js frontend
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	srv.router.GET("/players", srv.getPlayers)
	srv.router.GET("/players/:id", srv.getPlayerById)
	srv.router.POST("/players", srv.postPlayers)
	srv.router.GET("/players/winners", srv.getWinners)
	srv.router.GET("/players/:id/bet/:m", srv.getPlayerPlay)
	hp := fmt.Sprintf("%s:%s", conf.RestServiceHost, conf.RestServicePort)
	return srv.router.Run(hp)
}

type Fe_resp struct {
	Won   uint64  `json:"won"`
	Name  string  `json:"name"`
	Lines []uint8 `json:"lines"`
}

func (srv *Server) getPlayerPlay(gct *gin.Context) {
	allNodesWaitGroup.Add(1)
	go func(s *Server, ctx *gin.Context) {
		defer allNodesWaitGroup.Done()
		id := ctx.Param("id")
		bet := ctx.Param("m")
		tmp, _ := strconv.ParseUint(id, 10, 64) //TODO handle error laster
		beti, err := strconv.ParseUint(bet, 10, 64)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error in bet request to bet", "err": err})
			return
		}

		players := s.dbiface.DisplayPlayers()
		if len(players) > 0 {
			for _, i := range players {
				if i.Id == tmp {
					if i.Money < beti {
						ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Not enough cuurencty  to bet"})
						break
					}
					//do proto call
					if err := srv.winReq.SendWin(tmp); err != nil {
						ctx.IndentedJSON(http.StatusBadGateway, gin.H{"message": "Fail to talk to the game service"})
						break
					} else {
						fe := Fe_resp{Name: i.Name,
							Won: srv.winReq.PlayerResponse.GetMoneyWon(),
						}
						if fe.Won > 0 {
							i.Money += (fe.Won * beti)
							tmp2 := srv.winReq.PlayerResponse.GetLines()
							for i := 0; i < len(tmp2); i++ {
								fe.Lines = append(fe.Lines, tmp2[i])
							}
							playercache.PutToCache(&i)
						} else {
							i.Money = i.Money - beti
						}
						if _, err := s.dbiface.UpdatePlayerMoney(&i); err != nil {
							ctx.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "Fails to update player"})
						} else {
							ctx.IndentedJSON(http.StatusOK, fe)
						}
						break
					}
				}
			}
		} else {
			ctx.IndentedJSON(http.StatusNoContent, gin.H{"message": "No players to display"})
		}
	}(srv, gct) //clsr
	allNodesWaitGroup.Wait()
}

func (srv *Server) getPlayers(ctx *gin.Context) {
	p := srv.dbiface.DisplayPlayers()
	if len(p) > 0 {
		ctx.IndentedJSON(http.StatusOK, p)
	} else {
		ctx.IndentedJSON(http.StatusNoContent, gin.H{"message": "No players to display"})
	}

}

func (srv *Server) getWinners(ctx *gin.Context) {

	if playercache.CacheSize() > 0 {
		playercache.DropThem()
		pl := playercache.GetThem()
		fmt.Println(pl)
		ctx.IndentedJSON(http.StatusOK, pl)
		return
	}
	ctx.IndentedJSON(http.StatusNoContent, gin.H{"message": "No 5 winners present"})

}

func (srv *Server) getPlayerById(ctx *gin.Context) {
	id := ctx.Param("id")
	tmp, _ := strconv.ParseUint(id, 10, 64) //TODO handle error laster
	players := srv.dbiface.DisplayPlayers()
	for _, a := range players {
		if a.Id == tmp {
			ctx.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	errmsg := fmt.Sprintf("Player with id %d does not exists", tmp)
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": errmsg})
}

func findById(p *player.Player, players []player.Player) bool {
	for _, id := range players {
		if id.Id == p.Id {
			return true
		}
	}
	return false
}

func (srv *Server) postPlayers(ctx *gin.Context) {
	var komardjia player.Player

	// Call BindJSON to bind the received JSON to
	if err := ctx.BindJSON(&komardjia); err != nil {
		return
	}
	p := srv.dbiface.DisplayPlayers()
	if findById(&komardjia, p) {
		errmsg := fmt.Sprintf("Player with id %d does  exists", komardjia.Id)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": errmsg})
		return
	}
	// Add the new album to the slice.
	srv.dbiface.AddPlayer(&komardjia)
	//	players = append(players, komardjia)
	ctx.IndentedJSON(http.StatusCreated, komardjia)
}
