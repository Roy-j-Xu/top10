package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"top10/core"
	"top10/core/room"
	"top10/server"

	"github.com/gorilla/websocket"
)

func TestCreateRoom(t *testing.T) {
	gm := core.NewGameManager()
	server.InitServer(gm)

	// Create a test server using gm’s handlers
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	createRoom(t, ts, "test_room")

	// Check if the room exists
	r, err := gm.GetRoomSync("test_room")
	if err != nil {
		t.Fatalf("room not found after creation: %v", err)
	}
	r.Print()

	<-time.After(10 * time.Millisecond)
}

func TestJoinRoom(t *testing.T) {
	gm := core.NewGameManager()
	server.InitServer(gm)
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	createRoom(t, ts, "test_room")

	<-time.After(10 * time.Millisecond)

	// socket connection
	conn := joinRoom(t, ts, "test_room", "player42")

	wait10msAnd(func() { getRoomInfo(t, ts, "test_room") })

	wait10msAnd(func() { conn.Close() })

	<-time.After(10 * time.Millisecond)
}

func TestReady(t *testing.T) {
	gm := core.NewGameManager()
	server.InitServer(gm)
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	createRoom(t, ts, "test_room")

	ws := joinRoom(t, ts, "test_room", "player1")

	<-time.After(10 * time.Millisecond)

	ws.WriteJSON(map[string]any{
		"type": room.SP_READY,
		"msg":  "",
	})

	<-time.After(10 * time.Millisecond)

	g, err := gm.GetGameSync("test_room")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(g.GetGameInfoSync())
}

func createRoom(t *testing.T, ts *httptest.Server, roomName string) {
	body, _ := json.Marshal(map[string]any{
		"roomName": roomName,
		"roomSize": 4,
	})
	resp, err := http.Post(ts.URL+"/api/create-room", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()
}

func joinRoom(t *testing.T, ts *httptest.Server, roomName string, playerName string) *websocket.Conn {
	u, _ := url.Parse(ts.URL)
	u.Scheme = "ws"
	u.Path = "/ws"
	q := u.Query()
	q.Set("roomName", roomName)
	q.Set("playerName", playerName)
	u.RawQuery = q.Encode()

	// Connect with a real WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}

	return ws
}

func getRoomInfo(t *testing.T, ts *httptest.Server, roomName string) {
	resp, err := http.Get(ts.URL + fmt.Sprintf("/api/room-info?roomName=%s", roomName))
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	logResponse(t, resp)
	defer resp.Body.Close()
}
