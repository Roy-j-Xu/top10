package game

import "top10/core/room"

type GameMsgType string

const (
	G_BROADCAST      GameMsgType = "game:broadcast"
	G_NEW_QUESTIONS  GameMsgType = "game:new-questions"
	G_SET_QUESTION   GameMsgType = "game:set-question"
	G_ASSIGN_NUMBERS GameMsgType = "game:assign-numbers"
	G_FINISHED       GameMsgType = "game:finished"
	G_ERROR          GameMsgType = "game:error"

	GP_READY        GameMsgType = "game-player:ready"
	GP_SET_QUESTION GameMsgType = "game-player:set-question"
	GP_CHOOSE_ORDER GameMsgType = "game-player:choose-order"
)

func GameMsgOf(msgType GameMsgType, msg any) room.Message {
	return room.Message{
		Type: string(msgType),
		Msg:  msg,
	}
}
