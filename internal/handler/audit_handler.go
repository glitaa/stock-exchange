package handler

import (
	"net/http"

	"github.com/glitaa/stock-exchange/internal/domain"
	"github.com/glitaa/stock-exchange/internal/service"
)

// AuditHandler handles HTTP requests for the audit log.
type AuditHandler struct {
	auditService *service.AuditService
}

// NewAuditHandler creates a new instance of AuditHandler.
func NewAuditHandler(s *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: s}
}

// auditResponse defines the expected JSON structure for the audit log.
type auditResponse struct {
	Log []domain.LogEntry `json:"log"`
}

// GetLog handles the GET request to retrieve the full audit log.
func (h *AuditHandler) GetLog(w http.ResponseWriter, r *http.Request) {
	logs, err := h.auditService.GetLog(r.Context())
	if err != nil {
		respondWithError(w, err)
		return
	}

	if logs == nil {
		logs = []domain.LogEntry{}
	}

	response := auditResponse{Log: logs}
	respondWithJSON(w, http.StatusOK, response)
}
