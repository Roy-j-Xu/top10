package core

import (
	"encoding/json"
	"log"
)

type MessageType string

const (
	Joined       MessageType = "joined"
	Broadcast    MessageType = "broadcast"
	Questions    MessageType = "questions"
	Ready        MessageType = "ready"
	AssignNumber MessageType = "assign-number"
)

type Message struct {
	Type MessageType
	Msg  any
}

type QuestionsMsg struct {
	Questions []string
	Guesser   int
}

type Messager interface {
	Broadcast(msg any, msgType MessageType)
	Message(msg any, playerID int, msgType MessageType)
}

type DebugMessager struct{}

func (d *DebugMessager) Broadcast(msg any, msgType MessageType) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("[Broadcast] (%s) %s", msgType, data)
}

func (d *DebugMessager) Message(msg any, playerID int, msgType MessageType) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("[Player %d] (%s) %s\n", playerID, msgType, string(data))
}
