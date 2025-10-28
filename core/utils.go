package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode JSON: %v", err), http.StatusInternalServerError)
		return
	}
}

type RoomInfoResponse struct {
	RoomName string   `json:"roomName"`
	RoomSize int      `json:"roomSize"`
	Game     string   `json:"game"`
	Players  []string `json:"players"`
}
