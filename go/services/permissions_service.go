package services

import "context"

type PermissionsService struct{}

func (s *PermissionsService) CanExecute(ctx context.Context, tier string, toolName string) (bool, error) {
	return true, nil
}
