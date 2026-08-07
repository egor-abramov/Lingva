package out_docker

import (
	"context"
	"lingva/pkg/test"
	"lingva/services/runner/internal/core/domain"
	"log/slog"
	"testing"
)

type TestExecutingJob struct {
	domain.ExecutionJob `yaml:",inline"`
	Stdin               []string `yaml:"stdin"`
}

func TestDockerSandbox_Run(t *testing.T) {
	log := &slog.Logger{}
	runner, err := New(log)
	if err != nil {
		t.Fatal(err)
	}

	handler := func(req TestExecutingJob) ([]domain.OutputChunk, error) {
		stdinCh := make(chan string)

		go func() {
			defer close(stdinCh)

			for _, inputLine := range req.Stdin {
				stdinCh <- inputLine
			}
		}()

		ctx := context.Background()

		outCh, err := runner.Run(ctx, req.ExecutionJob, stdinCh)
		if err != nil {
			return nil, err
		}

		var stdout, stderr []domain.OutputChunk
		for chunk := range outCh {
			if chunk.IsStdErr {
				stderr = append(stderr, chunk)
			} else {
				stdout = append(stdout, chunk)
			}
		}

		return append(stdout, stderr...), nil
	}
	test.RunTests(t, handler)
}
