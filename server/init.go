package server

import (
	"net/http"
	"top10/core"
)

func InitServer(gm *core.GameManager) {
	http.HandleFunc("/api/create-room", handleNewRoom(gm))
	http.HandleFunc("/api/room-info", handleRoomInfo(gm))
	http.HandleFunc("/api/game-info", handleGameInfo(gm))
	// joining room and establish socket connection
	http.HandleFunc("/ws", joinHandler(gm))
}

func ServeFrontend() {
	fs := http.FileServer(http.Dir("./resource/frontend/dist"))
	http.Handle("/", fs)
}
