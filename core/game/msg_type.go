package game

import "top10/core/room"

type GameMsgType string

const (
	G_BROADCAST      GameMsgType = "topten:broadcast"
	G_START          GameMsgType = "topten:start"
	G_GAME_INFO      GameMsgType = "topten:game-info"
	G_SET_QUESTION   GameMsgType = "topten:set-question"
	G_ASSIGN_NUMBERS GameMsgType = "topten:assign-numbers"
	G_REVEAL_NUMBER  GameMsgType = "topten:reveal-number"
	G_FINISHED       GameMsgType = "topten:finished"
	G_ERROR          GameMsgType = "topten:error"

	GP_READY        GameMsgType = "topten-player:ready"
	GP_SET_QUESTION GameMsgType = "topten-player:set-question"
	GP_CHOOSE_ORDER GameMsgType = "topten-player:choose-order"
)

type GameInfo struct {
	Turn         int            `json:"turn"`
	MaxTurn      int            `json:"maxTurn"`
	TurnOrder    []string       `json:"turnOrder"`
	Guesser      string         `json:"guesser"`
	Questions    []string       `json:"questions"`
	UsedQuestion string         `json:"usedQuestion"`
	Numbers      map[string]int `json:"numbers"`
	State        string         `json:"state"`
}

func GameMsgOf(msgType GameMsgType, msg any) room.Message {
	return room.Message{
		Type: string(msgType),
		Msg:  msg,
	}
}
