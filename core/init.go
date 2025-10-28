package core

import (
	"net/http"
	"top10/core/game"
)

func InitCore() {
	game.LoadQuestionSet()

	gm := NewGameManager()
	gm.HandleHTTP()

	http.ListenAndServe(":8000", nil)
}
