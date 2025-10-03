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

	socket.InitSocketHandler(room)

	serveFrontend()

	log.Println("Server started")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func serveFrontend() {
	fs := http.FileServer(http.Dir("./resource"))
	http.Handle("/", fs)
}
