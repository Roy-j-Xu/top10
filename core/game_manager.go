package core

import (
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
	http.HandleFunc("/api/create-room", handleNewRoom(gm))
	http.HandleFunc("/api/room-info", handleRoomInfo(gm))
	// joining room and establish socket connection
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

func (gm *GameManager) NewRoomSync(roomName string, roomSize int) error {
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
	go func() {
		room.WaitForStartSync()
		game.NewGame(room)
	}()

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

	// block and wait for the next package
	var playerID string
	if err := conn.ReadJSON(&playerID); err != nil {
		log.Println("failed to read player name:", err)
		conn.WriteJSON("failed to read player name")
		conn.Close()
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

// deletes room when its context is done
func (gm *GameManager) watchRoom(r *room.Room) {
	<-r.StopCtx().Done() // wait until room stops
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	delete(gm.Rooms, r.ID)
	log.Printf("room %s removed", r.ID)
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
