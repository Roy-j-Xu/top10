package game

import "top10/core/room"

type GameMsgType string

const (
	G_BROADCAST      GameMsgType = "topten:broadcast"
	G_NEW_QUESTIONS  GameMsgType = "topten:new-questions"
	G_SET_QUESTION   GameMsgType = "topten:set-question"
	G_ASSIGN_NUMBERS GameMsgType = "topten:assign-numbers"
	G_FINISHED       GameMsgType = "topten:finished"
	G_ERROR          GameMsgType = "topten:error"

	GP_READY        GameMsgType = "topten-player:ready"
	GP_SET_QUESTION GameMsgType = "topten-player:set-question"
	GP_CHOOSE_ORDER GameMsgType = "topten-player:choose-order"
)

func GameMsgOf(msgType GameMsgType, msg any) room.Message {
	return room.Message{
		Type: string(msgType),
		Msg:  msg,
	}
}
