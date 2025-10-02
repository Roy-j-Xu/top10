package socket

import (
	"log"
	"net/http"
	"top10/core"

	"github.com/gorilla/websocket"
)

func InitSocketHandler(room *core.Room) {
	room.Messagers = append(room.Messagers, &WebSocketMessager{Room: room})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(room, w, r)
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

func wsHandler(room *core.Room, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	player := &core.Player{Conn: conn}

	room.AddPlayer(player)

	go handlePlayerMessages(room, player)
}

func handlePlayerMessages(room *core.Room, player *core.Player) {
	defer player.Conn.Close()

	for {
		var msg map[string]string
		if err := player.Conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			return
		}

		room.ReadyPlayer(player.ID)
	}
}
