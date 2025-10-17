package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"top10/core"
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
	resp, err := http.Post(ts.URL+"/create-room", "application/json", bytes.NewReader(body))
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
