package socket

import (
	"log"
	"top10/core"
)

type WebSocketMessager struct {
	Room *core.Room
}

type MessageType string

type Message struct {
	Type MessageType
	Msg  any
}

func (w *WebSocketMessager) Broadcast(msg any) {
	w.Room.Lock()
	defer w.Room.Unlock()

	for _, p := range w.Room.Players {
		if p.Conn != nil {
			if err := p.Conn.WriteJSON(Message{
				Type: "broadcast",
				Msg:  msg,
			}); err != nil {
				log.Printf("Error sending to player %d: %v", p.ID, err)
			}
		}
	}
}

func (w *WebSocketMessager) Message(msg any, playerID int) {
	w.Room.Lock()
	defer w.Room.Unlock()

	for _, p := range w.Room.Players {
		if p.ID == playerID && p.Conn != nil {
			if err := p.Conn.WriteJSON(Message{
				Type: "broadcast",
				Msg:  msg,
			}); err != nil {
				log.Printf("Error sending to player %d: %v", p.ID, err)
			}
			break
		}
	}
}
