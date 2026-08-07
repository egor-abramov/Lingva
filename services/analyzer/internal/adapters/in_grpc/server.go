package in_grpc

import (
	"context"
	"fmt"
	"lingva/api/gen"
	"lingva/services/analyzer/internal/core/domain"

	"google.golang.org/grpc"
)

type analyzerAPI struct {
	gen.UnimplementedCodeAnalyzeServiceServer
	analyzer AnalyzerService
}

func Register(gRPC *grpc.Server, analyzer AnalyzerService) {
	gen.RegisterCodeAnalyzeServiceServer(gRPC, &analyzerAPI{
		analyzer: analyzer,
	})
}

func (a *analyzerAPI) Analyze(ctx context.Context, req *gen.AnalyzeRequest) (*gen.AnalyzeResponse, error) {
	const op = "in_grpc.analyzerAPI.Analyze"

	lang := req.GetLang()
	code := req.GetCode()
	job := domain.AnalyzeJob{
		Lang:      domain.Language(lang),
		Code:      code,
		TimeoutMS: 5000,
	}

	resp, err := a.analyzer.Analyze(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var pbDiagnostics []*gen.Diagnostic
	for _, diag := range resp {
		pbDiag := &gen.Diagnostic{
			Line:     diag.Line,
			Column:   diag.Column,
			Message:  diag.Message,
			Severity: gen.Severity(diag.Severity),
		}
		pbDiagnostics = append(pbDiagnostics, pbDiag)
	}
	pbResp := &gen.AnalyzeResponse{
		Diagnostics: pbDiagnostics,
	}
	return pbResp, nil
}
