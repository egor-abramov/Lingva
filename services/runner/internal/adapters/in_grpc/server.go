package in_grpc

import (
	"fmt"
	"io"
	"lingva/api/gen"
	l "lingva/pkg/lang"
	"lingva/services/runner/internal/core/domain"

	"google.golang.org/grpc"
)

type runnerAPI struct {
	gen.UnimplementedCodeRunServiceServer
	runner RunnerService
}

func Register(gRPC *grpc.Server, runner RunnerService) {
	gen.RegisterCodeRunServiceServer(gRPC, &runnerAPI{
		runner: runner,
	})
}

func (r *runnerAPI) ExecuteCode(stream gen.CodeRunService_ExecuteCodeServer) error {
	const op = "in_grpc.runnerAPI.ExecuteCode"

	setupReq, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	setup := setupReq.GetSetup()
	if setup == nil {
		return fmt.Errorf("%s: first message must contain setup", op)
	}

	ctx := stream.Context()
	job := domain.ExecutionJob{
		Lang:       l.Language(setup.GetLang()),
		TimeoutMS:  5000,
		SourceCode: setup.GetCode(),
	}
	stdin := make(chan string)
	stdout, err := r.runner.Execute(ctx, job, stdin)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		defer close(stdin)

		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			if chunk := req.GetStdinChunk(); chunk != "" {
				select {
				case stdin <- chunk:
				case <-stream.Context().Done():
					return
				}
			}
		}
	}()

	for chunk := range stdout {
		resp := &gen.ExecutionResponse{
			IsFinished: false,
		}

		if chunk.IsStdErr {
			resp.Stderr = chunk.Data
		} else {
			resp.Stdout = chunk.Data
		}

		if err := stream.Send(resp); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	_ = stream.Send(&gen.ExecutionResponse{
		IsFinished: true,
	})

	return nil
}
