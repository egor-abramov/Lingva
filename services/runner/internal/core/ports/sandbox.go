package ports

import (
	"context"
	"lingva/services/runner/internal/core/domain"
)

type Sandbox interface {
	Run(ctx context.Context, job domain.ExecutionJob, stdin <-chan string) (<-chan domain.OutputChunk, error)
}
