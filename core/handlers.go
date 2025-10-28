package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleNewRoom(gm *GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Name string `json:"roomName"`
			Size int    `json:"roomSize"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		err := gm.NewRoomSync(req.Name, req.Size)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not create room: %v", err), http.StatusBadRequest)
			return
		}

		writeJson(w, RoomInfoResponse{
			RoomName: req.Name,
			RoomSize: req.Size,
			Game:     "Top10",
			Players:  []string{},
		}, 200)
	}
}

func handleRoomInfo(gm *GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		roomName := r.URL.Query().Get("roomName")

		rm, err := gm.GetRoomSync(roomName)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to find room \"%s\": %v", roomName, err), http.StatusBadRequest)
			return
		}

		writeJson(w, RoomInfoResponse{
			RoomName: rm.ID,
			RoomSize: rm.MaxSize,
			Game:     "Top10",
			Players:  rm.GetAllPlayerIDsSync(),
		}, 200)
	}
}
