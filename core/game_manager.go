package core

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"top10/core/game"
	"top10/core/room"
)

type GameManager struct {
	Rooms map[string]*room.Room
	Games map[string]*game.Game
	mutex sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		Rooms: make(map[string]*room.Room),
		Games: make(map[string]*game.Game),
	}
}

func (gm *GameManager) RunGame(roomName string) error {
	gm.mutex.Lock()
	r, ok := gm.Rooms[roomName]
	if !ok {
		return fmt.Errorf("running game in room %s: %w", roomName, room.ErrRoomNotFound)
	}
	game := game.NewGame(r)
	gm.Games[roomName] = game
	gm.mutex.Unlock()

	game.Start() // game starts here

	gm.mutex.Lock()
	delete(gm.Games, roomName)
	gm.mutex.Unlock()

	log.Printf("game %s removed", roomName)
	return nil
}

func (gm *GameManager) NewRoomSync(roomName string, roomSize int) (room.RoomInfo, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if _, ok := gm.Rooms[roomName]; ok {
		return room.RoomInfo{}, fmt.Errorf("creating room %s: %w", roomName, room.ErrRoomExists)
	}

	rm, err := room.NewRoomWebSocket(roomName, roomSize)
	if err != nil {
		return room.RoomInfo{}, fmt.Errorf("creating room %s: %w", roomName, err)
	}
	gm.Rooms[rm.ID] = rm

	go gm.watchRoom(rm)
	go func() {
		rm.WaitForStartSync()
		gm.RunGame(roomName)
	}()

	return rm.GetRoomInfoUnsafe(), nil
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

func (gm *GameManager) GetGameSync(roomID string) (*game.Game, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	if r, ok := gm.Games[roomID]; ok {
		return r, nil
	} else {
		return nil, errors.New("game not found")
	}
}
