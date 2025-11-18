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
	ID      string
	Conn    *websocket.Conn
	Ready   bool
	Left    bool         // represents disconnection
	msgChan chan Message // handles game related messages
	sysChan chan Message // handles system messages
}

type Room struct {
	ID       string
	MaxSize  int
	GameName string
	InGame   bool
	Paused   bool

	Players map[string]*Player

	Messengers []Messenger

	sysChan chan string

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
		GameName:   "top10",
		Players:    make(map[string]*Player),
		Messengers: msgrs,
		sysChan:    make(chan string, maxSize),
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

// translate raw messages from socket into corresponding channels
func (r *Room) handlePlayerMessages(playerID string) {
	player, err := r.GetPlayerSync(playerID)
	if err != nil {
		log.Printf("unable to listen message from player %s: %s", playerID, err.Error())
		return
	}
	log.Printf("listening for messages from player %s", playerID)
	defer func() {
		log.Printf("stopped listening messages from player %s", playerID)
		r.SendToSysChannelSync_LEFT(playerID)
		player.Conn.Close()
	}()

	for {
		var msg Message
		if err := player.Conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			return
		}

		if msg.Type == string(SP_READY) || msg.Type == string(SP_LEFT) {
			r.SendToSysChannelSync(playerID, msg)
		} else if !r.Paused {
			r.SendToPlayerChannelSync(playerID, msg)
		}
	}
}

func (r *Room) listenPlayerSysMessageSync(playerID string) error {
	player, err := r.GetPlayerSync(playerID)
	if err != nil {
		return fmt.Errorf("listening to player \"%s\" in room \"%s\": %w", playerID, r.ID, ErrPlayerNotFound)
	}

	for {
		select {
		case <-r.ctx.Done():
			r.RemovePlayerSync(playerID)
			return nil
		case msg := <-player.sysChan:
			log.Printf("received message from player \"%s\": %v", playerID, msg)
			switch msg.Type {
			case string(SP_READY):
				r.sysChan <- playerID
				r.ResetTimerSync()
			case string(SP_LEFT):
				r.LeavePlayerSync(playerID)
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
			ID:      playerID,
			Conn:    conn,
			sysChan: make(chan Message, 1),
			msgChan: make(chan Message, 10),
		}
		r.Players[playerID] = player
		r.ResetTimerUnsafe() // use unsafe to prevent deadlock

		go r.Broadcast(JoinedMsgOf(playerID, r.GetRoomInfoUnsafe()))
	}

	go r.handlePlayerMessages(playerID)
	go r.listenPlayerSysMessageSync(playerID)

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

	if r.noPlayerIsLeftUnsafe() {
		r.Paused = false
	}

	go r.Broadcast(JoinedMsgOf(playerID, r.GetRoomInfoUnsafe()))
	go r.handlePlayerMessages(playerID)
	return nil
}

func (r *Room) LeavePlayerSync(playerID string) error {
	r.Lock()
	defer r.Unlock()
	player, err := r.GetPlayerUnsafe(playerID)
	if err != nil {
		return fmt.Errorf("rejoining player \"%s\" to room \"%s\": %w", playerID, r.ID, ErrPlayerNotFound)
	}

	if r.InGame {
		player.Ready = false
		player.Left = true
		r.Paused = true
		r.Broadcast(SystemMsgOf(S_BROADCAST, fmt.Sprintf("player \"%s\" disconnected, game may pause", playerID)))
	} else {
		r.RemovePlayerSync(player.ID)
	}

	return nil
}

func (r *Room) RemovePlayerSync(playerID string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Players[playerID]; !ok {
		return fmt.Errorf("removing player \"%s\" from room \"%s\": %w", playerID, r.ID, ErrPlayerNotFound)
	}
	delete(r.Players, playerID)

	go r.Broadcast(LeftMsgOf(playerID, r.GetRoomInfoUnsafe()))

	if r.SizeSync() <= 0 {
		log.Println("no player in room, shutting down")
		r.Shutdown()
	}

	return nil
}

// Wait for every player to ready
func (r *Room) waitAllSync() error {
	for {
		select {
		case <-time.After(r.Timeout):
			return fmt.Errorf("waiting for players in room \"%s\": %w", r.ID, ErrTimeout)
		case playerID := <-r.sysChan:
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
				r.Broadcast(ReadyMsgOf(playerID, r.GetRoomInfoUnsafe()))
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
	if err := r.waitAllSync(); err != nil {
		r.Broadcast(SystemMsgOf(S_ERROR, "wait for start: waiting for players timed out"))
		return err
	}

	r.Lock()
	defer r.Unlock()
	r.InGame = true
	r.Broadcast(SystemMsgOf(S_START, r.GetRoomInfoUnsafe()))
	return nil
}
