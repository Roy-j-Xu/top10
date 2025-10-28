package room

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID        string
	Conn      *websocket.Conn
	Ready     bool
	Left      bool         // represents disconnection
	msgChan   chan Message // handles game related messages
	readyChan chan Message // handles connection
}

type Room struct {
	ID      string
	MaxSize int
	InGame  bool

	Players map[string]*Player

	Messengers []Messenger

	readyChan chan string

	// for timeout
	Timeout time.Duration
	Timer   *time.Timer
	ctx     context.Context
	cancel  context.CancelFunc

	mutex sync.Mutex
}

func NewRoomDebug(roomID string, maxSize int) (*Room, error) {
	if maxSize >= 20 {
		return nil, fmt.Errorf("room \"%s\" is too large: %w", roomID, ErrInvalidRoom)
	}

	msgrs := []Messenger{&DebugMessenger{}}

	timeout := 10 * time.Minute
	timer := time.NewTimer(timeout)
	ctx, cancel := context.WithCancel(context.Background())

	room := &Room{
		ID:         roomID,
		MaxSize:    maxSize,
		Players:    make(map[string]*Player),
		Messengers: msgrs,
		readyChan:  make(chan string, maxSize),
		Timer:      timer,
		Timeout:    timeout,
		ctx:        ctx,
		cancel:     cancel,
	}
	go room.ListenToTimeout()

	return room, nil
}

func NewRoomWebSocket(roomID string, maxSize int) (*Room, error) {
	room, err := NewRoomDebug(roomID, maxSize)
	if err != nil {
		return nil, err
	}
	room.Messengers = append(room.Messengers, &WebSocketMessenger{Players: room.Players})
	return room, nil
}

func (r *Room) ListenPlayerReadySync(playerID string) error {
	player, err := r.GetPlayerSync(playerID)
	if err != nil {
		return fmt.Errorf("listening to player \"%s\" in room \"%s\": %w", playerID, r.ID, ErrPlayerNotFound)
	}

	for {
		select {
		case <-r.ctx.Done():
			r.Message(SystemMsgOf(S_LEFT, "you are disconnected from game"), playerID)
			r.RemovePlayerSync(playerID)
			return nil
		case msg := <-player.readyChan:
			log.Printf("received message from player \"%s\": %v", playerID, msg)
			switch msg.Type {
			case string(SP_READY):
				r.readyChan <- playerID
				r.ResetTimerSync()
			case string(SP_LEFT):
				if r.InGame {
					r.Lock()
					player.Ready = false
					player.Left = true
					r.Unlock()
					r.Broadcast(SystemMsgOf(S_BROADCAST, fmt.Sprintf("player \"%s\" disconnected, game may pause", playerID)))
				} else {
					r.RemovePlayerSync(player.ID)
				}
				if r.SizeSync() <= 0 {
					log.Println("no player in room, shutting down")
					r.Shutdown()
				}
			default:
				r.Broadcast(SystemMsgOf(S_ERROR, fmt.Sprintf("unknown message type: %s", msg.Type)))
			}
		}
	}
}

func (r *Room) AddPlayerSync(playerID string, conn *websocket.Conn) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Players[playerID]; ok {
		return fmt.Errorf("adding player \"%s\" to room \"%s\": %w", playerID, r.ID, ErrPlayerExists)
	} else {
		// do not use r.SizeSync, otherwise deadlock
		if r.Size() >= r.MaxSize {
			return fmt.Errorf("adding player \"%s\" to room \"%s\", exceed max number: %w", playerID, r.ID, ErrInvalidRoom)
		}
		player := &Player{
			ID:        playerID,
			Conn:      conn,
			readyChan: make(chan Message, 1),
			msgChan:   make(chan Message, 10),
		}
		r.Players[playerID] = player
		r.ResetTimerUnsafe() // use unsafe to prevent deadlock

		go r.Broadcast(SystemMsgOf(S_JOINED, fmt.Sprintf("player \"%s\" joined", playerID)))
		go r.Message(SystemMsgOf(S_JOINED, playerID), playerID)
	}

	go r.ListenPlayerReadySync(playerID)

	return nil
}

func (r *Room) RejoinPlayerSync(playerID string, conn *websocket.Conn) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Players[playerID]; !ok {
		return fmt.Errorf("rejoining player \"%s\" to room \"%s\": %w", playerID, r.ID, ErrPlayerNotFound)
	}

	r.Players[playerID].Left = false
	r.Players[playerID].Conn = conn
	r.ResetTimerUnsafe() // use unsafe to prevent deadlock

	go r.Broadcast(SystemMsgOf(S_JOINED, fmt.Sprintf("player \"%s\" rejoined", playerID)))
	return nil
}

func (r *Room) RemovePlayerSync(playerID string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Players[playerID]; !ok {
		return fmt.Errorf("removing player \"%s\" from room \"%s\": %w", playerID, r.ID, ErrPlayerNotFound)
	}
	delete(r.Players, playerID)

	go r.Broadcast(SystemMsgOf(S_LEFT, fmt.Sprintf("player \"%s\" left", playerID)))

	return nil
}

// Wait for every player to ready
func (r *Room) WaitAllSync() error {
	for {
		select {
		case <-time.After(r.Timeout):
			return fmt.Errorf("waiting for players in room \"%s\": %w", r.ID, ErrTimeout)
		case playerID := <-r.readyChan:
			// no lock here, use Sync methods
			player, err := r.GetPlayerSync(playerID)
			if err != nil {
				r.Broadcast(SystemMsgOf(S_ERROR, "readying player who is not in room"))
				continue
			}

			if !player.Ready {
				player.Ready = true
				numberOfReadies := r.GetNumberOfReadiesSync()
				roomSize := r.SizeSync()
				go r.Broadcast(SystemMsgOf(
					S_BROADCAST,
					fmt.Sprintf("player \"%s\" is ready (%d/%d)", playerID, numberOfReadies, roomSize),
				))
				if numberOfReadies >= roomSize {
					r.UnreadyAllSync()
					return nil
				}
			}

		}
	}
}

func (r *Room) WaitForStartSync() error {
	r.Broadcast(SystemMsgOf(S_BROADCAST, "wait for start"))
	if err := r.WaitAllSync(); err != nil {
		r.Broadcast(SystemMsgOf(S_ERROR, "wait for start: waiting for players timed out"))
		return err
	}

	r.Lock()
	r.InGame = true
	r.Unlock()

	r.Broadcast(SystemMsgOf(S_START, "all players are ready, game starting"))
	return nil
}
