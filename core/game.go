package core

import (
	"log"
	"math/rand"
)

func (room *Room) Run() {
	room.SetStatus(Waiting)
	room.Wait()
	for {
		<-room.continueChan
		switch room.Status {
		case Waiting:
			room.SetStatus(Playing)
			room.NextRound()
			room.Wait()
		case Playing:
			room.SetStatus(Playing)
			room.NextRound()
			room.Wait()
		case Finished:
		}
	}
}

// Wait for everyone to ready
func (room *Room) Wait() {
	numOfReady := 0
	playerReady := make([]bool, room.Size())

	for {
		playerID := <-room.readyChan
		if !playerReady[playerID] {
			playerReady[playerID] = true
			numOfReady++
			log.Printf("Player %d ready (%d of %d)", playerID, numOfReady, room.Size())
		}

		if numOfReady == room.Size() {
			room.continueChan <- true
			log.Println("All players ready, game continues")
			return
		}
	}
}

func (room *Room) NextRound() {
	room.Questions = RandomQuestions(4)
	assignNumber(room.Players, room.Size())
	room.GuesserID = (room.GuesserID + 1) % room.Size()
}

func assignNumber(players []*Player, size int) {
	for i, k := range randomKFromN(size, 10) {
		players[i].Number = k
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
