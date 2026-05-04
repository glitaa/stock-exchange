package service

import (
	"context"

	"github.com/glitaa/stock-exchange/internal/domain"
)

// AuditService handles business logic related to retrieving the audit log.
type AuditService struct {
	repo domain.AuditLogRepository
}

// NewAuditService creates a new instance of AuditService.
func NewAuditService(repo domain.AuditLogRepository) *AuditService {
	return &AuditService{repo: repo}
}

// GetLog retrieves the entire audit log in order of occurrence.
func (s *AuditService) GetLog(ctx context.Context) ([]domain.LogEntry, error) {
	return s.repo.GetAll(ctx)
}
