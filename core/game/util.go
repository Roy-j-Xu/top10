package game

import (
	"fmt"
	"log"
	"math/rand"
	"top10/core/room"
)

func (g *Game) SetStatus(status *Status) {
	log.Printf("Game status: %s -> %s", g.Status.Name, status.Name)
	g.Status = status
	g.Status.OnStatus(g)
}

func (g *Game) Room() *room.Room {
	return g.room
}

func (g *Game) Size() int {
	return len(g.PlayerStates)
}

func (g *Game) AddNewPlayerState(playerID string) error {
	if _, ok := g.PlayerStates[playerID]; ok {
		return fmt.Errorf("player %s already has a game state", playerID)
	}
	g.PlayerStates[playerID] = &PlayerState{}
	return nil
}

func randomKFromN(k, n int) []int {
	if k > n {
		panic("k greater than n")
	}

	nums := make([]int, n)
	for i := range n {
		nums[i] = i + 1
	}

	rand.Shuffle(n, func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})

	return nums[:k]
}

func (g *Game) Print() {
	log.Println()
	log.Println("----- Game State -----")
	log.Printf("Room ID: %s", g.Room().ID)
	log.Printf("Game status: %s", g.Status.Name)
	log.Printf("Turn order: %v", g.TurnOrder)
	log.Printf("Turn number: %d", g.TurnNumber)
	log.Printf("Question: %s", g.Question)
	log.Printf("Guesser ID: %s", g.GuesserID)
	log.Println("----------------------")
	log.Println()
}
