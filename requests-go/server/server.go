package server

import (
	"casino/rest-backend/db"
	"casino/rest-backend/player"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

// todo: temp for testing - remove later
var players = []player.Player{
	{Id: 1, Name: "Lubaka F", Money: 123456},
	{Id: 2, Name: "Lubaka K", Money: 1},
	{Id: 3, Name: "Kucheto", Money: 5},
	{Id: 4, Name: "Kalniq", Money: 5},
	{Id: 5, Name: "Potniq", Money: 5},
	{Id: 6, Name: "Bavniq", Money: 5},
	{Id: 7, Name: "Burziq", Money: 5},
}

// end todo
type Server struct {
	Host       string
	Port       uint16
	router     *gin.Engine
	autocrt    autocert.Manager //member for certificates with Let's encrypt
	sqliteconn db.DBconnection
}

func (srv *Server) SetHostPort(s string, p uint16) {
	srv.Host = s
	srv.Port = p
}

func (srv *Server) GetHostPortStr() string {
	return fmt.Sprintf("%s:%d", srv.Host, srv.Port)
}

func (srv *Server) DoRun() error {
	srv.sqliteconn.Init()
	srv.router = gin.Default()

	defer srv.sqliteconn.Deinit()

	srv.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"}, // Next.js frontend
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	srv.router.GET("/players", srv.getPlayers)
	srv.router.GET("/players/:id", getPlayerById) //todo
	srv.router.POST("/players", srv.postPlayers)
	return srv.router.Run("localhost:8080")
}

// priv
func (srv *Server) getPlayers(ctx *gin.Context) {
	p := srv.sqliteconn.DisplayPlayers()
	ctx.IndentedJSON(http.StatusOK, p)
}

func getPlayerById(ctx *gin.Context) {
	id := ctx.Param("id")
	tmp, _ := strconv.ParseUint(id, 10, 64) //TODO handle error laster

	for _, a := range players {
		if a.Id == tmp {
			ctx.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	errmsg := fmt.Sprintf("Player with id %d does not exists", tmp)
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": errmsg})
}

func findById(p *player.Player) bool {
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

	if findById(&komardjia) {
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
