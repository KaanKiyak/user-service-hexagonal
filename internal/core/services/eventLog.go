package services

import (
	"context"
	"user-service-hexagonal/internal/core/domain"
	"user-service-hexagonal/internal/core/ports"
)

type eventLogService struct {
	repo ports.EventLogRepository
}

func NewEventLogService(repo ports.EventLogRepository) *eventLogService {
	return &eventLogService{
		repo: repo,
	}
}

func (s *eventLogService) SaveEvent(ctx context.Context, log *domain.EventLog) error {
	// İstersen burada log türüne göre filtreleme, ön işlem vs. yapabilirsin
	return s.repo.Save(ctx, log)
}
