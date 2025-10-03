package core

import "log"

type MessageType string

const (
	Questions    MessageType = "questions"
	AssignNumber MessageType = "number"
	Broadcast    MessageType = "broadcast"
)

type Message struct {
	Type MessageType
	Msg  any
}

type Messager interface {
	Broadcast(msg any, msgType MessageType)
	Message(msg any, playerID int, msgType MessageType)
}

type DebugMessager struct{}

func (d *DebugMessager) Broadcast(msg any, msgType MessageType) {
	log.Printf("[Broadcast] (%s) %s", msgType, msg)
}

func (d *DebugMessager) Message(msg any, playerID int, msgType MessageType) {
	log.Printf("[Player %d] (%s) %s\n", playerID, msgType, msg)
}
