package ports

import (
	"context"
	"lingva/services/analyzer/internal/core/domain"
)

type Sandbox interface {
	Run(ctx context.Context, job domain.AnalyzeJob) (domain.AnalyzeResult, error)
}
