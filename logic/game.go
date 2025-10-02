package logic

import (
	"math/rand"
)

func (room *Room) Run() {
	for {
		<-room.continueChan
		switch room.Status {
		case Waiting:
			room.Status = Playing
			room.NewRound()
			room.Wait()
		case Playing:
		case Finished:
		}
	}
}

func (room *Room) NewRound() {
	room.Question = RandomQuestion()
	for i, k := range randomKFromN(room.Size(), 10) {
		room.Players[i].Number = k
	}
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
