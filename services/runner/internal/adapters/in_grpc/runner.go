package in_grpc

import (
	"context"
	"lingva/services/runner/internal/core/domain"
)

type RunnerService interface {
	Execute(ctx context.Context, job domain.ExecutionJob, stdin <-chan string) (<-chan domain.OutputChunk, error)
}
