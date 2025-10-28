package core

import "top10/core/game"

func InitCore() {
	game.LoadQuestionSet()

	gm := NewGameManager()
	gm.HandleHTTP()
}
