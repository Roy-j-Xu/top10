package room

import (
	"fmt"
	"sync"
)

type GameManager struct {
	Rooms map[string]*Room
	mutex sync.Mutex
}

func (gm *GameManager) NewRoom(roomName string, roomSize int) (*Room, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if _, ok := gm.Rooms[roomName]; ok {
		return nil, fmt.Errorf("creating room %s: %w", roomName, ErrRoomExists)
	}

	room, err := NewRoomWebSocket(roomName, roomSize)
	if err != nil {
		return nil, err
	}

	gm.Rooms[room.ID] = room

	return room, nil
}

func (gm *GameManager) GetRoomSync(roomID string) (*Room, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	if r, ok := gm.Rooms[roomID]; ok {
		return r, nil
	} else {
		return nil, ErrRoomNotFound
	}
}
