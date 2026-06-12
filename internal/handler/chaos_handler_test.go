package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

func TestChaosHandler_Crash(t *testing.T) {
	if os.Getenv("BE_CRASH") == "1" {
		h := NewChaosHandler()
		req := httptest.NewRequest(http.MethodPost, "/chaos", nil)
		rr := httptest.NewRecorder()
		h.Crash(rr, req)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestChaosHandler_Crash")
	cmd.Env = append(os.Environ(), "BE_CRASH=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		// Expecting an exit error due to os.Exit(1)
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
