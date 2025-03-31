package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ExtraWhy/internal-libs/config"
	"github.com/ExtraWhy/internal-libs/db"
	"github.com/ExtraWhy/internal-libs/player"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

const (
	LAST_PLAYERS = 5
)

type Server struct {
	Host       string
	Port       uint16
	Config     *config.RequestService
	router     *gin.Engine
	autocrt    autocert.Manager //member for certificates with Let's encrypt
	sqliteconn db.DBconnection
}

func (srv *Server) DoRun() error {
	srv.sqliteconn.Init("players.db")
	srv.router = gin.Default()

	defer srv.sqliteconn.Deinit()

	srv.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"}, // Next.js frontend
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	srv.router.GET("/players", srv.getPlayers)
	srv.router.GET("/players/:id", srv.getPlayerById)
	srv.router.POST("/players", srv.postPlayers)
	srv.router.GET("/players/winners", srv.getWinners)

	hp := fmt.Sprintf("%s:%s", srv.Config.RestServiceHost, srv.Config.RestServicePort)
	return srv.router.Run(hp)
}

// priv
func (srv *Server) getPlayers(ctx *gin.Context) {
	p := srv.sqliteconn.DisplayPlayers()
	ctx.IndentedJSON(http.StatusOK, p)

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

// do not use yet, will nbe needed for certificates for Let's encrypt
func (srv *Server) autoCertManagement() {
	srv.autocrt = autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example.com", "example1.com"), // etc - todo for later
		Cache:      autocert.DirCache("/var/www/.cache"),
	}
	log.Fatal(autotls.RunWithManager(srv.router, &srv.autocrt))
}

//todo - add the context and run with golang context also handle sigint / sigterm
// for now will start with cert manager
/*
  // Create context that listens for the interrupt signal from the OS.
  ctx, stop := signal.NotifyContext(
    context.Background(),
    syscall.SIGINT,
    syscall.SIGTERM,
  )
  defer stop()

  r := gin.Default()

  // Ping handler
  r.GET("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "pong")
  })

  log.Fatal(autotls.RunWithContext(ctx, r, "example1.com", "example2.com"))
}
*/
