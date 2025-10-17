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

	wait10msAnd(func() { room.SendToReadyChannel_READY("player1") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("fake player") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("No. 2") })
	wait10msAnd(func() { room.SendToReadyChannel_LEFT("player1") })
	wait10msAnd(func() { room.AddPlayerSync("player 4", nil) })
	wait10msAnd(func() { room.SendToReadyChannel_READY("P3") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("player 4") })

	<-time.After(2 * time.Second)
}

func TestGameDisconnect(t *testing.T) {
	game.LoadQuestionSet()

	room, err := room.NewRoomDebug("new room", 10)
	if err != nil {
		t.Fatal(err)
	}

	room.AddPlayerSync("1", nil)

	go room.WaitForStartSync()

	room.AddPlayerSync("2", nil)
	room.AddPlayerSync("3", nil)

	wait10msAnd(func() { room.SendToReadyChannel_READY("1") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("2") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("3") })

	<-time.After(500 * time.Millisecond)

	g := game.NewGame(room)

	go g.Start()

	wait10msAnd(g.Print)

	wait10msAnd(func() { room.SendToReadyChannel_READY("1") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("2") })
	wait10msAnd(func() { room.SendToReadyChannel_LEFT("1") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("3") })
	wait10msAnd(func() { room.RejoinPlayerSync("1", nil) })
	wait10msAnd(func() { room.SendToReadyChannel_READY("1") })

	wait10msAnd(g.Print)
}

func TestGameFlow(t *testing.T) {
	game.LoadQuestionSet()

	room, err := room.NewRoomDebug("new room", 10)
	if err != nil {
		t.Fatal(err)
	}

	room.AddPlayerSync("1", nil)
	room.AddPlayerSync("2", nil)
	room.AddPlayerSync("3", nil)

	go room.WaitForStartSync()

	wait10msAnd(func() { room.SendToReadyChannel_READY("1") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("2") })
	wait10msAnd(func() { room.SendToReadyChannel_READY("3") })

	<-time.After(500 * time.Millisecond)

	g := game.NewGame(room)

	go g.Start()

	wait10msAnd(g.Print)

	guesser := g.GuesserID
	room.SendToPlayerChannel(guesser, game.GameMsgOf(game.GP_SET_QUESTION, "12321"))

	wait10msAnd(g.Print)
}
