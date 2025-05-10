package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
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
	max := new(big.Int).SetUint64(^uint64(0))
	// we will seed once with dev/randum crypto seed then will continue with mt
	// assume each user session gets it's own seed it\s almst impossible or only via a bug : )
	// to predict the outcome of the rng
	//
	el1, _ := rand.Int(rand.Reader, max)
	el2, _ := rand.Int(rand.Reader, max)
	el3, _ := rand.Int(rand.Reader, max)
	el4, _ := rand.Int(rand.Reader, max)
	el5, _ := rand.Int(rand.Reader, max)
	el6, _ := rand.Int(rand.Reader, max)
	el7, _ := rand.Int(rand.Reader, max)
	el8, _ := rand.Int(rand.Reader, max)

	var arr = []uint64{el1.Uint64(), el2.Uint64(), el3.Uint64(), el4.Uint64(),
		el5.Uint64(), el6.Uint64(), el7.Uint64(), el8.Uint64()}
	mt19937.Init_by_array64(arr, uint64(len(arr)))
	//mt19937.Init_genrand64(19650218)
	for {
		// Simulate a number to send every 500ms
		number := mt19937.Genrand64_int64() % 32

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
