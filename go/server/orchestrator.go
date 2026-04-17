package server

import "context"

type Orchestrator interface {
	Start(ctx context.Context) error
}
