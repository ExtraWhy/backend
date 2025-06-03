package wsserver

import (
	"casino/game/crwcleopatra"
	"casino/game/models"
	feresponse "casino/game/models"
	playercache "casino/game/player-cache"
	"casino/recovery"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ExtraWhy/internal-libs/config"
	"github.com/ExtraWhy/internal-libs/db"
	"github.com/ExtraWhy/internal-libs/models/player"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	Host    string
	Port    uint16
	router  *gin.Engine
	dbiface db.DbIface
}

const (
	CRW_Ok                   = 0
	CRW_No_money             = 1
	CRW_Db_err_write         = 2
	CRW_Db_err_read          = 4
	CRW_Db_no_players        = 8
	CRW_Db_no_player_with_id = 16
	CRW_No_Win               = 32
	CRW_Unknown              = 0xff
)

type cachedPlayer struct {
	Hits uint64
	Pl   player.Player
}
type skv struct {
	k uint64
	v cachedPlayer
}

var clients = make(map[*websocket.Conn]*player.Player)
var broadcast = make(chan models.MessageBet)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

func testMode(s *string) bool {
	if s != nil && (*s == "on" || *s == "On" || *s == "Enabled" ||
		*s == "enabled" || *s == "True" || *s == "true") {
		return true
	} else {
		return false
	}
}

func (srv *WSServer) DoRun(conf *config.RequestService) error {

	crwcleopatra.TestMode = testMode(&conf.TestMode)

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

	srv.router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		clients[conn] = &player.Player{Id: 0, Name: "", Money: 0}
		go handleWebSocketConnection(conn)
	})
	go srv.handleBroadcast()
	hp := fmt.Sprintf("%s:%s", conf.WsServiceHost, conf.WsServicePort)
	return srv.router.Run(hp)
}

func (srv *WSServer) getPlayerPlayCleo(msg *models.MessageBet, fer *feresponse.CRW_Fe_resp_slots) uint {

	var player *player.Player
	players := srv.dbiface.DisplayPlayers()
	for _, i := range players {
		if i.Id == msg.Id {
			if i.Money < msg.Money {
				return CRW_No_money
			}
			player = &i
			break
		} //
	}

	if player == nil {
		return CRW_Db_no_player_with_id
	}
	resp := crwcleopatra.GetWinForCleopatra(msg)
	var j = 0
	for x := 0; x < 5; x++ {
		for y := 0; y < 3; y++ {
			fer.Scr[x][y] = resp.Syms[x*3+y]
		}
	}

	//todo sym the win response and decide the player to display in the most recent played
	if len(resp.Wins) > 0 {
		playercache.PutToCache(player)
	}
	fer.Cleo = make([]feresponse.CRW_Fe_resp_cleo, len(resp.Wins))
	var freegames uint64 = 0
	for i := range resp.Wins {

		fer.Cleo[j].XY = make([]uint32, 1)
		fer.Cleo[j].BID = resp.Wins[i].BID

		fer.Cleo[j].Free = resp.Wins[i].Free
		freegames += uint64(resp.Wins[i].Free)
		fer.Cleo[j].JID = resp.Wins[i].JID

		fer.Cleo[j].Jack = resp.Wins[i].Jack

		fer.Cleo[j].Line = resp.Wins[i].Line

		fer.Cleo[j].Mult = resp.Wins[i].Mult

		fer.Cleo[j].Pay = resp.Wins[i].Pay

		fer.Cleo[j].Sym = resp.Wins[i].Sym

		fer.Cleo[j].Num = resp.Wins[i].Num

		if resp.Wins[i].Linex != nil {
			for k := 0; k < len(*&resp.Wins[i].Linex); k++ {
				fer.Cleo[j].XY = append(fer.Cleo[j].XY, *&resp.Wins[i].Linex[k])
			}
		}
		j++
	}
	recovery.AddRecord(*player)
	return CRW_Ok

}

func (srv *WSServer) getPlayers(ctx *gin.Context) {
	p := srv.dbiface.DisplayPlayers()
	if len(p) > 0 {
		ctx.IndentedJSON(http.StatusOK, p)
	} else {
		ctx.IndentedJSON(http.StatusNoContent, gin.H{"message": "No players to display"})
	}

}

func (srv *WSServer) getWinners(ctx *gin.Context) {

	if playercache.CacheSize() > 0 {
		playercache.DropThem()
		pl := playercache.GetThem()
		ctx.IndentedJSON(http.StatusOK, pl)
		return
	}
	ctx.IndentedJSON(http.StatusNoContent, gin.H{"message": "No 5 winners present"})
}

func (srv *WSServer) getPlayerById(ctx *gin.Context) {
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

func (srv *WSServer) postPlayers(ctx *gin.Context) {
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

func (ws *WSServer) handleBroadcast() {
	for {
		msg := <-broadcast
		for client, _ := range clients {
			fecleo := feresponse.CRW_Fe_resp_slots{}
			res := ws.getPlayerPlayCleo(&msg, &fecleo)
			if res == CRW_Ok {
				err := client.WriteJSON(fecleo)

				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func handleWebSocketConnection(conn *websocket.Conn) {
	for {
		var message models.MessageBet
		err := conn.ReadJSON(&message)
		if err != nil {
			conn.Close()
			delete(clients, conn)
			break
		}
		if v, ok := clients[conn]; ok {
			v.Id = message.Id
		}
		broadcast <- message
	}
}
