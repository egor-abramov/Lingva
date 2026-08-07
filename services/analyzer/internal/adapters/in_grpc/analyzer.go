package in_grpc

import (
	"context"
	"lingva/services/analyzer/internal/core/domain"
)

type AnalyzerService interface {
	Analyze(ctx context.Context, job domain.AnalyzeJob) (domain.AnalyzeResult, error)
}
