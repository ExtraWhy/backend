package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections
		return true
	},
}

type Message struct {
	Id     int `json:"id"`
	Action int `json:"bet"`
}

type MessageLogin struct {
	Id int `json:"id"`
}

type Foo struct {
	F int `json:"F"`
	B int `json:"B"`
}
type Resp struct {
	Msg string `json:"response"`
}

var clients = make(map[*websocket.Conn]int)
var broadcast = make(chan Message)

var logins = make(chan MessageLogin)

func main() {

	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		go handleWebSocketConnection(conn)
	})
	go handleBroadcast()
	r.Run(":8083")
}

func handleWebSocketConnection(conn *websocket.Conn) {
	for {
		var message Message
		var login MessageLogin
		err := conn.ReadJSON(&login)
		if err != nil {
			conn.Close()
			delete(clients, conn)
		} else {
			if ok := clients[conn]; ok == 0 {
				clients[conn] = 1
			} else {
				logins <- login
			}
		}
		err = conn.ReadJSON(&message)
		if err != nil {
			conn.Close()
			delete(clients, conn)

		}
		broadcast <- message
	}
}

func handleBroadcast() {
	for {
		<-broadcast
		for client, v := range clients {
			if v == 1 {
				f := Foo{F: 1, B: 2}
				err := client.WriteJSON(f)
				if err != nil {
					client.Close()
					delete(clients, client)

				}
				err = client.WriteJSON(l)
				if err != nil {
					client.Close()
					delete(clients, client)

				}
			}
		}
	}
}
