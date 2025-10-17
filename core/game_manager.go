package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"top10/core/game"
	"top10/core/room"

	"github.com/gorilla/websocket"
)

type GameManager struct {
	Rooms map[string]*room.Room
	mutex sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		Rooms: make(map[string]*room.Room),
	}
}

var validRoomName = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,32}$`)

func (gm *GameManager) HandleHTTP() {
	http.HandleFunc("/api/create-room", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name string `json:"name"`
			Size int    `json:"size"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		err := gm.NewRoom(req.Name, req.Size)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not create room: %v", err), http.StatusBadRequest)
			return
		}
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Parse room name from query or headers
		roomName := r.URL.Query().Get("room")
		rm, err := gm.GetRoomSync(roomName)
		if err != nil {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}
		wsHandler(rm, w, r)
	})
}

func (gm *GameManager) RunGame(roomName string) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	r, ok := gm.Rooms[roomName]
	if !ok {
		return fmt.Errorf("running game in room %s: %w", roomName, room.ErrRoomNotFound)
	}
	if r.InGame {
		return fmt.Errorf("running game in room %s: %w", roomName, room.ErrInvalidRoom)
	}
	game := game.NewGame(r)
	go game.Start()
	return nil
}

func (gm *GameManager) NewRoom(roomName string, roomSize int) error {
	if !validRoomName.MatchString(roomName) {
		return fmt.Errorf("creating room %s: invalid name", roomName)
	}
	if roomSize <= 0 {
		return fmt.Errorf("creating room %s: invalid room size", roomName)
	}

	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if _, ok := gm.Rooms[roomName]; ok {
		return fmt.Errorf("creating room %s: %w", roomName, room.ErrRoomExists)
	}

	room, err := room.NewRoomWebSocket(roomName, roomSize)
	if err != nil {
		return fmt.Errorf("creating room %s: %w", roomName, err)
	}
	gm.Rooms[room.ID] = room

	go gm.watchRoom(room)

	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

func wsHandler(rm *room.Room, w http.ResponseWriter, r *http.Request) {
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

	rm.AddPlayerSync(playerID, conn)

	go handlePlayerMessages(rm, playerID)
}

func handlePlayerMessages(r *room.Room, playerID string) {
	player, err := r.GetPlayerSync(playerID)
	if err != nil {
		log.Printf("unable to listen for message from player %s: %s", playerID, err.Error())
		return
	}
	log.Printf("listening for messages from player %s", player.ID)
	defer func() {
		r.SendToReadyChannel_LEFT(player.ID)
		player.Conn.Close()
	}()

	for {
		var msg room.Message
		if err := player.Conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			return
		}

		if msg.Type == string(room.SP_READY) || msg.Type == string(room.SP_LEFT) {
			r.SendToReadyChannel(player.ID, msg)
		} else {
			r.SendToPlayerChannel(player.ID, msg)
		}
	}
}

// deletes room when its context is done
func (gm *GameManager) watchRoom(r *room.Room) {
	go func() {
		<-r.StopCtx().Done() // wait until room stops
		gm.mutex.Lock()
		defer gm.mutex.Unlock()
		delete(gm.Rooms, r.ID)
		log.Printf("room %s removed from manager", r.ID)
	}()
}

func (gm *GameManager) GetRoomSync(roomID string) (*room.Room, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	if r, ok := gm.Rooms[roomID]; ok {
		return r, nil
	} else {
		return nil, room.ErrRoomNotFound
	}
}
