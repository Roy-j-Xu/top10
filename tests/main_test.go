package tests

import (
	"testing"
	"time"
	"top10/core/game"
	"top10/core/room"
)

func TestRandomQuestion(t *testing.T) {
	game.LoadQuestionSet()
	t.Log(game.RandomQuestions(4))
}

func TestReadyRoom(t *testing.T) {
	game.LoadQuestionSet()

	room, err := room.NewRoomDebug("new room", 10) // use debug messager
	if err != nil {
		t.Fatal(err)
	}

	room.Timeout = 1 * time.Second

	room.AddPlayerSync("player1", nil)
	room.AddPlayerSync("No. 2", nil)
	room.AddPlayerSync("P3", nil)

	room.Print()

	go room.WaitForStartSync()

	wait10msAnd(func() { room.SendToPlayerChannel_READY("player1") })
	wait10msAnd(func() { room.SendToPlayerChannel_READY("fake player") })
	wait10msAnd(func() { room.SendToPlayerChannel_READY("No. 2") })
	wait10msAnd(func() { room.SendToPlayerChannel_LEFT("player1") })
	wait10msAnd(func() { room.AddPlayerSync("player 4", nil) })
	wait10msAnd(func() { room.SendToPlayerChannel_READY("P3") })
	wait10msAnd(func() { room.SendToPlayerChannel_READY("player 4") })

	<-time.After(2 * time.Second)
}

func TestGame(t *testing.T) {
	game.LoadQuestionSet()

	room, err := room.NewRoomDebug("new room", 10)
	if err != nil {
		t.Fatal(err)
	}

	room.AddPlayerSync("1", nil)

	go room.WaitForStartSync()

	room.AddPlayerSync("2", nil)
	room.AddPlayerSync("3", nil)

	wait10msAnd(func() { room.SendToPlayerChannel_READY("1") })
	wait10msAnd(func() { room.SendToPlayerChannel_READY("2") })
	wait10msAnd(func() { room.SendToPlayerChannel_READY("3") })
	<-time.After(500 * time.Millisecond)

	game := game.NewGame(room)

	go game.Start()

	wait10msAnd(game.Print)

	wait10msAnd(func() { room.SendToPlayerChannel_READY("1") })
	wait10msAnd(func() { room.SendToPlayerChannel_READY("2") })
	wait10msAnd(func() { room.SendToPlayerChannel_READY("3") })

	wait10msAnd(game.Print)
}
