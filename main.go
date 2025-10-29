package main

import (
	"log"
	"net/http"
	"top10/core"
	"top10/server"
)

func main() {
	gm := core.InitCore()
	server.InitServer(gm)
	server.ServeFrontend()

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
