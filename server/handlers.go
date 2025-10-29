package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"top10/core"
	"top10/core/room"

	"github.com/gorilla/websocket"
)

var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,32}$`)

func handleNewRoom(gm *core.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name string `json:"roomName"`
			Size int    `json:"roomSize"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		if !validName.MatchString(req.Name) {
			http.Error(w, fmt.Sprintf("invalid room name \"%s\"", req.Name), http.StatusBadRequest)
			return
		}
		if req.Size <= 0 || req.Size >= 21 {
			http.Error(w, fmt.Sprintf("invalid room size \"%d\"", req.Size), http.StatusBadRequest)
			return
		}

		err := gm.NewRoomSync(req.Name, req.Size)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not create room: %v", err), http.StatusBadRequest)
			return
		}

		writeJson(w, RoomInfoResponse{
			RoomName: req.Name,
			RoomSize: req.Size,
			Game:     "Top10",
			Players:  []string{},
		}, 200)
	}
}

func handleRoomInfo(gm *core.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		roomName := r.URL.Query().Get("roomName")

		rm, err := gm.GetRoomSync(roomName)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to find room \"%s\": %v", roomName, err), http.StatusBadRequest)
			return
		}

		writeJson(w, RoomInfoResponse{
			RoomName: rm.ID,
			RoomSize: rm.MaxSize,
			Game:     "Top10",
			Players:  rm.GetAllPlayerIDsSync(),
		}, 200)
	}
}

func joinHandler(gm *core.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse room name from query or headers
		roomName := r.URL.Query().Get("roomName")
		playerID := r.URL.Query().Get("playerName")

		if !validName.MatchString(playerID) {
			http.Error(w, fmt.Sprintf("invalid player name \"%s\"", playerID), http.StatusBadRequest)
			return
		}

		rm, err := gm.GetRoomSync(roomName)
		if err != nil {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}

		if rm.PlayerExistsAndLeftSync(playerID) {
			err = rm.RejoinPlayerSync(playerID, conn)
		} else {
			err = rm.AddPlayerSync(playerID, conn)
		}

		if err != nil {
			log.Println("failed to join:", err)
			conn.WriteJSON(fmt.Sprint("failed to join", err.Error()))
			conn.Close()
			return
		}

		go handlePlayerMessages(rm, playerID)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

func handlePlayerMessages(r *room.Room, playerID string) {
	player, err := r.GetPlayerSync(playerID)
	if err != nil {
		log.Printf("unable to listen for message from player %s: %s", playerID, err.Error())
		return
	}
	log.Printf("listening for messages from player %s", playerID)
	defer func() {
		r.SendToReadyChannel_LEFT(playerID)
		player.Conn.Close()
	}()

	for {
		var msg room.Message
		if err := player.Conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			return
		}

		if msg.Type == string(room.SP_READY) || msg.Type == string(room.SP_LEFT) {
			r.SendToReadyChannel(playerID, msg)
		} else {
			r.SendToPlayerChannel(playerID, msg)
		}
	}
}
