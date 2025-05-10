package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/ExtraWhy/internal-libs/mt19937"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")
	var arr = []uint64{0x123, 0x234, 0x345, 0x456, 0xdeadbeef, 0xdecafbad}
	mt19937.Init_by_array64(arr, uint64(len(arr)))
	//mt19937.Init_genrand64(19650218)
	for {
		// Simulate a number to send every 500ms
		number := mt19937.Genrand64_int64() % math.MaxInt32

		err := conn.WriteJSON(map[string]int{"value": int(number)})
		if err != nil {
			log.Println("Write error:", err)
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/ws", serveWs)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
