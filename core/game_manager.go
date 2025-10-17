package core

import (
	"fmt"
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

func (gm *GameManager) NewRoom(roomName string, roomSize int) (*room.Room, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if _, ok := gm.Rooms[roomName]; ok {
		return nil, fmt.Errorf("creating room %s: %w", roomName, room.ErrRoomExists)
	}

	room, err := room.NewRoomWebSocket(roomName, roomSize)
	if err != nil {
		return nil, err
	}

	gm.Rooms[room.ID] = room

	return room, nil
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
