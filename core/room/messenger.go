package room

import (
	"encoding/json"
	"log"
)

type Messenger interface {
	Broadcast(msg Message)
	Message(msg Message, playerID string)
}

type DebugMessenger struct{}

func (d *DebugMessenger) Broadcast(msg Message) {
	data, err := stringify(msg.Msg)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("[Broadcast] (%s) %s", msg.Type, data)
}

func (d *DebugMessenger) Message(msg Message, playerID string) {
	data, err := stringify(msg.Msg)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("[Message to %s] (%s) %s\n", playerID, msg.Type, string(data))
}

func stringify(input any) (string, error) {
	var out string

	switch v := input.(type) {
	case string:
		out = v
	default:
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return "", err
		}
		out = string(data)
	}

	return out, nil
}
