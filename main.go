package main

import (
	"log"
	"net/http"
	"top10/core"
	"top10/core/room"
)

func main() {
	core.InitCore()

	room, err := room.NewRoomDebug("room", 10)
	go room.Run()

	room.InitSocketHandler(room)

	serveFrontend()

	log.Println("Server started")
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	log.Println(err.Error())
}

func serveFrontend() {
	fs := http.FileServer(http.Dir("./resource"))
	http.Handle("/", fs)
}
