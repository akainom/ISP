package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/ws", nil)
	if err != nil {
		log.Fatal("[FATAL] dial:", err)
	}
	defer conn.Close()

	log.Println("[INFO] connected to server")

	for i := 1; i <= 5; i++ {
		msg := []byte("hiii haii")
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("[ERROR] write:", err)
			return
		}

		_, response, err := conn.ReadMessage()
		if err != nil {
			log.Println("[ERROR] read:", err)
			return
		}
		log.Printf("[LOG] received: %s", response)

		time.Sleep(1 * time.Second)
	}

	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("[ERROR] close:", err)
	}
	log.Println("[INFO] client shutting down")
}
