package logic

import (
	"errors"
	"sync"
)

type Player struct {
	ID      int
	Guesser bool
	Number  int
}

type Status int

const (
	Waiting Status = iota
	Playing
	Finished
)

type Room struct {
	Players      []*Player
	Question     string
	Status       Status
	readyChan    chan int
	continueChan chan bool
	mutex        sync.Mutex
}

func NewRoom() *Room {
	return &Room{
		Players:   []*Player{},
		Status:    Waiting,
		readyChan: make(chan int),
	}
}

func (room *Room) AddPlayer(player *Player) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	if len(room.Players) > 10 {
		return
	}

	room.Players = append(room.Players, player)
}

func (room *Room) Size() int {
	return len(room.Players)
}

// Wait for everyone to ready
func (room *Room) Wait() {
	numOfReady := 0
	playerReady := make([]bool, room.Size())

	for {
		playerID := <-room.readyChan
		if !playerReady[playerID] {
			playerReady[playerID] = true
			numOfReady++
		}

		if numOfReady == room.Size() {
			room.continueChan <- true
			return
		}
	}
}

func (room *Room) ReadyPlayer(playerID int) error {
	if 0 > playerID || playerID > room.Size() {
		return errors.New("player does not exist")
	}
	room.readyChan <- playerID
	return nil
}
