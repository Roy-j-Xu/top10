package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func wait10msAnd(f func()) {
	time.Sleep(10 * time.Millisecond)
	f()
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
