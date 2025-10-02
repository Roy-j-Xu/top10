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

func TestRoomWait(t *testing.T) {
	core.LoadQuestionSet()
	room := core.NewRoom()

	go room.Run()

	room.AddPlayer(&core.Player{})
	room.AddPlayer(&core.Player{})
	room.AddPlayer(&core.Player{})

	room.ReadyPlayer(0)
	room.ReadyPlayer(2)

	t.Log(room.Status)

	room.ReadyPlayer(1)

	<-time.After(1 * time.Second)
	t.Log(room.Status)
}
