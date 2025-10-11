package game

import "top10/core/room"

type GameMsgType string

const (
	G_BROADCAST      GameMsgType = "game:broadcast"
	G_NEW_QUESTIONS  GameMsgType = "game:new-questions"
	G_ASSIGN_NUMBERS GameMsgType = "game:assign-numbers"
	G_FINISHED       GameMsgType = "game:finished"
)

// for G_NEW_QUESTIONS
type QuestionsMsg struct {
	Questions []string
	Guesser   string
}

func GameMsgOf(msgType GameMsgType, msg any) room.Message {
	return room.Message{
		Type: string(msgType),
		Msg:  msg,
	}
}
