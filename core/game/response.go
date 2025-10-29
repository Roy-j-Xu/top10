package game

type GameInfoResponse struct {
	Turn         int      `json:"turn"`
	Guesser      string   `json:"guesser"`
	Questions    []string `json:"questions"`
	UsedQuestion string   `json:"usedQuestion"`
	State        string   `json:"state"`
}
