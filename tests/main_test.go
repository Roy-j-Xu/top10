package tests

import (
	"testing"
	"time"
	"top10/core"
)

func TestRandomQuestion(t *testing.T) {
	core.LoadQuestionSet()
	t.Log(core.RandomQuestions(4))
}

func TestRun(t *testing.T) {
	core.LoadQuestionSet()
	room := core.NewRoom(nil) // use debug messager

	go room.Run()

	room.AddPlayer(&core.Player{})
	room.AddPlayer(&core.Player{})
	room.AddPlayer(&core.Player{})

	room.ReadyPlayer(0)
	room.ReadyPlayer(1)
	room.ReadyPlayer(2)

	<-time.After(1 * time.Second)
}
