package core

import (
	"fmt"
	"log"
	"strings"
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
			return
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
	room.GuesserID = (room.GuesserID + 1) % room.Size()
	room.Broadcast(fmt.Sprint("Gusser: player ", room.GuesserID))

	room.Questions = RandomQuestions(4)
	room.Message(strings.Join(room.Questions, "\n"), room.GuesserID)

	room.assignNumbers()
}

func (room *Room) assignNumbers() {
	for playerID, k := range randomKFromN(room.Size(), 10) {
		room.Players[playerID].Number = k
		room.Message(fmt.Sprint("Your number: ", k), playerID)
	}
}
