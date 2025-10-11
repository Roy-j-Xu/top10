package room

import (
	"fmt"
	"log"
)

func (r *Room) Size() int {
	return len(r.Players)
}

func (r *Room) SizeSync() int {
	r.Lock()
	defer r.Unlock()
	return len(r.Players)
}

func (r *Room) GetPlayerSync(playerID string) (*Player, error) {
	r.Lock()
	defer r.Unlock()
	if p, ok := r.Players[playerID]; ok {
		return p, nil
	} else {
		return nil, ErrPlayerNotFound
	}
}

func (r *Room) GetAllPlayerIDsSync() []string {
	r.Lock()
	defer r.Unlock()
	ids := make([]string, 0, len(r.Players))
	for id := range r.Players {
		ids = append(ids, id)
	}
	return ids
}

func (r *Room) GetNumberOfReadiesSync() int {
	r.Lock()
	defer r.Unlock()
	count := 0
	for _, p := range r.Players {
		if p.Ready && !p.Left {
			count++
		}
	}
	return count
}

func (r *Room) UnreadyAllSync() {
	r.Lock()
	defer r.Unlock()
	for _, p := range r.Players {
		p.Ready = false
	}
}

func (r *Room) Message(msg Message, playerID string) {
	for _, msgr := range r.Messagers {
		msgr.Message(msg, playerID)
	}
}

func (r *Room) Broadcast(msg Message) {
	for _, msgr := range r.Messagers {
		msgr.Broadcast(msg)
	}
}

func (r *Room) BroadcastError(err error) {
	r.Broadcast(SystemMsgOf(S_ERROR, err.Error()))
}

func (r *Room) ListenToTimeout() error {
	r.Lock()
	if r.Timer == nil {
		return fmt.Errorf("unable to find timer in room %s: %w", r.ID, ErrInvalidRoom)
	}
	r.Unlock()

	<-r.Timer.C
	r.cancel()
	return nil
}

func (r *Room) ResetTimerSync() {
	r.Lock()
	defer r.Unlock()
	if r.Timer != nil {
		r.Timer.Reset(r.Timeout)
	}
}

func (r *Room) ResetTimerUnsafe() {
	if r.Timer != nil {
		r.Timer.Reset(r.Timeout)
	}
}

func (r *Room) SendToPlayerChannel(playerID string, msg Message) error {
	if p, err := r.GetPlayerSync(playerID); err == nil {
		p.msgChan <- msg
		return nil
	} else {
		return ErrPlayerNotFound
	}
}

func (r *Room) SendToPlayerChannel_READY(playerID string) error {
	return r.SendToPlayerChannel(playerID, SystemMsgOf(P_READY, "ready"))
}

func (r *Room) SendToPlayerChannel_LEFT(playerID string) error {
	r.SendToPlayerChannel(playerID, SystemMsgOf(P_LEFT, "left"))
	return nil
}

func (r *Room) Lock() {
	r.mutex.Lock()
}

func (r *Room) Unlock() {
	r.mutex.Unlock()
}

func (r *Room) Print() {
	r.Lock()
	defer r.Unlock()
	log.Println("----------")
	log.Println("Room ID: ", r.ID)
	log.Println("InGame: ", r.InGame)
	log.Println("Max Size: ", r.MaxSize)
	log.Println("Players: ")
	for id, p := range r.Players {
		log.Printf(" - ID: %s, Left: %v\n", id, p.Left)
	}
	log.Println("----------")
}
