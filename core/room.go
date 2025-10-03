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

func (room *Room) AddPlayerSync(player *Player) {
	room.Lock()
	defer room.Unlock()

	if room.Status != Waiting {
		return
	}

	// do not use room.SizeSync, otherwise deadlock
	player.ID = room.Size()

	if player.ID > 10 {
		return
	}

	room.Players = append(room.Players, player)

	// avoid locking
	go room.Broadcast(fmt.Sprintf("Player %d joined", player.ID), Broadcast)
}

func (room *Room) Size() int {
	return len(room.Players)
}

func (room *Room) SizeSync() int {
	room.Lock()
	defer room.Unlock()
	return len(room.Players)
}

func (room *Room) ReadyPlayerSync(playerID int) error {
	if playerID < 0 || playerID >= room.SizeSync() {
		return errors.New("player does not exist")
	}
	room.readyChan <- playerID
	go room.Broadcast(fmt.Sprintf("Player %d ready", playerID), Broadcast)
	return nil
}

func (room *Room) SetStatus(status Status) {
	if room.Status == status {
		return
	}
	log.Printf("Game status: %s -> %s", room.Status, status)
	room.Status = status
}

func (room *Room) Message(msg string, playerID int, msgType MessageType) {
	for _, msgr := range room.Messagers {
		msgr.Message(msg, playerID, msgType)
	}
}

func (room *Room) Broadcast(msg string, msgType MessageType) {
	for _, msgr := range room.Messagers {
		msgr.Broadcast(msg, msgType)
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
