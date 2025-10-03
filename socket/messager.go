package socket

import (
	"log"
	"top10/core"
)

type WebSocketMessager struct {
	Room *core.Room
}

func (w *WebSocketMessager) Broadcast(msg any, msgType core.MessageType) {
	w.Room.Lock()
	defer w.Room.Unlock()

	for _, p := range w.Room.Players {
		if p.Conn != nil {
			if err := p.Conn.WriteJSON(core.Message{
				Type: msgType,
				Msg:  msg,
			}); err != nil {
				log.Printf("Error sending to player %d: %v", p.ID, err)
			}
		}
	}
}

func (w *WebSocketMessager) Message(msg any, playerID int, msgType core.MessageType) {
	w.Room.Lock()
	defer w.Room.Unlock()

	for _, p := range w.Room.Players {
		if p.ID == playerID && p.Conn != nil {
			if err := p.Conn.WriteJSON(core.Message{
				Type: msgType,
				Msg:  msg,
			}); err != nil {
				log.Printf("Error sending to player %d: %v", p.ID, err)
			}
			break
		}
	}
}
