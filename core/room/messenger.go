package room

import (
	"encoding/json"
	"log"
)

type Message struct {
	Type string
	Msg  any
}

type Messenger interface {
	Broadcast(msg Message)
	Message(msg Message, playerID string)
}

type SystemMsgType string

const (
	// Messages from system to players
	S_JOINED    SystemMsgType = "system:joined"
	S_LEFT      SystemMsgType = "system:left"
	S_START     SystemMsgType = "system:start"
	S_BROADCAST SystemMsgType = "system:broadcast"
	S_ERROR     SystemMsgType = "system:error"

	// Messages from players to system
	SP_READY SystemMsgType = "system-player:ready"
	SP_LEFT  SystemMsgType = "system-player:leave"
)

func SystemMsgOf(msgType SystemMsgType, msg any) Message {
	return Message{
		Type: string(msgType),
		Msg:  msg,
	}
}

type DebugMessenger struct{}

func (d *DebugMessenger) Broadcast(msg Message) {
	data, err := json.Marshal(msg.Msg)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("[Broadcast] (%s) %s", msg.Type, data)
}

func (d *DebugMessenger) Message(msg Message, playerID string) {
	data, err := json.Marshal(msg.Msg)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("[Message to %s] (%s) %s\n", playerID, msg.Type, string(data))
}
