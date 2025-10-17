package room

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

func (r *Room) WaitUntilGetMessage(
	playerID string,
	types ...string,
) (Message, error) {
	player, err := r.GetPlayerSync(playerID)
	if err != nil {
		return Message{}, fmt.Errorf("waiting for message from player %s: %w", playerID, err)
	}
	for {
		select {
		case message := <-player.msgChan:
			if slices.Contains(types, message.Type) {
				return message, nil
			}
		case <-time.After(r.Timeout):
			return Message{}, fmt.Errorf("waiting for message from player %s: %w", playerID, ErrTimeout)
		}
	}
}

func (r *Room) WaitAndHandleAllMessages(
	handler func(Message),
	types ...string,
) error {
	var wg sync.WaitGroup
	done := make(chan struct{})

	for _, player := range r.Players {
		wg.Add(1)
		go func(p *Player) {
			defer wg.Done()

			for {
				select {
				case msg := <-p.msgChan:
					if slices.Contains(types, msg.Type) {
						handler(msg)
						return
					}
				case <-time.After(r.Timeout): // player listener must also timeout
					return
				}
			}
		}(player)
	}

	// Close done channel when all goroutines finish
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(r.Timeout):
		return fmt.Errorf("timeout waiting for all messages")
	}
}

func (r *Room) WaitForAllMessages(types ...string) error {
	return r.WaitAndHandleAllMessages(
		func(msg Message) {},
		types...,
	)
}
