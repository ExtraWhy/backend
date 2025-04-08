package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Game struct {
	reels [][]uint16
}

var (
	g = Game{reels: [][]uint16{{1, 2, 3}}}
)

type Player struct {
	ID string `json:"id"`
}

type Message struct {
	Method string `json:"action"`
	ID     string `json:"id,omitempty"`
}

var (
	players   = make(map[string]string)
	playersMu sync.RWMutex
	upgrader  = websocket.Upgrader{}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // allow all origins
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		switch message.Method {
		case "POST":
			playersMu.Lock()
			players[message.ID] = message.ID
			playersMu.Unlock()
			resp := map[string]string{"status": "Player saved"}
			conn.WriteJSON(resp)
		case "GET":
			playersMu.RLock()
			player, found := players[message.ID]
			playersMu.RUnlock()

			if found {
				conn.WriteJSON(player)
			} else {
				resp := map[string]string{"error": "Player not found"}
				conn.WriteJSON(resp)
			}
		default:
			resp := map[string]string{"error": "Invalid method"}
			conn.WriteJSON(resp)
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
