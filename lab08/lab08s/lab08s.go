package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] upgrade failed:", err)
		return
	}
	defer conn.Close()

	log.Println("[INFO] client connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[WARN] unexpected close: %v", err)
			} else {
				log.Println("[INFO] client disconnected")
			}
			break
		}

		log.Printf("[LOG] received: %s", msg)

		response := "from server " + string(msg)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			log.Println("[ERROR] write failed:", err)
			break
		}
	}
}

func main() {
	file, err := os.OpenFile("08_01.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("[FATAL] unable to open log file, aborting")
		os.Exit(1)
	}
	log.SetOutput(file)

	http.HandleFunc("/ws", handleConnections)

	fmt.Println("WS server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
