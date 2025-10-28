package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"top10/core"

	"github.com/gorilla/websocket"
)

func TestCreateRoom(t *testing.T) {
	gm := core.NewGameManager()
	gm.HandleHTTP()

	// Create a test server using gm’s handlers
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	body, _ := json.Marshal(map[string]any{
		"name": "test_room",
		"size": 4,
	})
	resp, err := http.Post(ts.URL+"/api/create-room", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	logResponse(t, resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %v", resp.Status)
	}

	// Check if the room exists
	if _, err := gm.GetRoomSync("test_room"); err != nil {
		t.Fatalf("room not found after creation: %v", err)
	}

	<-time.After(10 * time.Millisecond)
}

func TestJoinRoom(t *testing.T) {
	gm := core.NewGameManager()
	gm.HandleHTTP()
	ts := httptest.NewServer(http.DefaultServeMux)
	defer ts.Close()

	body, _ := json.Marshal(map[string]any{
		"name": "test_room",
		"size": 4,
	})
	resp, err := http.Post(ts.URL+"/api/create-room", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	<-time.After(10 * time.Millisecond)

	// socket connection
	conn := wsConnect(t, ts, "test_room")
	wait10msAnd(func() { conn.WriteJSON("player42") })
	wait10msAnd(func() { conn.Close() })

	<-time.After(10 * time.Millisecond)
}

func wsConnect(t *testing.T, ts *httptest.Server, roomName string) *websocket.Conn {
	u, _ := url.Parse(ts.URL)
	u.Scheme = "ws"
	u.Path = "/ws"
	q := u.Query()
	q.Set("room", roomName)
	u.RawQuery = q.Encode()

	// Connect with a real WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}

	return ws
}

func logResponse(t *testing.T, resp *http.Response) {
	var pretty bytes.Buffer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Indent(&pretty, body, "", "  "); err != nil {
		t.Logf("Response (not JSON): %s", string(body))
		return
	}
	t.Logf("Response JSON:\n%s", pretty.String())
}
