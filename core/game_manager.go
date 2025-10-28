package core

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"top10/core/game"
	"top10/core/room"
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

func (gm *GameManager) HandleHTTP() {
	http.HandleFunc("/api/create-room", handleNewRoom(gm))
	http.HandleFunc("/api/room-info", handleRoomInfo(gm))
	// joining room and establish socket connection
	http.HandleFunc("/ws", joinHandler(gm))
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
