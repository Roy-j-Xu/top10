package room

import (
	"context"
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

func (r *Room) PlayerExistsAndLeftSync(playerID string) bool {
	r.Lock()
	defer r.Unlock()
	if p, ok := r.Players[playerID]; ok {
		return p.Left
	} else {
		return false
	}
}

func (r *Room) UnreadyAllSync() {
	r.Lock()
	defer r.Unlock()
	for _, p := range r.Players {
		p.Ready = false
	}
}

func (r *Room) Message(msg Message, playerID string) {
	for _, msgr := range r.Messengers {
		msgr.Message(msg, playerID)
	}
}

func (r *Room) Broadcast(msg Message) {
	for _, msgr := range r.Messengers {
		msgr.Broadcast(msg)
	}
}

func (r *Room) BroadcastError(err error) {
	r.Broadcast(SystemMsgOf(S_ERROR, err.Error()))
}

func (r *Room) ListenToTimeout() {
	<-r.Timer.C
	r.cancel()
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

func (r *Room) SendToReadyChannel(playerID string, msg Message) error {
	if p, err := r.GetPlayerSync(playerID); err == nil {
		p.readyChan <- msg
		return nil
	} else {
		return ErrPlayerNotFound
	}
}

func (r *Room) SendToReadyChannel_READY(playerID string) error {
	return r.SendToReadyChannel(playerID, SystemMsgOf(SP_READY, "player ready"))
}

func (r *Room) SendToReadyChannel_LEFT(playerID string) error {
	r.SendToReadyChannel(playerID, SystemMsgOf(SP_LEFT, "player left or hang up"))
	return nil
}

func (r *Room) Lock() {
	r.mutex.Lock()
}

func (r *Room) Unlock() {
	r.mutex.Unlock()
}

func (r *Room) Shutdown() {
	r.cancel()
}

func (r *Room) StopCtx() context.Context {
	return r.ctx
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
