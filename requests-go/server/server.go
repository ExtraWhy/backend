package server

import (
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

const (
	LAST_PLAYERS = 5
)

var allNodesWaitGroup sync.WaitGroup

type Server struct {
	Host       string
	Port       uint16
	Config     *config.RequestService
	router     *gin.Engine
	sqliteconn db.DBConnection
	winReq     server.WinRequest
}

func (srv *Server) DoRun() error {
	srv.sqliteconn.Init("sqlite3", "players.db")
	srv.router = gin.Default()
	defer srv.sqliteconn.Deinit()
	srv.sqliteconn.CreatePlayersTable()
	srv.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"}, // Next.js frontend
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	srv.router.GET("/players", srv.getPlayers)
	srv.router.GET("/players/:id", srv.getPlayerById)
	srv.router.POST("/players", srv.postPlayers)
	srv.router.GET("/players/winners", srv.getWinners)
	srv.router.GET("/players/:id/play", srv.getPlayerPlay)
	hp := fmt.Sprintf("%s:%s", srv.Config.RestServiceHost, srv.Config.RestServicePort)
	return srv.router.Run(hp)
}

// priv

func (srv *Server) getPlayerPlay(gct *gin.Context) {
	allNodesWaitGroup.Add(1)
	go func(s *Server, ctx *gin.Context) {
		defer allNodesWaitGroup.Done()
		id := ctx.Param("id")
		tmp, _ := strconv.ParseUint(id, 10, 64) //TODO handle error laster
		players := s.sqliteconn.DisplayPlayers()
		if len(players) > 0 {
			for _, i := range players {
				if i.Id == tmp {
					winner := player.Player{Id: i.Id, Name: i.Name}
					//do proto call
					if err := srv.winReq.SendWin(tmp); err != nil {
						ctx.IndentedJSON(http.StatusBadGateway, gin.H{"message": "Fail to talk to the game service"})
						break
					} else {
						winner.Money = srv.winReq.PlayerResponse.GetMoneyWon()
						fmt.Printf("Player %s won %d \r\n", winner.Name, winner.Money)
						ctx.IndentedJSON(http.StatusOK, winner)
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
	p := srv.sqliteconn.DisplayPlayers()
	if len(p) > 0 {
		ctx.IndentedJSON(http.StatusOK, p)
	} else {
		ctx.IndentedJSON(http.StatusNoContent, gin.H{"message": "No players to display"})
	}

}

func (srv *Server) getWinners(ctx *gin.Context) {
	p := srv.sqliteconn.DisplayPlayers()
	if len(p) >= LAST_PLAYERS {
		ctx.IndentedJSON(http.StatusOK, p[len(p)-LAST_PLAYERS:])
	} else {
		ctx.IndentedJSON(http.StatusNoContent, gin.H{"message": "No content for winners"})
	}

}

func (srv *Server) getPlayerById(ctx *gin.Context) {
	id := ctx.Param("id")
	tmp, _ := strconv.ParseUint(id, 10, 64) //TODO handle error laster
	players := srv.sqliteconn.DisplayPlayers()
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
	p := srv.sqliteconn.DisplayPlayers()
	if findById(&komardjia, p) {
		errmsg := fmt.Sprintf("Player with id %d does  exists", komardjia.Id)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": errmsg})
		return
	}
	// Add the new album to the slice.
	srv.sqliteconn.AddPlayer(&komardjia)
	//	players = append(players, komardjia)
	ctx.IndentedJSON(http.StatusCreated, komardjia)
}
