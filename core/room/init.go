package room

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func InitSocketHandler(room *Room) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(room, w, r)
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

func wsHandler(room *Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	var playerID string
	if err := conn.ReadJSON(&playerID); err != nil {
		log.Println("failed to read player name:", err)
		conn.WriteJSON("failed to read player name")
		conn.Close()
		return
	}

	room.AddPlayerSync(playerID, conn)

	go handlePlayerMessages(room, playerID)
}

func handlePlayerMessages(room *Room, playerID string) {
	player := room.Players[playerID]
	log.Printf("Listening for messages from player %s", player.ID)
	defer func() {
		room.SendToReadyChannel_LEFT(player.ID)
		player.Conn.Close()
	}()

	for {
		var msg Message
		if err := player.Conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			continue
		}

		if msg.Type == string(SP_READY) || msg.Type == string(SP_LEFT) {
			room.SendToReadyChannel(player.ID, msg)
		} else {
			room.SendToPlayerChannel(player.ID, msg)
		}
	}
}
