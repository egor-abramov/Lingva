package service

import (
	"context"
	"fmt"
	"lingva/services/runner/internal/core/domain"
	"lingva/services/runner/internal/core/ports"
	"time"
)

type Runner struct {
	sandbox ports.Sandbox
}

func New(sandbox ports.Sandbox) *Runner {
	return &Runner{sandbox: sandbox}
}

func (r *Runner) Execute(ctx context.Context, job domain.ExecutionJob, stdin <-chan string) (<-chan domain.OutputChunk, error) {
	const op = "core.service.Runner.Execute"

	timeoutDuration := time.Duration(job.TimeoutMS) * time.Millisecond
	execCtx, cancel := context.WithTimeout(ctx, timeoutDuration)

	outputChan, err := r.sandbox.Run(execCtx, job, stdin)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	proxyChan := make(chan domain.OutputChunk)

	go func() {
		defer close(proxyChan)
		defer cancel()

		for chunk := range outputChan {
			proxyChan <- chunk
		}
	}()

	return proxyChan, nil
}
