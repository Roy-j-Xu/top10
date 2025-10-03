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

	room.AddPlayerSync(&core.Player{})
	room.AddPlayerSync(&core.Player{})
	room.AddPlayerSync(&core.Player{})

	room.ReadyPlayerSync(0)
	room.ReadyPlayerSync(1)
	room.ReadyPlayerSync(2)

	<-time.After(1 * time.Second)
}
