package main

import (
	"log"
	"net/http"
	"top10/core"
	"top10/socket"
)

func main() {
	core.InitCore()

	room := core.NewRoom(nil)
	go room.Run()

	socket.InitSocketHandler(room)

	serveFrontend()

	log.Println("Server started")
	err := http.ListenAndServe("0.0.0.0:8080", nil)

	log.Println(err.Error())
}

func serveFrontend() {
	fs := http.FileServer(http.Dir("./resource"))
	http.Handle("/", fs)
}
