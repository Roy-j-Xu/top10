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
			g.Room().WaitAllSync()
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
	Questions    []string

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
	g.generateQuestions()
	g.assignNumbers()
}

func (g *Game) assignNumbers() {
	numbers := randomKFromN(g.Size(), 10)
	for k, playerID := range g.TurnOrder {
		playerNumber := numbers[k]
		g.PlayerStates[playerID].Number = playerNumber
		g.Room().Message(GameMsgOf(G_ASSIGN_NUMBERS, playerNumber), playerID)
	}
}

func (g *Game) generateQuestions() {
	g.Questions = RandomQuestions(4)
	msgData := QuestionsMsg{
		Questions: g.Questions,
		Guesser:   g.GuesserID,
	}
	g.Room().Broadcast(GameMsgOf(G_NEW_QUESTIONS, msgData))
}
