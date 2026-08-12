package out_docker

import (
	"context"
	"fmt"
	"io"
	l "lingva/pkg/lang"
	"lingva/services/runner/internal/core/domain"
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

func (s *DockerSandbox) Run(ctx context.Context, job domain.ExecutionJob, stdin <-chan string) (<-chan domain.OutputChunk, error) {
	const op = "adapter.out_docker.Run"

	var image string
	var cmd []string

	switch job.Lang {
	case l.Python:
		image = "python:3.11-alpine"
		cmd = []string{"python3", "-u", "-c", job.SourceCode}
	case l.Go:
		image = "golang:1.24-alpine"
		cmd = []string{
			"sh",
			"-c",
			"printf \"%s\" \"$1\" > main.go && go run main.go",
			"run-go",
			job.SourceCode,
		}
	case l.C:
		image = "gcc:latest"
		cmd = []string{
			"sh",
			"-c",
			"printf \"%s\" \"$1\" > main.c && gcc -Wall main.c -o main && ./main",
			"run-c",
			job.SourceCode,
		}
	case l.Cpp:
		image = "gcc:latest"
		cmd = []string{
			"sh",
			"-c",
			"printf \"%s\" \"$1\" > main.cpp && g++ -Wall main.cpp -o main && ./main",
			"run-cpp",
			job.SourceCode,
		}
	default:
		return nil, fmt.Errorf("%s: unsupported language %s", op, job.Lang)
	}

	config := &container.Config{
		Image:           image,
		Cmd:             cmd,
		Tty:             false,
		OpenStdin:       true,
		StdinOnce:       true,
		AttachStdin:     true,
		AttachStdout:    true,
		AttachStderr:    true,
		NetworkDisabled: true,

		WorkingDir: "/workspace",

		Env: []string{
			"GOCACHE=/gocache",
			"TMPDIR=/workspace",
			"CGO_ENABLED=0",
		},
	}
	pidsLimit := int64(512)

	hostConfig := &container.HostConfig{
		AutoRemove: true,

		Binds: []string{
			"/var/lib/lingva/gocache:/gocache",
		},

		Resources: container.Resources{
			Memory:     2048 * 1024 * 1024,
			MemorySwap: -1,
			NanoCPUs:   1_000_000_000,
			PidsLimit:  &pidsLimit,
		},

		ReadonlyRootfs: true,
		CapDrop:        []string{"ALL"},
		SecurityOpt:    []string{"no-new-privileges:true"},

		Tmpfs: map[string]string{
			"/workspace": "rw,exec,nosuid,size=256m",
		},
	}
	resp, err := s.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	attachResp, err := s.cli.ContainerAttach(ctx, resp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		attachResp.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	outChan := make(chan domain.OutputChunk)

	go func() {
		for input := range stdin {
			_, err := attachResp.Conn.Write([]byte(input))
			if err != nil {
				s.log.Error(fmt.Sprintf("%s: failed to write to docker container: %s", op, err))
				break
			}
		}
	}()

	go func() {
		defer attachResp.Close()
		defer close(outChan)

		stdoutWriter := &channelWriter{outChan: outChan, isErr: false}
		stderrWriter := &channelWriter{outChan: outChan, isErr: true}

		_, err := stdcopy.StdCopy(stdoutWriter, stderrWriter, attachResp.Reader)
		if err != nil && err != io.EOF {
			s.log.Error(fmt.Sprintf("failed to copy stdout from docker container: %s", err))
		}
	}()

	return outChan, nil
}
