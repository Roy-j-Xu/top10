package main

import (
	"log"
	"net/http"
	"top10/core"
)

func main() {
	core.InitCore()

	serveFrontend()

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

func serveFrontend() {
	fs := http.FileServer(http.Dir("./resource/frontend/dist"))
	http.Handle("/", fs)
}
