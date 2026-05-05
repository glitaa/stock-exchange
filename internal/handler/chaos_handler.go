package handler

import (
	"log"
	"net/http"
	"os"
)

// ChaosHandler handles requests for injecting faults into the system.
type ChaosHandler struct{}

// NewChaosHandler creates a new instance of ChaosHandler.
func NewChaosHandler() *ChaosHandler {
	return &ChaosHandler{}
}

// Crash handles the POST /chaos request by killing the application process.
func (h *ChaosHandler) Crash(w http.ResponseWriter, r *http.Request) {
	log.Println("ChaosHandler triggered. Terminating process...")
	os.Exit(1)
}
