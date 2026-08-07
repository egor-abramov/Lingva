package out_docker

import "lingva/services/runner/internal/core/domain"

type channelWriter struct {
	outChan chan<- domain.OutputChunk
	isErr   bool
}

func (cw *channelWriter) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))
	copy(buf, p)
	cw.outChan <- domain.OutputChunk{
		Data:     string(buf),
		IsStdErr: cw.isErr,
	}
	return len(p), nil
}
