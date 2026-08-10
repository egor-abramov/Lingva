package out_docker

import (
	"context"
	"lingva/pkg/test"
	"lingva/services/runner/internal/core/domain"
	"log/slog"
	"strings"
	"testing"
)

type TestExecutingJob struct {
	domain.ExecutionJob `yaml:",inline"`
	Stdin               []string `yaml:"stdin"`
}

type ExecutionResult struct {
	Stdout string `yaml:"stdout"`
	Stderr string `yaml:"stderr"`
}

func TestDockerSandbox_Run(t *testing.T) {
	log := &slog.Logger{}
	runner, err := New(log)
	if err != nil {
		t.Fatal(err)
	}

	handler := func(req TestExecutingJob) (ExecutionResult, error) {
		stdinCh := make(chan string)

		go func() {
			defer close(stdinCh)

			for _, inputLine := range req.Stdin {
				if !strings.HasSuffix(inputLine, "\n") {
					inputLine += "\n"
				}
				stdinCh <- inputLine
			}
		}()

		ctx := context.Background()

		outCh, err := runner.Run(ctx, req.ExecutionJob, stdinCh)
		if err != nil {
			return ExecutionResult{}, err
		}

		var stdoutB, stderrB strings.Builder
		for chunk := range outCh {
			if chunk.IsStdErr {
				stderrB.WriteString(chunk.Data)
			} else {
				stdoutB.WriteString(chunk.Data)
			}
		}
		return ExecutionResult{
			Stdout: stdoutB.String(),
			Stderr: stderrB.String(),
		}, nil
	}
	test.RunTests(t, handler)
}
