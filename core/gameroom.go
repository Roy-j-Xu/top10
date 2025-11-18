package core

import (
	"log"
	"top10/core/game"
	"top10/core/room"
)

type GameRoom struct {
	Game *game.Game
	Room *room.Room
}

func newGameRoom(roomName string, maxSize int) (GameRoom, error) {
	r, err := room.NewRoomWebSocket(roomName, maxSize)
	if err != nil {
		return GameRoom{}, err
	}
	return GameRoom{Room: r}, nil
}

func (gr GameRoom) delete() {
	gr.Room.Shutdown()
}

func (gr GameRoom) runGame(gameName string) error {
	game := game.NewGame(gr.Room)
	gr.Game = game

	game.Start() // game starts here

	gr.delete()

	log.Printf("game %s stopped", gameName)
	return nil
}
