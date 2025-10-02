package core

import (
	"log"
	"math/rand"
)

type DebugMessager struct{}

func (d *DebugMessager) Broadcast(msg any) {
	log.Println("[Broadcast]", msg)
}

func (d *DebugMessager) Message(msg any, playerID int) {
	log.Printf("[Player %d] %s\n", playerID, msg)
}

func randomKFromN(k, n int) []int {
	if k > n {
		panic("k greater than n")
	}

	nums := make([]int, n)
	for i := range n {
		nums[i] = i
	}

	rand.Shuffle(n, func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})

	return nums[:k]
}
