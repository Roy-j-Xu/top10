package core

import "top10/core/game"

func InitCore() *GameManager {
	game.LoadQuestionSet()
	return NewGameManager()
}
