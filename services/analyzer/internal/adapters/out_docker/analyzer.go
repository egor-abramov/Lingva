package out_docker

import (
	"bytes"
	"context"
	"fmt"
	l "lingva/pkg/lang"
	"lingva/services/analyzer/internal/adapters/out_docker/parsers"
	"lingva/services/analyzer/internal/core/domain"
	"log/slog"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerSandbox struct {
	cli *client.Client
	log *slog.Logger
}

func New(log *slog.Logger) (*DockerSandbox, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerSandbox{cli: cli, log: log}, nil
}

func (a *DockerSandbox) Run(ctx context.Context, job domain.AnalyzeJob) (domain.AnalyzeResult, error) {
	const op = "out_docker.DockerAnalyzer.Run"

	var image string
	var cmd []string

	switch job.Lang {
	case l.Python:
		image = "alpine/flake8:latest"
		cmd = []string{"-"}
	case l.Go:
		image = "golang:1.24-alpine"
		cmd = []string{
			"sh",
			"-c",
			"cat > main.go && go build -o /dev/null main.go && go vet main.go",
		}
	case l.C:
		image = "gcc:latest"
		cmd = []string{"gcc", "-fsyntax-only", "-Wall", "-Wextra", "-x", "c", "-"}
	case l.Cpp:
		image = "gcc:latest"
		cmd = []string{"g++", "-fsyntax-only", "-Wall", "-Wextra", "-x", "c++", "-"}
	default:
		return nil, fmt.Errorf("%s: unsupported language %s", op, job.Lang)

	}

	config := &container.Config{
		Image:        image,
		Cmd:          cmd,
		Tty:          false,
		OpenStdin:    true,
		StdinOnce:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}
	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}

	if job.Lang == l.Go {
		config.WorkingDir = "/workspace"
		config.Env = []string{
			"CGO_ENABLED=0",
			"GOCACHE=/gocache",
			"TMPDIR=/workspace",
		}
		hostConfig.Binds = []string{
			"/var/lib/lingva/gocache:/gocache",
		}
		hostConfig.Tmpfs = map[string]string{
			"/workspace": "rw,exec,nosuid,size=128m",
		}
	}

	resp, err := a.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create container: %s", op, err)
	}

	attachResp, err := a.cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
		Stdin:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: failed to attach container: %s", op, err)
	}
	defer attachResp.Close()

	if err := a.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("%s: failed to start container: %s", op, err)
	}

	if _, err := attachResp.Conn.Write([]byte(job.Code)); err != nil {
		return nil, fmt.Errorf("%s: failed to write code to container: %s", op, err)
	}

	_ = attachResp.CloseWrite()

	var stdout, stderr bytes.Buffer
	copyErrCh := make(chan error, 1)
	go func() {
		_, err := stdcopy.StdCopy(&stdout, &stderr, attachResp.Reader)
		copyErrCh <- err
	}()

	statusCh, errCh := a.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("%s: container wait error: %w", op, err)
		}
	case <-statusCh:
	}

	if err := <-copyErrCh; err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to copy output", op), "error", err)
	}

	switch job.Lang {
	case l.Python:
		return parsers.Flake8(stdout.String(), stderr.String()), nil //[cite: 15]
	case l.Go:
		return parsers.Go(stderr.String()), nil
	case l.C, l.Cpp:
		return parsers.Gcc(stderr.String()), nil
	}
	return nil, fmt.Errorf("%s: unsupported language %s", op, job.Lang)
}
