package websocket

import (
	feresponse "casino/rest-backend/models"
	playercache "casino/rest-backend/player-cache"
	server "casino/rest-backend/proto-client"
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
	winReq  server.WinRequest
}

const (
	ok            = 0
	no_money      = 1
	db_err_write  = 2
	db_err_read   = 3
	db_no_players = 4
	unknownw      = 0xff
)

type cachedPlayer struct {
	Hits uint64
	Pl   player.Player
}
type skv struct {
	k uint64
	v cachedPlayer
}

type MessageBet struct {
	Id    uint64 `json:"id"`
	Money uint64 `json:"money"`
}

var clients = make(map[*websocket.Conn]*player.Player)
var broadcast = make(chan MessageBet)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

func (srv *WSServer) DoRun(conf *config.RequestService) error {

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

	//srv.router.Use(cors.New(cors.Config{
	//	AllowOrigins: []string{"http://localhost:3000"}, // Next.js frontend
	//	AllowMethods: []string{"GET", "POST", "OPTIONS"},
	//	AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	//}))

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
	hp := fmt.Sprintf(":%s", conf.RestServicePort)
	return srv.router.Run(hp)
}

func (srv *WSServer) getPlayerPlay(msg *MessageBet, fe *feresponse.Fe_resp) uint {

	players := srv.dbiface.DisplayPlayers()
	fe.Id = msg.Id
	if len(players) > 0 {
		for _, i := range players {
			if i.Id == fe.Id {
				if i.Money < msg.Money {
					return no_money
				}
				//do proto call
				if err := srv.winReq.SendWin(msg.Id); err != nil {
					return unknownw
				} else {
					fe.Won = srv.winReq.PlayerResponse.GetMoneyWon()
					if fe.Won > 0 {
						fe.Name = i.Name
						i.Money += (fe.Won * msg.Money)
						tmp2 := srv.winReq.PlayerResponse.GetLines()
						for i := 0; i < len(tmp2); i++ {
							fe.Lines = append(fe.Lines, tmp2[i])
						}
						playercache.PutToCache(&i)
					} else {
						i.Money = i.Money - msg.Money
					}
					fe.Reels = srv.winReq.PlayerResponse.GetReels()

					if _, err := srv.dbiface.UpdatePlayerMoney(&i); err != nil {
						return db_err_write
					}
					break
				}
			}
		}
	} else {
		return db_no_players
	}
	fmt.Println(fe)
	return ok
}

func (srv *WSServer) getPlayerPlayCleo(msg *MessageBet, fer *feresponse.Fe_resp_slots) uint {

	//do proto call
	if err := srv.winReq.SendWin4Cleo(msg.Id); err != nil {
		return unknownw
	} else {
		var j = 0
		for x := 0; x < 5; x++ {
			for y := 0; y < 3; y++ {
				fer.Scr[x][y] = srv.winReq.CleopatraResponse.Syms[x*3+y]
			}
		}
		fer.Cleo = make([]feresponse.Fe_resp_cleo, len(srv.winReq.CleopatraResponse.Wins))
		for i := range srv.winReq.CleopatraResponse.Wins {

			fer.Cleo[j].XY = make([]uint32, 1)
			fer.Cleo[j].Num = make([]uint32, 1)
			if srv.winReq.CleopatraResponse.Wins[i].BID != nil {
				fer.Cleo[j].BID = *srv.winReq.CleopatraResponse.Wins[i].BID
			}
			if srv.winReq.CleopatraResponse.Wins[i].Free != nil {
				fer.Cleo[j].Free = *srv.winReq.CleopatraResponse.Wins[i].Free
			}
			if srv.winReq.CleopatraResponse.Wins[i].JID != nil {
				fer.Cleo[j].JID = *srv.winReq.CleopatraResponse.Wins[i].JID
			}
			if srv.winReq.CleopatraResponse.Wins[i].Jack != nil {
				fer.Cleo[j].Jack = *srv.winReq.CleopatraResponse.Wins[i].Jack
			}
			if srv.winReq.CleopatraResponse.Wins[i].Line != nil {
				fer.Cleo[j].Line = *srv.winReq.CleopatraResponse.Wins[i].Line
			}
			if srv.winReq.CleopatraResponse.Wins[i].Mult != nil {
				fer.Cleo[j].Mult = *srv.winReq.CleopatraResponse.Wins[i].Mult
			}
			if srv.winReq.CleopatraResponse.Wins[i].Pay != nil {
				fer.Cleo[j].Pay = *srv.winReq.CleopatraResponse.Wins[i].Pay
			}
			if srv.winReq.CleopatraResponse.Wins[i].Sym != nil {
				fer.Cleo[j].Sym = *srv.winReq.CleopatraResponse.Wins[i].Sym
			}
			if srv.winReq.CleopatraResponse.Wins[i].Linex != nil {
				for k := 0; k < len(*&srv.winReq.CleopatraResponse.Wins[i].Linex); k++ {
					fer.Cleo[j].XY = append(fer.Cleo[j].XY, *&srv.winReq.CleopatraResponse.Wins[i].Linex[k])
				}
			}
			j++
		}

		return ok
	}

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
		fmt.Println(pl)
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
		//	fe := feresponse.Fe_resp{}
		fecleo := feresponse.Fe_resp_slots{}
		//		fecleo.Cleo = make([]feresponse.Fe_resp_cleo, 1)
		//test		res := ws.getPlayerPlay(&msg, &fe)
		res := ws.getPlayerPlayCleo(&msg, &fecleo)

		if res == 0 {
			for client, _ := range clients {
				err := client.WriteJSON(fecleo)
				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
			//for client, player := range clients {
			//	if player.Id == fe.Id {
			//		err := client.WriteJSON(fe)
			//		if err != nil {
			//			client.Close()
			//			delete(clients, client)
			//		}
			//	}
			//}
		}
	}
}

func handleWebSocketConnection(conn *websocket.Conn) {
	for {
		var message MessageBet
		err := conn.ReadJSON(&message)
		if err != nil {
			conn.Close()
			delete(clients, conn)
			break
		}
		clients[conn].Id = message.Id
		broadcast <- message
	}
}
