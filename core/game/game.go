package game

import (
	"top10/core/room"
)

type Status struct {
	Name     string
	OnStatus func(*Game)
}

var (
	Start = &Status{
		Name: "Start",
		OnStatus: func(g *Game) {
			if g.Size() > 0 {
				g.SetStatus(Playing)
			}
		},
	}
	Playing = &Status{
		Name: "Playing",
		OnStatus: func(g *Game) {
			g.nextTurn()
			if g.TurnNumber >= g.Size() {
				g.SetStatus(Finished)
				return
			} else {
				g.Status.OnStatus(g)
			}
		},
	}
	Finished = &Status{
		Name: "Finished",
		OnStatus: func(g *Game) {
			g.Room().Broadcast(GameMsgOf(G_FINISHED, "game finished"))
		},
	}
)

type PlayerState struct {
	Number int
}

type Game struct {
	Status *Status

	TurnOrder    []string
	TurnNumber   int
	MaxTurn      int
	PlayerStates map[string]*PlayerState
	GuesserID    string
	Question     string

	room *room.Room
}

func NewGame(r *room.Room) *Game {
	playerIDs := r.GetAllPlayerIDsSync()
	game := &Game{
		Status:       Start,
		PlayerStates: make(map[string]*PlayerState),
		TurnOrder:    playerIDs,
		TurnNumber:   0,
		MaxTurn:      len(playerIDs),
		room:         r,
	}

	for _, playerID := range playerIDs {
		game.AddNewPlayerState(playerID)
	}

	return game
}

func (g *Game) Start() {
	g.SetStatus(Start)
}

func (g *Game) nextTurn() {
	g.TurnNumber++
	g.GuesserID = g.TurnOrder[g.TurnNumber-1]
	g.setQuestion()
	g.assignNumbers()
	g.Room().WaitForAllMessages(string(GP_READY))
}

func (g *Game) assignNumbers() {
	numbers := randomKFromN(g.Size(), 10)
	for k, playerID := range g.TurnOrder {
		playerNumber := numbers[k]
		g.PlayerStates[playerID].Number = playerNumber
		g.Room().Message(GameMsgOf(G_ASSIGN_NUMBERS, playerNumber), playerID)
	}
}

func (g *Game) setQuestion() {
	questions := RandomQuestions(4)
	g.Room().Message(GameMsgOf(G_NEW_QUESTIONS, questions), g.GuesserID)
	setQsMsg, err := g.Room().WaitUntilGetMessage(g.GuesserID, string(GP_SET_QUESTION))
	question, ok := setQsMsg.Msg.(string)
	if err != nil || !ok {
		g.Room().Broadcast(GameMsgOf(G_ERROR, "error reading question message"))
		g.Room().Shutdown()
	}
	g.Question = question
	g.Room().Broadcast(GameMsgOf(G_SET_QUESTION, question))
}
