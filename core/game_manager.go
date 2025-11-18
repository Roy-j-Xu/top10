package core

import (
	"fmt"
	"log"
	"sync"
	"top10/core/room"

	"github.com/gorilla/websocket"
)

type GameManager struct {
	gameRooms map[string]GameRoom
	mutex     sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		gameRooms: make(map[string]GameRoom),
	}
}

func (gm *GameManager) runGame(roomID string, gameName string) error {
	gm.mutex.Lock()
	gr, err := gm.GetGameRoomUnsafe(roomID)
	if err != nil {
		return fmt.Errorf("running game in room %s: %w", roomID, room.ErrRoomNotFound)
	}
	gm.mutex.Unlock()

	gr.runGame(gameName)

	gm.deleteGameRoom(roomID)

	log.Printf("game room %s removed", roomID)
	return nil
}

func (gm *GameManager) NewRoomSync(roomID string, maxSize int) (room.RoomInfo, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if _, ok := gm.gameRooms[roomID]; ok {
		return room.RoomInfo{}, fmt.Errorf("creating room %s: %w", roomID, room.ErrRoomExists)
	}

	gr, err := newGameRoom(roomID, maxSize)
	if err != nil {
		return room.RoomInfo{}, fmt.Errorf("creating room %s: %w", roomID, err)
	}
	gm.gameRooms[roomID] = gr

	go gm.watchGameRoom(gr)
	go func() {
		gr.Room.WaitForStartSync()
		gm.runGame(roomID, "top10")
	}()

	return gr.Room.GetRoomInfoUnsafe(), nil
}

// deletes game room when its context is done
func (gm *GameManager) watchGameRoom(gr GameRoom) {
	<-gr.Room.StopCtx().Done() // wait until room stops
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	delete(gm.gameRooms, gr.Room.ID)
	log.Printf("room %s removed", gr.Room.ID)
}

func (gm *GameManager) deleteGameRoom(roomID string) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	gr, err := gm.GetGameRoomUnsafe(roomID)
	if err != nil {
		return err
	}

	gr.delete()
	delete(gm.gameRooms, roomID)
	return nil
}

func (gm *GameManager) GetGameRoomUnsafe(roomID string) (GameRoom, error) {
	if r, ok := gm.gameRooms[roomID]; ok {
		return r, nil
	} else {
		return GameRoom{}, room.ErrRoomNotFound
	}
}

func (gm *GameManager) GetGameRoomSync(roomID string) (GameRoom, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()
	if r, ok := gm.gameRooms[roomID]; ok {
		return r, nil
	} else {
		return GameRoom{}, room.ErrRoomNotFound
	}
}

func (gm *GameManager) GetRoomInfoSync(roomID string) (room.RoomInfo, error) {
	r, err := gm.GetGameRoomSync(roomID)
	if err != nil {
		return room.RoomInfo{}, err
	}
	return r.Room.GetRoomInfoSync(), nil
}

func (gm *GameManager) JoinPlayerSync(roomID string, playerID string, conn *websocket.Conn) error {
	gr, err := gm.GetGameRoomSync(roomID)
	if err != nil {
		return err
	}

	rm := gr.Room

	if rm.PlayerExistsAndLeftSync(playerID) {
		err = rm.RejoinPlayerSync(playerID, conn)
	} else {
		err = rm.AddPlayerSync(playerID, conn)
	}

	if err != nil {
		return err
	}

	return nil
}
