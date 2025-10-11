package room

import (
	"log"
	"sync"
)

type WebSocketMessager struct {
	Players map[string]*Player
	mutex   sync.Mutex
}

func (w *WebSocketMessager) Broadcast(msg Message) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, p := range w.Players {
		if p.Conn != nil {
			if err := p.Conn.WriteJSON(Message{
				Type: msg.Type,
				Msg:  msg.Msg,
			}); err != nil {
				log.Printf("Error sending to player %s: %v", p.ID, err)
			}
		}
	}
}

func (w *WebSocketMessager) Message(msg Message, playerID string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	for _, p := range w.Players {
		if p.ID == playerID && p.Conn != nil {
			if err := p.Conn.WriteJSON(Message{
				Type: msg.Type,
				Msg:  msg.Msg,
			}); err != nil {
				log.Printf("Error sending to player %s: %v", p.ID, err)
			}
			break
		}
	}
}
