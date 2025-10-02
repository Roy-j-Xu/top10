package core

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID     int
	Number int
	Conn   *websocket.Conn
}

type Status string

const (
	Waiting  Status = "Waiting"
	Playing  Status = "Playing"
	Finished Status = "Finished"
)

type Messager interface {
	Broadcast(msg string)
	Message(msg string, playerID int)
}

type Room struct {
	Players   []*Player
	Questions []string
	GuesserID int
	Status    Status

	Messager Messager

	readyChan    chan int
	continueChan chan bool

	mutex sync.Mutex
}

func NewRoom(msgr Messager) *Room {
	if msgr == nil {
		msgr = &DebugMessager{}
	}
	return &Room{
		Players:      []*Player{},
		Status:       Waiting,
		Messager:     msgr,
		readyChan:    make(chan int, 10),
		continueChan: make(chan bool, 1),
	}
}

func (room *Room) AddPlayer(player *Player) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	if room.Size() > 10 {
		return
	}

	player.ID = room.Size()

	room.Players = append(room.Players, player)
}

func (room *Room) Size() int {
	return len(room.Players)
}

func (room *Room) ReadyPlayer(playerID int) error {
	if playerID < 0 || playerID > room.Size() {
		return errors.New("player does not exist")
	}
	room.readyChan <- playerID
	return nil
}

func (room *Room) SetStatus(status Status) {
	log.Printf("Game status: %s -> %s", room.Status, status)
	room.Status = status
}

func (room *Room) Print() {
	log.Println("-----")
	log.Println("Player count: ", room.Size())
	log.Println("Questions:\n", strings.Join(room.Questions, "\n"))
	log.Println("Guesser: ", room.GuesserID)
	log.Println("Status: ", room.Status)
	log.Println("-----")
}
