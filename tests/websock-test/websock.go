package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var (
	upgrader    = websocket.Upgrader{}
	connections sync.Map
	ctx         = context.Background()
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)

type JackpotMessage struct {
	Jackpot float64 `json:"jackpot"`
}

func broadcast(message JackpotMessage) {
	data, _ := json.Marshal(message)

	connections.Range(func(key, value any) bool {
		conn := key.(*websocket.Conn)

		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Failed to send: %v", err)
			conn.Close()
			connections.Delete(conn)
		}
		return true
	})
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // Allow all for demo
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	connections.Store(conn, true)

	// Send current jackpot on connect
	if val, err := redisClient.Get(ctx, "jackpot").Result(); err == nil {
		amt, _ := strconv.ParseFloat(val, 64)
		conn.WriteJSON(JackpotMessage{Jackpot: amt})
	}
}

// Listen for jackpot updates via Redis Pub/Sub
func startRedisSubscriber() {
	sub := redisClient.Subscribe(ctx, "jackpot:updates")
	ch := sub.Channel()

	for msg := range ch {
		var jackpot JackpotMessage
		if err := json.Unmarshal([]byte(msg.Payload), &jackpot); err == nil {
			broadcast(jackpot)
		}
	}
}

// Simulate jackpot updates (every 1s)
func simulateJackpot() {
	for {
		time.Sleep(1 * time.Second)
		newAmount := rand.Float64() * 100000
		jackpot := JackpotMessage{Jackpot: newAmount}

		data, _ := json.Marshal(jackpot)
		redisClient.Set(ctx, "jackpot", fmt.Sprintf("%.2f", newAmount), 0)
		redisClient.Publish(ctx, "jackpot:updates", data)
	}
}

func main() {
	go startRedisSubscriber()
	go simulateJackpot()

	http.HandleFunc("/ws", handleWS)
	fmt.Println("WebSocket server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
