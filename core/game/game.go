package game

import "top10/core/room"

type Status struct {
	Name     string
	OnStatus func(*Game)
}

var (
	Start = &Status{
		Name: "Start",
		OnStatus: func(g *Game) {
			g.Room().Broadcast(GameMsgOf(G_START, g.GetGameInfoUnsafe()))
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
				g.RepeatStatus()
			}
		},
	}
	Finished = &Status{
		Name: "Finished",
		OnStatus: func(g *Game) {
			g.Room().Broadcast(GameMsgOf(G_FINISHED, ""))
		},
	}
)

type PlayerState struct {
	Number int
}

type Game struct {
	Status *Status

	TurnOrder     []string
	TurnNumber    int
	MaxTurn       int
	PlayerNumbers map[string]int
	GuesserID     string
	Questions     []string
	UsedQuestion  string

	room *room.Room
}

func NewGame(r *room.Room) *Game {
	playerIDs := r.GetAllPlayerIDsSync()
	game := &Game{
		Status:        Start,
		PlayerNumbers: make(map[string]int),
		TurnOrder:     playerIDs,
		TurnNumber:    0,
		MaxTurn:       len(playerIDs),
		room:          r,
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
	g.Questions = RandomQuestions(4)
	g.UsedQuestion = ""
	g.clearNumbers()
	g.Room().Broadcast(GameMsgOf(G_GAME_INFO, g.GetGameInfoUnsafe()))

	g.setQuestion()
	g.assignNumbers()
	g.Room().Broadcast(GameMsgOf(G_START_GUESSING, g.GetGameInfoUnsafe()))
	g.Room().WaitForAllMessages(string(GP_READY))
}

func (g *Game) assignNumbers() {
	numbers := randomKFromN(g.Size(), 10)
	for k, playerID := range g.TurnOrder {
		playerNumber := numbers[k]
		g.PlayerNumbers[playerID] = playerNumber
	}
}

func (g *Game) clearNumbers() {
	for _, playerID := range g.TurnOrder {
		g.PlayerNumbers[playerID] = 0
	}
}

// wait for guesser to choose one question
func (g *Game) setQuestion() {
	setQsMsg, err := g.Room().WaitUntilGetMessage(g.GuesserID, string(GP_SET_QUESTION))
	question, ok := setQsMsg.Msg.(string)
	if err != nil || !ok {
		g.Room().Broadcast(GameMsgOf(G_ERROR, "error reading question message"))
		g.Room().Shutdown()
	}
	g.UsedQuestion = question
	g.Questions = nil
}
