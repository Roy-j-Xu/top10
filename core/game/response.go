package game

type TurnInfoResponse struct {
	Turn      int      `json:"turn"`
	Guesser   string   `json:"guesser"`
	Questions []string `json:"questions"`
}
