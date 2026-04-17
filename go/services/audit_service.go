package services

import "context"

type AuditService struct{}

func (s *AuditService) LogToolCall(ctx context.Context, toolName string, payload any) error {
	return nil
}
