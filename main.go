package main

import (
	"net/http"
	"top10/core"
	"top10/socket"
)

func main() {
	core.InitCore()

	room := core.NewRoom(nil)

	socket.InitSocketHandler(room)

	serveFrontend()

	http.ListenAndServe("0.0.0.0:2357", nil)
}

func serveFrontend() {
	fs := http.FileServer(http.Dir("./resource"))
	http.Handle("/", fs)
}
