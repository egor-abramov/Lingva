package service

import (
	"context"
	"fmt"
	"lingva/services/analyzer/internal/core/domain"
	"lingva/services/analyzer/internal/core/ports"
	"time"
)

type Analyzer struct {
	sandbox ports.Sandbox
}

func New(sandbox ports.Sandbox) *Analyzer {
	return &Analyzer{
		sandbox: sandbox,
	}
}

func (a *Analyzer) Analyze(ctx context.Context, job domain.AnalyzeJob) (domain.AnalyzeResult, error) {
	const op = "core.service.Analyzer.Analyze"

	timeoutDuration := time.Duration(job.TimeoutMS) * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)

	res, err := a.sandbox.Run(ctx, job)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cancel()
	return res, nil
}
