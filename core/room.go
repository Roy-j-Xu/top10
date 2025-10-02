package core

import (
	"errors"
	"fmt"
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
	Broadcast(msg any)
	Message(msg any, playerID int)
}

type Room struct {
	Players   []*Player
	Questions []string
	GuesserID int
	Status    Status

	Messagers []Messager

	readyChan    chan int
	continueChan chan bool

	mutex sync.Mutex
}

func NewRoom(msgrs []Messager) *Room {
	if msgrs == nil {
		msgrs = []Messager{&DebugMessager{}}
	}
	return &Room{
		Players:      []*Player{},
		Status:       Waiting,
		Messagers:    msgrs,
		readyChan:    make(chan int, 10),
		continueChan: make(chan bool, 1),
	}
}

func (room *Room) AddPlayer(player *Player) {
	room.Lock()

	if room.Size() > 10 {
		return
	}

	player.ID = room.Size()

	room.Players = append(room.Players, player)

	room.Unlock()

	room.Broadcast(fmt.Sprintf("Player %d joined", player.ID))
}

func (room *Room) Size() int {
	return len(room.Players)
}

func (room *Room) ReadyPlayer(playerID int) error {
	if playerID < 0 || playerID > room.Size() {
		return errors.New("player does not exist")
	}
	room.readyChan <- playerID
	room.Broadcast(fmt.Sprintf("Player %d ready", playerID))
	return nil
}

func (room *Room) SetStatus(status Status) {
	log.Printf("Game status: %s -> %s", room.Status, status)
	room.Status = status
}

func (room *Room) Message(msg string, playerID int) {
	for _, msgr := range room.Messagers {
		msgr.Message(msg, playerID)
	}
}

func (room *Room) Broadcast(msg string) {
	for _, msgr := range room.Messagers {
		msgr.Broadcast(msg)
	}
}

func (room *Room) Lock() {
	room.mutex.Lock()
}

func (room *Room) Unlock() {
	room.mutex.Unlock()
}

func (room *Room) Print() {
	log.Println("-----")
	log.Println("Player count: ", room.Size())
	log.Println("Questions:\n", strings.Join(room.Questions, "\n"))
	log.Println("Guesser: ", room.GuesserID)
	log.Println("Status: ", room.Status)
	log.Println("-----")
}
