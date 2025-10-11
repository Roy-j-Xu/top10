package room

import "errors"

var (
	ErrInvalidRoom    = errors.New("invalid room configuration")
	ErrRoomNotFound   = errors.New("room not found")
	ErrPlayerNotFound = errors.New("player not found")
	ErrRoomExists     = errors.New("room name already exists")
	ErrPlayerExists   = errors.New("player name already exists in this room")
	ErrTimeout        = errors.New("operation timed out")
)
