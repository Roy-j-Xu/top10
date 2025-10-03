package core

import (
	"fmt"
)

func (room *Room) Run() {
	room.SetStatus(Waiting)
	room.WaitSync()
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

// Wait for everyone to join with AddPlayerSync
func (room *Room) WaitSync() {
	numOfReady := 0
	playerReady := make(map[int]bool)

	for {
		playerID := <-room.readyChan
		roomSize := room.SizeSync()
		if !playerReady[playerID] {
			playerReady[playerID] = true
			numOfReady++
			room.Broadcast(fmt.Sprintf("Player %d ready", playerID), Broadcast)
		}

		if numOfReady == roomSize {
			room.continueChan <- true
			room.Broadcast("All players ready, game starts", Broadcast)
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
			room.Broadcast(fmt.Sprintf("Player %d ready", playerID), Broadcast)
		}

		if numOfReady == room.Size() {
			room.continueChan <- true
			room.Broadcast("All players ready, game continues", Broadcast)
			return
		}
	}
}

func (room *Room) NextRound() {
	room.generateQuestions()

	room.assignNumbers()

	room.GuesserID = (room.GuesserID + 1) % room.Size()
}

func (room *Room) assignNumbers() {
	for playerID, k := range randomKFromN(room.Size(), 10) {
		room.Players[playerID].Number = k
		room.Message(fmt.Sprint("Your number: ", k), playerID, AssignNumber)
	}
}

func (room *Room) generateQuestions() {
	room.Questions = RandomQuestions(4)
	room.Broadcast(QuestionsMsg{
		Questions: room.Questions,
		Guesser:   room.GuesserID,
	}, Questions)
}
