package server

import (
	"casino/rest-backend/player"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// todo: temp for testing - remove later
var players = []player.Player{
	{Id: 1, Name: "Lubaka F"},
	{Id: 2, Name: "Lubaka K"},
	{Id: 3, Name: "Kucheto"},
}

// end todo
type Server struct {
	Host   string
	Port   uint16
	router *gin.Engine
}

func (srv *Server) SetHostPort(s string, p uint16) {
	srv.Host = s
	srv.Port = p
}

func (srv *Server) GetHostPortStr() string {
	return fmt.Sprintf("%s:%d", srv.Host, srv.Port)
}

func (srv *Server) DoRun() error {
	srv.router = gin.Default()
	srv.router.GET("/players", getPlayers)
	srv.router.POST("/players", postPlayers)
	return srv.router.Run("localhost:8080")
}

// priv
func getPlayers(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, players)
}

func postPlayers(ctx *gin.Context) {
	var komardjia player.Player

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := ctx.BindJSON(&komardjia); err != nil {
		fmt.Errorf("Failed to add player\r\n")
		return
	}

	fmt.Println("Ok , added player")
	// Add the new album to the slice.
	players = append(players, komardjia)
	ctx.IndentedJSON(http.StatusCreated, komardjia)
}
